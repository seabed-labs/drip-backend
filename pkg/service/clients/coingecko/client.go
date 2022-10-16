package coinGecko

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dcaf-labs/drip/pkg/service/clients"
)

type CoinGeckoClient interface {
	GetCoinGeckoMetadata(ctx context.Context, contractAddress string) (cgMeta *CoinGeckoMetadataResponse, err error)
	sendUnAuthenticatedGetRequest(ctx context.Context, urlString string) (*http.Response, error)
}

func NewCoinGeckoClient() CoinGeckoClient {
	return newClient()
}

type client struct {
	clients.RetryableHTTPClient
}

func newClient() *client {
	httpClient := clients.GetRateLimitedHTTPClient(callsPerSecond)
	apiClient := client{httpClient}
	return &apiClient
}

func (client *client) GetCoinGeckoMetadata(ctx context.Context, contractAddress string) (cgMeta *CoinGeckoMetadataResponse, err error) {

	urlString := baseUrl + "/coins/solana/contract/" + contractAddress
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
