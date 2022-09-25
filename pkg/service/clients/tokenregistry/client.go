package tokenregistry

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/clients"
	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type TokenRegistry interface {
	GetTokenRegistry(ctx context.Context) (*TokenRegistryResponse, error)
	GetTokenRegistryToken(ctx context.Context, mint string) (*Token, error)
}

func NewTokenRegistry() TokenRegistry {
	return newClient()
}

type client struct {
	*retryablehttp.Client
	*rate.Limiter
}

func newClient() *client {
	rateLimiter := rate.NewLimiter(rate.Every(time.Second/callsPerSecond), 1)
	httpClient := retryablehttp.NewClient()
	httpClient.CheckRetry = clients.DefaultCheckRetry
	httpClient.RetryWaitMin = clients.DefaultRetryMin
	httpClient.RetryMax = clients.MaxRetries
	apiClient := client{
		Client:  httpClient,
		Limiter: rateLimiter,
	}
	apiClient.RequestLogHook = func(logger retryablehttp.Logger, req *http.Request, retry int) {
		if err := apiClient.Limiter.Wait(context.Background()); err != nil {
			log.WithError(err).Errorf("waiting for rate limit")
			return
		}
	}
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
	retryableRequest, err := retryablehttp.FromRequest(httpReq)
	if err != nil {
		return nil, err
	}
	return apiClient.Do(retryableRequest)
}
