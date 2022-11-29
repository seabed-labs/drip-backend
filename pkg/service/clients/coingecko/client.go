package coingecko

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/dcaf-labs/drip/pkg/service/clients"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	"github.com/samber/lo"
)

type CoinGeckoClient interface {
	GetSolanaCoinsList(ctx context.Context) (CoinsListResponse, error)
	GetCoinGeckoMetadata(ctx context.Context, contractAddress string) (cgMeta CoinGeckoMetadataResponse, err error)
	GetMarketPriceForTokens(ctx context.Context, coinGeckoIDs ...string) (CoinGeckoTokensMarketPriceResponse, error)
	sendUnAuthenticatedGetRequest(ctx context.Context, urlString string) (*http.Response, error)
}

func NewCoinGeckoClient(retryClientProvider clients.RetryableHTTPClientProvider) CoinGeckoClient {
	return newClient(retryClientProvider)
}

type client struct {
	clients.RetryableHTTPClient
	cache *cache.Cache
}

func newClient(retryClientProvider clients.RetryableHTTPClientProvider) *client {
	retryClient := retryClientProvider(clients.RateLimitHTTPClientOptions{
		CallsPerSecond: utils.GetIntPtr(callsPerSecond),
	})
	apiClient := client{retryClient, cache.New(60*time.Minute, 60*time.Minute)}
	return &apiClient
}

func (client *client) GetSolanaCoinsList(ctx context.Context) (CoinsListResponse, error) {
	if res, found := client.cache.Get(cacheCoinsListPath); found {
		return res.([]CoinResponse), nil
	}
	urlString := fmt.Sprintf("%s%s?include_platform=true", baseUrl, coinsListPath)
	resp, err := client.sendUnAuthenticatedGetRequest(ctx, urlString)
	if err != nil {
		return nil, err
	}
	res, err := clients.DecodeRequestBody(resp, CoinsListResponse{})
	if err != nil {
		return nil, err
	}
	filteredRes := lo.Filter[CoinResponse](res, func(coin CoinResponse, _ int) bool {
		return coin.Platforms.Solana != nil && *coin.Platforms.Solana != ""
	})
	client.cache.Set(cacheCoinsListPath, filteredRes, cache.DefaultExpiration)
	return filteredRes, nil
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
	return clients.DecodeRequestBody(resp, CoinGeckoTokensMarketPriceResponse{})
}

func (client *client) GetCoinGeckoMetadata(ctx context.Context, contractAddress string) (CoinGeckoMetadataResponse, error) {
	cacheKey := cacheSolanaContractPath + contractAddress
	if res, found := client.cache.Get(cacheKey); found {
		return res.(CoinGeckoMetadataResponse), nil
	}
	urlString := fmt.Sprintf("%s%s/%s", baseUrl, solanaContractPath, contractAddress)
	resp, err := client.sendUnAuthenticatedGetRequest(ctx, urlString)
	if err != nil {
		return CoinGeckoMetadataResponse{}, err
	}
	res, err := clients.DecodeRequestBody(resp, CoinGeckoMetadataResponse{})
	if err != nil {
		return CoinGeckoMetadataResponse{}, err
	}
	client.cache.Set(cacheKey, res, cache.DefaultExpiration)
	return res, nil
}

func (client *client) sendUnAuthenticatedGetRequest(ctx context.Context, urlString string) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(httpReq)
}
