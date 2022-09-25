package solana

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dcaf-labs/drip/pkg/service/configs"

	"github.com/dcaf-labs/drip/pkg/service/utils"

	bin "github.com/gagliardetto/binary"
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/mr-tron/base58"
	"github.com/sirupsen/logrus"
)

const ErrNotFound = "not found"

type Solana interface {
	MintToWallet(context.Context, string, string, uint64) (string, error)
	signAndBroadcast(context.Context, rpc.CommitmentType, ...solana.Instruction) (string, error)
	GetUserBalances(context.Context, string) (*rpc.GetTokenAccountsResult, error)
	GetAccount(context.Context, string, interface{}) error
	GetAccounts(context.Context, []string, func(string, []byte)) error
	GetProgramAccounts(context.Context, string) ([]string, error)
	GetAccountInfo(context.Context, solana.PublicKey) (*rpc.GetAccountInfoResult, error)
	ProgramSubscribe(context.Context, string, func(string, []byte)) error

	GetTokenMetadataAccount(ctx context.Context, mintAddress string) (token_metadata.Metadata, error)
	GetTokenMint(ctx context.Context, mintAddress string) (token.Mint, error)
	GetLargestTokenAccounts(ctx context.Context, mint string) ([]*rpc.TokenLargestAccountsResult, error)

	GetWalletPubKey() solana.PublicKey
	getWalletPrivKey() solana.PrivateKey
	GetVersion(context.Context) (*rpc.GetVersionResult, error)
	GetNetwork() configs.Network
}

func NewSolanaClient(
	config *configs.AppConfig,
) (Solana, error) {
	return createClient(config)
}

type impl struct {
	network configs.Network
	client  *rpc.Client
	wallet  *solana.Wallet
}

func (s impl) GetNetwork() configs.Network {
	return s.network
}

func createClient(
	config *configs.AppConfig,
) (impl, error) {
	url := GetURL(config.Network)
	solanaClient := impl{
		client:  rpc.NewWithCustomRPCClient(rpc.NewWithRateLimit(url, 2)),
		network: config.Network,
	}
	resp, err := solanaClient.GetVersion(context.Background())
	if err != nil {
		logrus.WithError(err).Fatalf("failed to get clients version info")
		return impl{}, err
	}
	logrus.
		WithFields(logrus.Fields{
			"version": resp.SolanaCore,
			"url":     url,
		}).
		Info("created solana clients")

	var accountBytes []byte
	if err := json.Unmarshal([]byte(config.Wallet), &accountBytes); err != nil {
		return impl{}, err
	}
	priv := base58.Encode(accountBytes)
	solWallet, err := solana.WalletFromPrivateKeyBase58(priv)
	if err != nil {
		return impl{}, err
	}
	solanaClient.wallet = solWallet
	logrus.
		WithFields(logrus.Fields{"publicKey": solanaClient.GetWalletPubKey()}).
		Infof("loaded wallet")

	return solanaClient, nil
}

func (s impl) GetAccounts(ctx context.Context, addresses []string, decode func(string, []byte)) error {
	var pubkeys []solana.PublicKey
	for _, address := range addresses {
		pubkey, err := solana.PublicKeyFromBase58(address)
		if err != nil {
			return err
		}
		pubkeys = append(pubkeys, pubkey)
	}
	resp, err := s.client.GetMultipleAccountsWithOpts(ctx, pubkeys, &rpc.GetMultipleAccountsOpts{
		Encoding:   solana.EncodingBase64,
		Commitment: "confirmed",
	})
	if err != nil {
		logrus.
			WithError(err).
			Errorf("couldn't get multiple account infos")
		return err
	}
	if len(resp.Value) != len(addresses) {
		return fmt.Errorf("response does not match length of addresses")
	}
	for i, val := range resp.Value {
		if val == nil || val.Data == nil {
			continue
		}
		decode(addresses[i], val.Data.GetBinary())
	}
	return nil
}

