package solanaclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type Wallet struct{}

type Solana struct {
	Client *rpc.Client
	Wallet *solana.Wallet
}

func NewSolanaClient(
	config *configs.Config,
) (*Solana, error) {
	var solanaClient Solana

	url := getURL(config.Environment)
	solClient := *rpc.New(url)
	resp, err := solClient.GetVersion(context.Background())
	if err != nil {
		log.WithError(err).Fatalf("failed to get client version info")
		return nil, err
	}
	log.
		WithFields(log.Fields{
			"version": resp.SolanaCore,
			"url":     url}).
		Info("created solana client")

	solanaClient.Client = &solClient
	if configs.IsLocal(config.Environment) {
		log.Infof("creating and funding test wallet")
		solanaClient.Wallet = solana.NewWallet()
		solClient.RequestAirdrop(context.Background(), solanaClient.Wallet.PublicKey(), solana.LAMPORTS_PER_SOL*1, "confirmed")
	} else {
		var accountBytes []byte
		if err := json.Unmarshal([]byte(config.Wallet), &accountBytes); err != nil {
			return nil, err
		}
		priv := base58.Encode(accountBytes)
		solWallet, err := solana.WalletFromPrivateKeyBase58(priv)
		if err != nil {
			return nil, err
		}
		solanaClient.Wallet = solWallet
	}
	log.
		WithFields(logrus.Fields{"publicKey": solanaClient.Wallet.PublicKey()}).
		Infof("loaded wallet")
	return &solanaClient, nil
}

func (s *Solana) MintToAccount(
	ctx context.Context, mint, destinationAccount string, amount uint64, instructions ...solana.Instruction,
) (string, error) {
	txBuilder := token.NewMintToInstructionBuilder()
	txBuilder.SetAuthorityAccount(s.Wallet.PublicKey())
	txBuilder.SetDestinationAccount(solana.MustPublicKeyFromBase58(destinationAccount))
	txBuilder.SetMintAccount(solana.MustPublicKeyFromBase58(mint))
	txBuilder.SetAmount(amount)
	tx, err := txBuilder.ValidateAndBuild()
	if err != nil {
		return "", err
	}
	instructions = append(instructions, tx)
	return s.signAndBroadcast(ctx, instructions...)
}

func (s *Solana) MintToWallet(
	ctx context.Context, mint, destWallet string, amount uint64,
) (string, error) {
	mintPubKey := solana.MustPublicKeyFromBase58(mint)
	destWalletPubKey := solana.MustPublicKeyFromBase58(destWallet)
	destAccount, _, err := solana.FindAssociatedTokenAddress(destWalletPubKey, mintPubKey)
	if err != nil {
		return "", err
	}
	var instructions []solana.Instruction
	if _, err := s.Client.GetTokenAccountBalance(ctx, destAccount, "confirmed"); err != nil {
		txBuilder := associatedtokenaccount.NewCreateInstructionBuilder()
		txBuilder.SetMint(mintPubKey)
		txBuilder.SetPayer(s.Wallet.PublicKey())
		txBuilder.SetWallet(destWalletPubKey)
		instruction, err := txBuilder.ValidateAndBuild()
		if err != nil {
			return "", err
		}
		instructions = append(instructions, instruction)
	}
	return s.MintToAccount(ctx, mint, destAccount.String(), amount, instructions...)
}

func (s *Solana) signAndBroadcast(
	ctx context.Context, instructions ...solana.Instruction,
) (string, error) {
	recent, err := s.Client.GetRecentBlockhash(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		return "", err
	}
	logFields := logrus.Fields{"numInstructions": len(instructions), "block": recent.Value.Blockhash}

	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(s.Wallet.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction, err %s", err)
	}
	logrus.WithFields(logFields).Infof("built transaction")

	if _, err := tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if s.Wallet.PublicKey().Equals(key) {
				return &s.Wallet.PrivateKey
			}
			return nil
		},
	); err != nil {
		return "", fmt.Errorf("failed to sign transaction, err %s", err)
	}
	logrus.WithFields(logFields).Info("signed transaction")

	txHash, err := s.Client.SendTransactionWithOpts(
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
