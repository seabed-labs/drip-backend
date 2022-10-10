package clients

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/time/rate"

	log "github.com/sirupsen/logrus"
)

const (
	defaultRetryMin = time.Second * 10
	maxRetries      = 3
)

func defaultCheckRetry(_ context.Context, resp *http.Response, err error) (bool, error) {
	if resp != nil && resp.StatusCode >= http.StatusTooManyRequests {
		if err != nil {
			log.WithFields(
				log.Fields{
					"err":        err.Error(),
					"statusCode": resp.StatusCode,
				}).Errorf("waiting for rate limit")
		} else {
			log.WithFields(
				log.Fields{
					"statusCode": resp.StatusCode,
				}).Errorf("waiting for rate limit")
		}
		return true, err
	}
	return false, nil
}

type RetryableHTTPClient struct {
	*retryablehttp.Client
}

func (r RetryableHTTPClient) Do(request *http.Request) (*http.Response, error) {
	return r.HTTPClient.Do(request)
}

func (r RetryableHTTPClient) CloseIdleConnections() {
	r.HTTPClient.CloseIdleConnections()
}

func GetRateLimitedHTTPClient(callsPerSecond int) RetryableHTTPClient {
	rateLimiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(callsPerSecond)), 1)
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil
	httpClient.CheckRetry = defaultCheckRetry
	httpClient.RetryWaitMin = defaultRetryMin
	httpClient.RetryMax = maxRetries
	httpClient.RequestLogHook = func(logger retryablehttp.Logger, req *http.Request, retry int) {
		if err := rateLimiter.Wait(context.Background()); err != nil {
			log.WithError(err).Errorf("waiting for rate limit")
			return
		}
	}
	return RetryableHTTPClient{httpClient}
}
