package solana

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go/rpc/ws"

	"github.com/dcaf-protocol/drip/internal/configs"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type Solana interface {
	MintToWallet(context.Context, string, string, uint64) (string, error)
	signAndBroadcast(context.Context, ...solana.Instruction) (string, error)

	// Wrappers

	GetWalletPubKey() solana.PublicKey
	getWalletPrivKey() solana.PrivateKey
	GetVersion(context.Context) (*rpc.GetVersionResult, error)
	GetAccountInfo(context.Context, solana.PublicKey) (*rpc.GetAccountInfoResult, error)
	ProgramSubscribe(context.Context, string, func(string, []byte)) error
}

func CreateSolanaClient(
	config *configs.AppConfig,
) (Solana, error) {
	return createsolanaImplClient(config)
}

type solanaImpl struct {
	environment configs.Environment
	client      *rpc.Client
	wallet      *solana.Wallet
}

func createsolanaImplClient(
	config *configs.AppConfig,
) (solanaImpl, error) {
	url := getURL(config.Environment)
	solanaClient := solanaImpl{
		client:      rpc.NewWithCustomRPCClient(rpc.NewWithRateLimit(url, 10)),
		environment: config.Environment,
	}
	resp, err := solanaClient.GetVersion(context.Background())
	if err != nil {
		log.WithError(err).Fatalf("failed to get clients version info")
		return solanaImpl{}, err
	}
	log.
		WithFields(log.Fields{
			"version": resp.SolanaCore,
			"url":     url,
		}).
		Info("created solana clients")

	var accountBytes []byte
	if err := json.Unmarshal([]byte(config.Wallet), &accountBytes); err != nil {
		return solanaImpl{}, err
	}
	priv := base58.Encode(accountBytes)
	solWallet, err := solana.WalletFromPrivateKeyBase58(priv)
	if err != nil {
		return solanaImpl{}, err
	}
	solanaClient.wallet = solWallet
	log.
		WithFields(logrus.Fields{"publicKey": solanaClient.GetWalletPubKey()}).
		Infof("loaded wallet")

	return solanaClient, nil
}

func (s solanaImpl) MintToWallet(
	ctx context.Context, mint, destWallet string, amount uint64,
) (string, error) {
	mintPubKey := solana.MustPublicKeyFromBase58(mint)
	destWalletPubKey := solana.MustPublicKeyFromBase58(destWallet)
	destAccount, _, err := solana.FindAssociatedTokenAddress(destWalletPubKey, mintPubKey)
	if err != nil {
		return "", err
	}
	var instructions []solana.Instruction
	if _, err := s.GetTokenAccountBalance(ctx, destAccount, "confirmed"); err != nil {
		txBuilder := associatedtokenaccount.NewCreateInstructionBuilder()
		txBuilder.SetMint(mintPubKey)
		txBuilder.SetPayer(s.GetWalletPubKey())
		txBuilder.SetWallet(destWalletPubKey)
		instruction, err := txBuilder.ValidateAndBuild()
		if err != nil {
			return "", err
		}
		instructions = append(instructions, instruction)
	}
	txBuilder := token.NewMintToInstructionBuilder()
	txBuilder.SetAuthorityAccount(s.GetWalletPubKey())
	txBuilder.SetDestinationAccount(destAccount)
	txBuilder.SetMintAccount(solana.MustPublicKeyFromBase58(mint))
	txBuilder.SetAmount(amount)
	tx, err := txBuilder.ValidateAndBuild()
	if err != nil {
		return "", err
	}
	instructions = append(instructions, tx)
	return s.signAndBroadcast(ctx, instructions...)
}

