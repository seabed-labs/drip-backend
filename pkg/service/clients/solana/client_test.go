package solana

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/configs"

	"github.com/dcaf-labs/solana-go-clients/pkg/tokenswap"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/test-go/testify/assert"
)

func TestSolanaClient(t *testing.T) {
	mint := "31nFDfb3b4qw8JPx4FaXGgEk8omt7NuHpPkwWCSym5rC"
	privKey := "[95,189,40,215,74,154,138,123,245,115,184,90,2,187,104,25,241,164,79,247,14,69,207,235,40,245,13,157,149,20,13,227,252,155,201,43,89,96,76,119,162,241,148,53,80,193,126,159,80,213,140,166,144,139,205,143,160,238,11,34,192,249,59,31]"
	config := configs.AppConfig{
		Network:     configs.DevnetNetwork,
		Environment: configs.StagingEnv,
		Wallet:      privKey,
	}
	client, err := NewSolanaClient(&config)
	assert.NoError(t, err)

	// Genesys go devnet RPC doesn't support airdrops for some reason
	rpcClient := rpc.New(rpc.DevNet_RPC)
	_, err = rpcClient.RequestAirdrop(
		context.Background(), client.GetWalletPubKey(), 100000000, "confirmed")
	assert.NoError(t, err)

	t.Run("GetWalletPubKey should return public key", func(t *testing.T) {
		assert.Equal(t, client.GetWalletPubKey().String(), "J15X2DWTRPVTJsofDrf5se4zkNv1sD1eJPgEHwvuNJer")
	})

	t.Run("getWalletPrivKey should return private key", func(t *testing.T) {
		assert.Equal(t, client.getWalletPrivKey().String(),
			"2v28DUBnjz9eGHbgMH4fVzixzpyP8SfdBVmo19vdhgzDddqnD4HMNiFgNKtQsKErEfhnRYKFY9k4WbaGyyFKQzai")
	})

	t.Run("should getWalletPrivKey", func(t *testing.T) {
		versionResponse, err := client.GetVersion(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, versionResponse.FeatureSet)
	})

	t.Run("mintToWallet should mint when wallet doesn't have a token account", func(t *testing.T) {
		destWallet := solana.NewWallet()
		txHash, err := client.MintToWallet(context.Background(), mint, destWallet.PublicKey().String(), 100)
		assert.NoError(t, err)
		assert.NotEmpty(t, txHash)
	})

	//t.Run("mintToWallet should mint when wallet has a token account", func(t *testing.T) {
	//	destWallet := solana.NewWallet()
	//	_, err = client.MintToWallet(context.Background(), mint, destWallet.PublicKey().String(), 100)
	//	assert.NoError(t, err)
	//	txHash, err := client.MintToWallet(context.Background(), mint, destWallet.PublicKey().String(), 100)
	//	assert.NoError(t, err)
	//	assert.NotEmpty(t, txHash)
	//})

	t.Run("getURL should return correct RPC url", func(t *testing.T) {
		assert.Equal(t, GetURL(configs.NilNetwork, true), rpc.LocalNet_RPC)
		assert.Equal(t, GetURL(configs.LocalNetwork, true), rpc.LocalNet_RPC)
		assert.Equal(t, GetURL(configs.DevnetNetwork, true), rpc.DevNet_RPC)
		//assert.Equal(t, GetURL(configs.MainnetNetwork, true), "https://dimensional-young-cloud.solana-mainnet.discover.quiknode.pro/a5a0fb3cfa38ab740ed634239fd502a99dbf028d/")
	})

	t.Run("ProgramSubscribe should subscribe to event", func(t *testing.T) {
		timeout := time.Second * 5
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err := client.ProgramSubscribe(ctx, solana.TokenProgramID.String(), func(address string, data []byte) {
			assert.NotEmpty(t, address)
			assert.NotEmpty(t, data)
		})
		assert.NoError(t, err)
		time.Sleep(timeout)
	})

	t.Run("GetProgramAccounts should return all account publickeys", func(t *testing.T) {
		timeout := time.Second * 5
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		offset, err := client.GetProgramAccounts(ctx, tokenswap.ProgramID.String())
		fmt.Println(len(offset))
		assert.NoError(t, err)
		assert.NotZero(t, offset)
		assert.True(t, len(offset) > 4500)
	})

	t.Run("GetTokenMetadataAccount should return token metadata for mint with valid token metadata account", func(t *testing.T) {
		tokenMetadata, err := client.GetTokenMetadataAccount(context.Background(), "2CsU92EN1AwEcJdWaVDa9n9o6AFB5nVBozEXZWDJcDj6")
		assert.NoError(t, err)
		assert.Equal(t, tokenMetadata.Data.Symbol, "DP")
		assert.Equal(t, tokenMetadata.Data.Name, "Drip Position")
	})

	t.Run("GetTokenMetadataAccount should return err for mint without token metadata account", func(t *testing.T) {
		_, err := client.GetTokenMetadataAccount(context.Background(), "2Kp1LB7Jo5RAt5gnisbtpu9tvK24WtV9zo9tZHKqFGqL")
		assert.Error(t, err)
	})

	t.Run("GetTokenMetadataAccount should return err for non-mint address", func(t *testing.T) {
		_, err := client.GetTokenMetadataAccount(context.Background(), "F1E3YejjuQ87d73enFG1w2toeeR4igZq8eQum1UgqGeF")
		assert.Error(t, err)
	})
}
