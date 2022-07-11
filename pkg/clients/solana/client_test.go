package solana

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dcaf-labs/drip/pkg/clients/solana/token_swap"
	"github.com/dcaf-labs/drip/pkg/configs"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/test-go/testify/assert"
)

func TestSolanaClient(t *testing.T) {
	mint := "31nFDfb3b4qw8JPx4FaXGgEk8omt7NuHpPkwWCSym5rC"
	privKey := "[95,189,40,215,74,154,138,123,245,115,184,90,2,187,104,25,241,164,79,247,14,69,207,235,40,245,13,157,149,20,13,227,252,155,201,43,89,96,76,119,162,241,148,53,80,193,126,159,80,213,140,166,144,139,205,143,160,238,11,34,192,249,59,31]"
	config := configs.AppConfig{
		Environment: configs.DevnetEnv,
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
		assert.Equal(t, getURL(configs.NilEnv), rpc.LocalNet_RPC)
		assert.Equal(t, getURL(configs.LocalnetEnv), rpc.LocalNet_RPC)
		assert.Equal(t, getURL(configs.DevnetEnv), "https://devnet.genesysgo.net")
		assert.Equal(t, getURL(configs.MainnetEnv), "https://ssc-dao.genesysgo.net")
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
		offset, err := client.GetProgramAccounts(ctx, token_swap.ProgramID.String())
		fmt.Println(len(offset))
		assert.NoError(t, err)
		assert.NotZero(t, offset)
		assert.True(t, len(offset) > 4500)
	})

	//t.Run("GetA should return token accounts", func(t *testing.T) {
	//	res, err := client.GetUserBalances(context.Background(), "DG43NUhq6gxUAcNGCJ45zCrBHgPDNNjMovocPHdHUqUD")
	//	assert.NoError(t, err)
	//	assert.NotNil(t, res)
	//	assert.NotNil(t, res.Value)
	//	assert.NotEmpty(t, res.Value)
	//	for _, a := range res.Value {
	//		rpc.mustJSONToInterface(mustAnyToJSON(out))
	//		//var tokenAccount token.Account
	//		//tokenAccount.
	//		//err := bin.NewBinDecoder(a.Account.Data.GetBinary()).Decode(&tokenAccount)
	//		//assert.NoError(t, err)
	//		//fmt.Println(tokenAccount)
	//	}
	//})
}
