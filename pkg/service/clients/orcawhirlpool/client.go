package orcawhirlpool

import (
	"context"
	"fmt"
	"net/http"

	dripextension "github.com/dcaf-labs/drip-client/drip-extension-go"
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
	options := api.GetDefaultRateLimitHTTPClientOptions()
	options.CallsPerSecond = callsPerSecond
	retryClient := retryClientProvider(options)
	cfg := dripextension.NewConfiguration()
	cfg.HTTPClient = retryClient.StandardClient()
	cfg.Host = host
	cfg.UserAgent = "drip-backend"
	cfg.Scheme = "https"

	return &client{
		APIClient:     dripextension.NewAPIClient(cfg),
		connectionUrl: appConfig.GetSolanaRPCURL(),
	}
}

type client struct {
	*dripextension.APIClient
	connectionUrl string
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
