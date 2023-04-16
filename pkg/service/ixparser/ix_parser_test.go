package ixparser

import (
	"testing"

	"github.com/test-go/testify/assert"

	"github.com/dcaf-labs/drip/pkg/unittest"

	api "github.com/dcaf-labs/solana-go-retryable-http-client"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
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
}
