package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	lo2 "github.com/samber/lo"

	"github.com/dcaf-labs/drip/pkg/service/clients"
	"github.com/dcaf-labs/drip/pkg/service/utils"
)

type CoinGeckoClient interface {
	GetSolanaCoinsList(ctx context.Context) (CoinsListResponse, error)
	GetCoinGeckoMetadata(ctx context.Context, contractAddress string) (cgMeta *CoinGeckoMetadataResponse, err error)
	GetMarketPriceForTokens(ctx context.Context, coinGeckoIDs ...string) (CoinGeckoTokensMarketPriceResponse, error)
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

func (client *client) GetSolanaCoinsList(ctx context.Context) (CoinsListResponse, error) {
	urlString := fmt.Sprintf("%s%s?include_platform=true", baseUrl, coinsList)
	resp, err := client.sendUnAuthenticatedGetRequest(ctx, urlString)
	if err != nil {
		return nil, err
	}
	var res CoinsListResponse
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return lo2.Filter[CoinResponse](res, func(coin CoinResponse, _ int) bool {
		return coin.Platforms.Solana != nil && *coin.Platforms.Solana != ""
	}), nil
}

func (client *client) GetMarketPriceForTokens(ctx context.Context, coinGeckoIDs ...string) (CoinGeckoTokensMarketPriceResponse, error) {
	if len(coinGeckoIDs) == 0 {
		return nil, fmt.Errorf("missing coinGeckoIDs")
	}
	urlString := fmt.Sprintf("%s%s?vs_currency=usd&ids=%s&order=market_cap_desc&per_page=%d&page=1&sparkline=false",
		baseUrl, coinsMarketsPath,
		url.QueryEscape(strings.Join(coinGeckoIDs, ",")), CoinsMarketsPathLimit)
	resp, err := client.sendUnAuthenticatedGetRequest(ctx, urlString)
	if err != nil {
		return nil, err
	}
	var res CoinGeckoTokensMarketPriceResponse
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (client *client) GetCoinGeckoMetadata(ctx context.Context, contractAddress string) (*CoinGeckoMetadataResponse, error) {
	urlString := fmt.Sprintf("%s/coins/solana/contract/%s", baseUrl, contractAddress)
	resp, err := client.sendUnAuthenticatedGetRequest(ctx, urlString)
	if err != nil {
		return nil, err
	}
	var res *CoinGeckoMetadataResponse
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (client *client) sendUnAuthenticatedGetRequest(ctx context.Context, urlString string) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(httpReq)
}
