package tokenregistry

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dcaf-labs/drip/pkg/service/utils"

	"github.com/gagliardetto/solana-go/rpc/jsonrpc"

	"github.com/dcaf-labs/drip/pkg/service/clients"
)

type TokenRegistry interface {
	GetTokenRegistry(ctx context.Context) (*TokenRegistryResponse, error)
	GetTokenRegistryToken(ctx context.Context, mint string) (*Token, error)
}

func NewTokenRegistry(retryClientProvider clients.RetryableHTTPClientProvider) TokenRegistry {
	return newClient(retryClientProvider)
}

type client struct {
	jsonrpc.HTTPClient
}

func newClient(retryClientProvider clients.RetryableHTTPClientProvider) *client {
	retryClient := retryClientProvider(clients.RateLimitHTTPClientOptions{
		CallsPerSecond: utils.GetIntPtr(callsPerSecond),
	})
	apiClient := client{retryClient}
	return &apiClient
}

func (apiClient *client) GetTokenRegistry(ctx context.Context) (tokenRegistry *TokenRegistryResponse, err error) {
	resp, err := apiClient.sendUnAuthenticatedGetRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &tokenRegistry); err != nil {
		return nil, err
	}
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

func (apiClient *client) sendUnAuthenticatedGetRequest(ctx context.Context, urlString string) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	return apiClient.Do(httpReq)
}
