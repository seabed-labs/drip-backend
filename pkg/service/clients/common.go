package clients

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

const (
	defaultRetryMin = time.Second * 10
	maxRetries      = 3
)

type RetryableHTTPClient struct {
	*retryablehttp.Client
}

type RateLimitHTTPClientOptions struct {
	CallsPerSecond *int
	HttpClient     *http.Client
}

type RetryableHTTPClientProvider func(options RateLimitHTTPClientOptions) RetryableHTTPClient

func DefaultClientProvider() RetryableHTTPClientProvider {
	return func(options RateLimitHTTPClientOptions) RetryableHTTPClient {
		return getRateLimitedHTTPClient(options)
	}
}

func (r RetryableHTTPClient) Do(request *http.Request) (*http.Response, error) {
	retryableRequest, err := retryablehttp.FromRequest(request)
	if err != nil {
		return nil, err
	}
	return r.Client.Do(retryableRequest)
}

func (r RetryableHTTPClient) CloseIdleConnections() {
	r.HTTPClient.CloseIdleConnections()
}

func defaultCheckRetry(_ context.Context, resp *http.Response, err error) (bool, error) {
	log := logrus.NewEntry(logrus.New())
	if err != nil {
		log = log.WithError(err)
	}
	if resp != nil && resp.Request != nil && resp.Request.URL != nil {
		log = log.WithField("url", resp.Request.URL.String())
	}
	if resp != nil {
		log = log.WithField("statusCode", resp.StatusCode)
	}
	if resp != nil && resp.StatusCode >= http.StatusTooManyRequests {
		log.Error("waiting for rate limit")
		return true, err
	}
	return false, nil
}

func getRateLimitedHTTPClient(options RateLimitHTTPClientOptions) RetryableHTTPClient {
	callsPerSecond := 1
	if options.CallsPerSecond != nil {
		callsPerSecond = *options.CallsPerSecond
	}
	httpClient := retryablehttp.NewClient().HTTPClient
	if options.HttpClient != nil {
		httpClient = options.HttpClient
	}

	client := retryablehttp.NewClient()
	client.Logger = nil
	client.CheckRetry = defaultCheckRetry
	client.RetryWaitMin = defaultRetryMin
	client.RetryMax = maxRetries
	client.HTTPClient = httpClient

	rateLimiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(callsPerSecond)), 1)
	client.RequestLogHook = func(logger retryablehttp.Logger, req *http.Request, retry int) {
		if err := rateLimiter.Wait(context.Background()); err != nil {
			logrus.WithError(err).Errorf("failed to wait for rate limiter")
			return
		}
	}
	return RetryableHTTPClient{client}
}

func DecodeRequestBody[V any](resp *http.Response, res V) (V, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return res, err
	}
	return res, nil
}