// TODO(Mocha): Pass in an error channel so that subscribers can handle errors
func (s solanaImpl) ProgramSubscribe(
	ctx context.Context, program string, onReceive func(string, []byte),
) error {
	url := getWSURL(s.environment)
	client, err := ws.Connect(ctx, url)
	if err != nil {
		return err
	}
	sub, err := client.ProgramSubscribeWithOpts(
		solana.MustPublicKeyFromBase58(program),
		rpc.CommitmentRecent,
		solana.EncodingBase64Zstd,
		nil,
	)
	if err != nil {
		return err
	}
	go func() {
		defer sub.Unsubscribe()
		for {
			msg, err := sub.Recv()
			if err != nil {
				log.
					WithFields(log.Fields{
						"event": program,
					}).
					Error("failed to get next msg from event ws")
				continue
			}
			if msg.Value.Account == nil || msg.Value.Account.Data == nil {
				log.
					WithFields(log.Fields{
						"event": program,
					}).
					Warning("event ws msg account or account data is nil")
				continue
			}
			decodedBinary := msg.Value.Account.Data.GetBinary()
			if decodedBinary == nil {
				log.
					WithFields(log.Fields{
						"event": program,
					}).
					Warning("event ws msg decoded binary is nil")
				continue
			}
			onReceive(msg.Value.Pubkey.String(), decodedBinary)
		}
	}()
	return nil
}

////////////////////////////////////////////////////////////
/// Wallet Wrapper
////////////////////////////////////////////////////////////

func (s solanaImpl) GetWalletPubKey() solana.PublicKey {
	return s.wallet.PublicKey()
}

func (s solanaImpl) getWalletPrivKey() solana.PrivateKey {
	return s.wallet.PrivateKey
}

////////////////////////////////////////////////////////////
/// RPC Client Wrapper
////////////////////////////////////////////////////////////

func (s solanaImpl) GetTokenAccountBalance(
	ctx context.Context, destAccount solana.PublicKey, commitmentType rpc.CommitmentType,
) (*rpc.GetTokenAccountBalanceResult, error) {
	return s.client.GetTokenAccountBalance(ctx, destAccount, commitmentType)
}

func (s solanaImpl) GetAccountInfo(
	ctx context.Context, account solana.PublicKey,
) (*rpc.GetAccountInfoResult, error) {
	return s.client.GetAccountInfo(ctx, account)
}

func (s solanaImpl) GetVersion(ctx context.Context) (*rpc.GetVersionResult, error) {
	return s.client.GetVersion(ctx)
}

func (s solanaImpl) signAndBroadcast(
	ctx context.Context, instructions ...solana.Instruction,
) (string, error) {
	recent, err := s.client.GetRecentBlockhash(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		return "", err
	}
	logFields := logrus.Fields{"numInstructions": len(instructions), "block": recent.Value.Blockhash}

	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(s.GetWalletPubKey()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction, err %s", err)
	}
	logrus.WithFields(logFields).Infof("built transaction")

	if _, err := tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if s.GetWalletPubKey().Equals(key) {
				priv := s.getWalletPrivKey()
				return &priv
			}

			return nil
		},
	); err != nil {
		return "", fmt.Errorf("failed to sign transaction, err %s", err)
	}
	logrus.WithFields(logFields).Info("signed transaction")

	txHash, err := s.client.SendTransactionWithOpts(
		ctx, tx, false, rpc.CommitmentConfirmed,
	)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction, err %s", err)
	}
	logFields["txHash"] = txHash

	return txHash.String(), nil
}

func getURL(env configs.Environment) string {
	switch env {
	case configs.DevnetEnv:
		return rpc.DevNet_RPC
	case configs.MainnetEnv:
		return rpc.MainNetBeta_RPC
	case configs.NilEnv:
		fallthrough
	case configs.LocalnetEnv:
		fallthrough
	default:
		return rpc.LocalNet_RPC
	}
}

func getWSURL(env configs.Environment) string {
	switch env {
	case configs.DevnetEnv:
		return rpc.DevNet_WS
	case configs.MainnetEnv:
		return rpc.MainNetBeta_WS
	case configs.NilEnv:
		fallthrough
	case configs.LocalnetEnv:
		fallthrough
	default:
		return rpc.LocalNet_WS
	}
}
