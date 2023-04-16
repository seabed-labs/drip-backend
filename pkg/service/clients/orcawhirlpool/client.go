package orcawhirlpool

import (
	"context"
	"fmt"
	"net/http"

	dripextension "github.com/dcaf-labs/drip-client/drip-extension-go"
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/config"
	api "github.com/dcaf-labs/solana-go-retryable-http-client"
)

type OrcaWhirlpoolClient interface {
	GetOrcaWhirlpoolQuoteEstimate(ctx context.Context, whirlpool, inputTokenMint, inputAmount string) (*dripextension.V1OrcawhirlpoolQuote200Response, error)
}

func NewOrcaWhirlpoolClient(
	appConfig config.AppConfig,
	retryClientProvider api.RetryableHTTPClientProvider,
) OrcaWhirlpoolClient {
	return newClient(appConfig.GetNetwork(), retryClientProvider)
}

type client struct {
	*dripextension.APIClient
	connectionUrl string
}

func newClient(network config.Network, retryClientProvider api.RetryableHTTPClientProvider) *client {
	options := api.GetDefaultRateLimitHTTPClientOptions()
	options.CallsPerSecond = callsPerSecond
	retryClient := retryClientProvider(options)
	config := dripextension.NewConfiguration()
	config.HTTPClient = retryClient.StandardClient()
	config.Host = host
	config.UserAgent = "drip-backend"
	config.Scheme = "https"

	connectionUrl, _ := solana.GetURLWithRateLimit(network)
	return &client{
		APIClient:     dripextension.NewAPIClient(config),
		connectionUrl: connectionUrl,
	}
}

func (c *client) GetOrcaWhirlpoolQuoteEstimate(
	ctx context.Context,
	whirlpool, inputTokenMint, inputAmount string,
) (res *dripextension.V1OrcawhirlpoolQuote200Response, err error) {
	res, httpRes, err := c.DefaultApi.
		V1OrcawhirlpoolQuote(ctx).
		V1OrcawhirlpoolQuoteRequest(dripextension.V1OrcawhirlpoolQuoteRequest{
			ConnectionUrl:  c.connectionUrl,
			Whirlpool:      whirlpool,
			InputTokenMint: inputTokenMint,
			InputAmount:    inputAmount,
		}).
		Execute()
	if err != nil {
		return nil, err
	} else if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("V1OrcawhirlpoolQuoteExecute returned non-200, statusCode: %d", httpRes.StatusCode)
	} else if res == nil {
		return nil, fmt.Errorf("nil V1OrcawhirlpoolQuote200Response with 200 status")
	}
	return res, nil
}