func (s impl) GetTokenMint(ctx context.Context, mintAddress string) (token.Mint, error) {
	var tokenMint token.Mint
	err := s.GetAccount(ctx, mintAddress, &tokenMint)
	return tokenMint, err
}

func (s impl) GetTokenMetadataAccount(ctx context.Context, mintAddress string) (token_metadata.Metadata, error) {
	var tokenMetadata token_metadata.Metadata

	tokenMetadataAddress, err := utils.GetTokenMetadataPDA(mintAddress)
	if err != nil {
		return tokenMetadata, err
	}

	resp, err := s.client.GetAccountInfoWithOpts(
		ctx,
		solana.MustPublicKeyFromBase58(tokenMetadataAddress),
		&rpc.GetAccountInfoOpts{
			Encoding:   solana.EncodingBase64,
			Commitment: "confirmed",
			DataSlice:  nil,
		})
	if err != nil {
		return tokenMetadata, err
	}
	decoder := bin.NewBorshDecoder(resp.Value.Data.GetBinary())
	if err := tokenMetadata.UnmarshalWithDecoder(decoder); err != nil {
		return tokenMetadata, err
	}
	if _, err := json.MarshalIndent(tokenMetadata, "", "  "); err != nil {
		return tokenMetadata, err
	}
	tokenMetadata.Data.Name = strings.Trim(tokenMetadata.Data.Name, "\u0000")
	tokenMetadata.Data.Symbol = strings.Trim(tokenMetadata.Data.Symbol, "\u0000")
	tokenMetadata.Data.Uri = strings.Trim(tokenMetadata.Data.Uri, "\u0000")
	return tokenMetadata, nil
}

func (s impl) GetAccount(ctx context.Context, address string, v interface{}) error {
	resp, err := s.client.GetAccountInfoWithOpts(
		ctx,
		solana.MustPublicKeyFromBase58(address),
		&rpc.GetAccountInfoOpts{
			Encoding:   solana.EncodingBase64,
			Commitment: "confirmed",
			DataSlice:  nil,
		})
	if err != nil {
		logrus.
			WithError(err).
			WithField("address", address).
			Errorf("couldn't get acount info")
		return err
	}
	if err := bin.NewBinDecoder(resp.Value.Data.GetBinary()).Decode(v); err != nil {
		logrus.
			WithError(err).
			WithField("address", address).
			Errorf("failed to decode")
		return err
	}
	return nil
}

func (s impl) GetLargestTokenAccounts(ctx context.Context, mint string) ([]*rpc.TokenLargestAccountsResult, error) {
	pubkey, err := solana.PublicKeyFromBase58(mint)
	if err != nil {
		return nil, err
	}
	out, err := s.client.GetTokenLargestAccounts(ctx, pubkey, "confirmed")
	if out == nil {
		return nil, err
	}
	return out.Value, err
}

func (s impl) GetProgramAccounts(ctx context.Context, address string) ([]string, error) {
	offset := uint64(0)
	length := uint64(0)
	var res []string
	resp, err := s.client.GetProgramAccountsWithOpts(
		ctx,
		solana.MustPublicKeyFromBase58(address),
		&rpc.GetProgramAccountsOpts{
			Encoding:   solana.EncodingBase64,
			Commitment: "confirmed",
			DataSlice: &rpc.DataSlice{
				Offset: &offset,
				Length: &length,
			},
			Filters: []rpc.RPCFilter{
				{
					Memcmp: &rpc.RPCFilterMemcmp{
						Offset: 0,
					},
				},
			},
		})
	if err != nil {
		logrus.
			WithError(err).
			WithField("address", address).
			Errorf("couldn't get acount info")
		return nil, err
	}
	for i := range resp {
		res = append(res, resp[i].Pubkey.String())
	}
	return res, nil
}

