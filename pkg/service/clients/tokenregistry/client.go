package tokenregistry

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	api "github.com/dcaf-labs/solana-go-retryable-http-client"

	"github.com/samber/lo"

	"github.com/patrickmn/go-cache"
)

type TokenRegistry interface {
	GetTokenRegistry(ctx context.Context) (*TokenRegistryResponse, error)
	GetTokenRegistryToken(ctx context.Context, mint string) (*Token, error)
	GetTokenRegistryTokens(ctx context.Context, mints ...string) ([]*Token, error)
}

func NewTokenRegistry(retryClientProvider api.RetryableHTTPClientProvider) TokenRegistry {
	return newClient(retryClientProvider)
}

type client struct {
	api.RetryableHTTPClient
	cache *cache.Cache
}

func newClient(retryClientProvider api.RetryableHTTPClientProvider) *client {
	options := api.GetDefaultRateLimitHTTPClientOptions()
	options.CallsPerSecond = callsPerSecond
	retryClient := retryClientProvider(options)
	apiClient := client{retryClient, cache.New(60*time.Minute, 60*time.Minute)}
	return &apiClient
}

func (apiClient *client) GetTokenRegistry(ctx context.Context) (tokenRegistry *TokenRegistryResponse, err error) {
	if res, found := apiClient.cache.Get(url); found {
		return res.(*TokenRegistryResponse), nil
	}
	resp, err := apiClient.sendUnAuthenticatedGetRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &tokenRegistry); err != nil {
		return nil, err
	}
	apiClient.cache.Set(url, tokenRegistry, cache.DefaultExpiration)
	return tokenRegistry, nil
}

func (apiClient *client) GetTokenRegistryToken(ctx context.Context, mint string) (*Token, error) {
	tokenRegistry, err := apiClient.GetTokenRegistry(ctx)
	if err != nil {
		return nil, err
	}
	for i := range tokenRegistry.Tokens {
		if tokenRegistry.Tokens[i].Address == mint {
			return &tokenRegistry.Tokens[i], nil
		}
	}
	return nil, fmt.Errorf("mint not found in token registry")
}
func (apiClient *client) GetTokenRegistryTokens(ctx context.Context, mints ...string) ([]*Token, error) {
	tokenRegistry, err := apiClient.GetTokenRegistry(ctx)
	if err != nil {
		return nil, err
	}
	mintSet := lo.KeyBy[string, string](mints, func(mint string) string { return mint })
	return lo.FilterMap[Token, *Token](tokenRegistry.Tokens, func(token Token, _ int) (*Token, bool) {
		if _, ok := mintSet[token.Address]; ok {
			return &token, true
		}
		return nil, false
	}), nil
}

func (apiClient *client) sendUnAuthenticatedGetRequest(ctx context.Context, urlString string) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	return apiClient.Do(httpReq)
}
