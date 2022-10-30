package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dcaf-labs/drip/pkg/service/utils"

	"github.com/dcaf-labs/drip/pkg/service/clients"
)

type CoinGeckoClient interface {
	GetCoinGeckoMetadata(ctx context.Context, contractAddress string) (cgMeta *CoinGeckoMetadataResponse, err error)
	sendUnAuthenticatedGetRequest(ctx context.Context, urlString string) (*http.Response, error)
}

func NewCoinGeckoClient(retryClientProvider clients.RetryableHTTPClientProvider) CoinGeckoClient {
	return newClient(retryClientProvider)
}

type client struct {
	clients.RetryableHTTPClient
}

func newClient(retryClientProvider clients.RetryableHTTPClientProvider) *client {
	retryClient := retryClientProvider(clients.RateLimitHTTPClientOptions{
		CallsPerSecond: utils.GetIntPtr(callsPerSecond),
	})
	apiClient := client{retryClient}
	return &apiClient
}

func (client *client) GetCoinGeckoMetadata(ctx context.Context, contractAddress string) (cgMeta *CoinGeckoMetadataResponse, err error) {

	urlString := fmt.Sprintf("%s/coins/solana/contract/%s", baseUrl, contractAddress)
	resp, err := client.sendUnAuthenticatedGetRequest(ctx, urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &cgMeta); err != nil {
		return nil, err
	}
	return cgMeta, nil
}

func (client *client) sendUnAuthenticatedGetRequest(ctx context.Context, urlString string) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(httpReq)
}