func (s impl) GetUserBalances(ctx context.Context, wallet string) (*rpc.GetTokenAccountsResult, error) {
	return s.client.GetTokenAccountsByOwner(
		ctx,
		solana.MustPublicKeyFromBase58(wallet),
		&rpc.GetTokenAccountsConfig{
			ProgramId: &solana.TokenProgramID,
		},
		&rpc.GetTokenAccountsOpts{
			Commitment: rpc.CommitmentMax,
			Encoding:   solana.EncodingJSONParsed,
		})
}

func (s impl) MintToWallet(
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
	return s.signAndBroadcast(ctx, rpc.CommitmentFinalized, instructions...)
}

func (s impl) ProgramSubscribe(
	ctx context.Context, program string, onReceive func(string, []byte),
) error {
	url := getWSURL(s.network)
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
			// TODO(Mocha): This err block is super ugly
			if err != nil {
				logrus.
					WithError(err).
					WithFields(logrus.Fields{
						"event": program,
					}).
					Error("failed to get next msg from event ws")
				// TODO(Mocha): need to handle the case where this fails
				client, err = ws.Connect(ctx, url)
				if err != nil {
					logrus.
						WithError(err).
						WithFields(logrus.Fields{
							"event": program,
						}).
						Error("failed to get new ws client")
				}
				sub, err = client.ProgramSubscribeWithOpts(
					solana.MustPublicKeyFromBase58(program),
					rpc.CommitmentRecent,
					solana.EncodingBase64Zstd,
					nil,
				)
				if err != nil {
					logrus.
						WithError(err).
						WithFields(logrus.Fields{
							"event": program,
						}).
						Error("failed to get new program websocket subscription")
				}
				continue
			}
			if msg.Value.Account == nil || msg.Value.Account.Data == nil {
				logrus.
					WithFields(logrus.Fields{
						"event": program,
					}).
					Warning("event ws msg account or account data is nil")
				continue
			}
			decodedBinary := msg.Value.Account.Data.GetBinary()
			if decodedBinary == nil {
				logrus.
					WithFields(logrus.Fields{
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

func (s impl) GetWalletPubKey() solana.PublicKey {
	return s.wallet.PublicKey()
}

func (s impl) getWalletPrivKey() solana.PrivateKey {
	return s.wallet.PrivateKey
}

////////////////////////////////////////////////////////////
/// RPC Client Wrapper
////////////////////////////////////////////////////////////

func (s impl) GetTokenAccountBalance(
	ctx context.Context, destAccount solana.PublicKey, commitmentType rpc.CommitmentType,
) (*rpc.GetTokenAccountBalanceResult, error) {
	return s.client.GetTokenAccountBalance(ctx, destAccount, commitmentType)
}

func (s impl) GetAccountInfo(
	ctx context.Context, account solana.PublicKey,
) (*rpc.GetAccountInfoResult, error) {
	return s.client.GetAccountInfo(ctx, account)
}

func (s impl) GetVersion(ctx context.Context) (*rpc.GetVersionResult, error) {
	return s.client.GetVersion(ctx)
}

func (s impl) signAndBroadcast(
	ctx context.Context, commitment rpc.CommitmentType, instructions ...solana.Instruction,
) (string, error) {
	recent, err := s.client.GetRecentBlockhash(ctx, commitment)
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
		ctx, tx, false, commitment,
	)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction, err %s", err)
	}
	logFields["txHash"] = txHash

	return txHash.String(), nil
}

func GetURL(env configs.Network) string {
	switch env {
	case configs.DevnetNetwork:
		return rpc.DevNet_RPC
	case configs.MainnetNetwork:
		return rpc.MainNetBeta_RPC
	case configs.NilNetwork:
		fallthrough
	case configs.LocalNetwork:
		fallthrough
	default:
		return rpc.LocalNet_RPC
	}
}

func getWSURL(env configs.Network) string {
	switch env {
	case configs.DevnetNetwork:
		return rpc.DevNet_WS
	case configs.MainnetNetwork:
		return rpc.MainNetBeta_WS
	case configs.NilNetwork:
		fallthrough
	case configs.LocalNetwork:
		fallthrough
	default:
		return rpc.LocalNet_WS
	}
}