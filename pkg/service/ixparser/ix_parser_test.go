package ixparser

import (
	"context"
	"testing"

	"github.com/samber/lo"

	"github.com/dcaf-labs/drip/pkg/unittest"
	drip "github.com/dcaf-labs/solana-drip-go/pkg"
	api "github.com/dcaf-labs/solana-go-retryable-http-client"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"github.com/test-go/testify/assert"
)

func TestIxParser(t *testing.T) {

	url := "https://palpable-warmhearted-hexagon.solana-mainnet.discover.quiknode.pro/5793cf44e6e16325347e62d571454890f16e0388"
	recorderProvider, recorderTeardown := unittest.GetHTTPRecorderClientProvider("./fixtures/drip_1")
	defer recorderTeardown()
	options := api.GetDefaultRateLimitHTTPClientOptions()
	opts := &jsonrpc.RPCClientOpts{
		HTTPClient: recorderProvider()(options),
	}
	rpcClient := rpc.NewWithCustomRPCClient(jsonrpc.NewClientWithOpts(url, opts))
	assert.NotNil(t, rpcClient)

	parser := ixParser{}

	t.Run("can parse v1 dripOrcaWhirlpool ix", func(t *testing.T) {
		txRaw, err := rpcClient.GetTransaction(context.Background(), solana.MustSignatureFromBase58("5uECxzjML1a5sXuMPqwcpP8BMAKCXCpUVwFydPHWBnaW9V72ockWQKevfSQBaqiSxtGKtmstkVaxmo5Jcfnp9Lcb"), nil)
		assert.NoError(t, err)
		assert.NotNil(t, txRaw)
		assert.NotNil(t, txRaw.Transaction)
		assert.Equal(t, drip.GetIdlVersion(txRaw.Slot), 1)

		tx, err := txRaw.Transaction.GetTransaction()
		assert.NoError(t, err)
		assert.Len(t, tx.Message.Instructions, 2)

		ix := tx.Message.Instructions[1]
		accounts, err := ix.ResolveInstructionAccounts(&tx.Message)
		assert.NoError(t, err)
		parsedIx, ixName, err := parser.MaybeParseV1DripOrcaWhirlpool(accounts, ix.Data)
		assert.NoError(t, err)
		assert.NotNil(t, parsedIx)
		assert.Equal(t, ixName, "DripOrcaWhirlpool")
		assert.Equal(t, parsedIx.GetDripOrcaWhirlpoolAccounts().Common.Vault.String(), "6PnzoW2Bcs6WGqYvecfSxN9C2EeDmQCjUCeFA7JDDfZG")
		assert.Equal(t, parsedIx.GetDripOrcaWhirlpoolAccounts().Whirlpool.String(), "HJPjoWUrhoZzkNfRpHuieeFk9WcZWjwy6PBjZ81ngndJ")
	})

	t.Run("can parse token transfer ix", func(t *testing.T) {
		txRaw, err := rpcClient.GetTransaction(context.Background(), solana.MustSignatureFromBase58("5uECxzjML1a5sXuMPqwcpP8BMAKCXCpUVwFydPHWBnaW9V72ockWQKevfSQBaqiSxtGKtmstkVaxmo5Jcfnp9Lcb"), nil)
		assert.NoError(t, err)
		assert.NotNil(t, txRaw)
		assert.NotNil(t, txRaw.Transaction)
		assert.Equal(t, drip.GetIdlVersion(txRaw.Slot), 1)

		tx, err := txRaw.Transaction.GetTransaction()
		assert.NoError(t, err)
		assert.Len(t, tx.Message.Instructions, 2)

		innerInstructions := lo.Filter(txRaw.Meta.InnerInstructions, func(innerIx rpc.InnerInstruction, _ int) bool {
			return innerIx.Index == 1
		})
		innerIx := innerInstructions[0].Instructions[0]
		accounts, err := innerIx.ResolveInstructionAccounts(&tx.Message)
		assert.NoError(t, err)
		parsedIx, ixName, err := parser.MaybeParseTokenTransfer(accounts, innerIx.Data)
		assert.NoError(t, err)
		assert.NotNil(t, parsedIx)
		assert.Equal(t, ixName, "Transfer")
		assert.Equal(t, *parsedIx.Amount, uint64(0xe34))
	})
}
