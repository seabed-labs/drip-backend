package orcawhirlpool

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/clients/solana"

	"github.com/dcaf-labs/drip/pkg/service/configs"

	dripextension "github.com/dcaf-labs/drip-client/drip-extension-go"
	"github.com/dcaf-labs/drip/pkg/service/clients"
	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type OrcaWhirlpoolClient interface {
	GetOrcaWhirlpoolQuoteEstimate(ctx context.Context, whirlpool, inputTokenMint, inputAmount string) (*dripextension.V1OrcawhirlpoolQuote200Response, error)
}

func NewOrcaWhirlpoolClient(
	network configs.Network,
) OrcaWhirlpoolClient {
	return newClient(network)
}

type client struct {
	*dripextension.APIClient
	connectionUrl string
}

func newClient(network configs.Network) *client {
	rateLimiter := rate.NewLimiter(rate.Every(time.Second/callsPerSecond), 1)
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil
	httpClient.CheckRetry = clients.DefaultCheckRetry
	httpClient.RetryWaitMin = clients.DefaultRetryMin
	httpClient.RetryMax = clients.MaxRetries

	httpClient.RequestLogHook = func(logger retryablehttp.Logger, req *http.Request, retry int) {
		if err := rateLimiter.Wait(context.Background()); err != nil {
			log.WithError(err).Errorf("waiting for rate limit")
			return
		}
	}

	config := dripextension.NewConfiguration()
	config.HTTPClient = httpClient.StandardClient()
	config.Host = host
	config.UserAgent = "drip-backend"
	config.Scheme = "https"

	return &client{
		APIClient:     dripextension.NewAPIClient(config),
		connectionUrl: solana.GetURL(network, true),
	}
}

func (c *client) GetOrcaWhirlpoolQuoteEstimate(
	ctx context.Context,
	whirlpool, inputTokenMint, inputAmount string,
) (res *dripextension.V1OrcawhirlpoolQuote200Response, err error) {
	//

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
	}
	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("V1OrcawhirlpoolQuoteExecute returned non-200, statusCode: %d", httpRes.StatusCode)
	}
	return res, nil
}
