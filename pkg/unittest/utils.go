package unittest

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/dcaf-labs/drip/pkg/service/config"
	api "github.com/dcaf-labs/solana-go-retryable-http-client"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func GetTestPrivateKey() string {
	return "[95,189,40,215,74,154,138,123,245,115,184,90,2,187,104,25,241,164,79,247,14,69,207,235,40,245,13,157,149,20,13,227,252,155,201,43,89,96,76,119,162,241,148,53,80,193,126,159,80,213,140,166,144,139,205,143,160,238,11,34,192,249,59,31]"
}

func GetTestRequestRecorder(e *echo.Echo, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, "/", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func GetMockMainnetProductionConfig(ctrl *gomock.Controller) *config.MockAppConfig {
	mockConfig := config.NewMockAppConfig(ctrl)
	mockConfig.EXPECT().GetWalletPrivateKey().Return(GetTestPrivateKey()).AnyTimes()
	mockConfig.EXPECT().GetDiscordWebhookID().Return("").AnyTimes()
	mockConfig.EXPECT().GetDiscordWebhookAccessToken().Return("").AnyTimes()
	mockConfig.EXPECT().GetNetwork().Return(config.MainnetNetwork).AnyTimes()
	mockConfig.EXPECT().GetEnvironment().Return(config.ProductionEnv).AnyTimes()
	mockConfig.EXPECT().GetDripProgramID().Return("dripTrkvSyQKvkyWg7oi4jmeEGMA5scSYowHArJ9Vwk").AnyTimes()
	mockConfig.EXPECT().GetServerPort().Return(8080).AnyTimes()
	return mockConfig
}
func GetMockDevnetProductionConfig(ctrl *gomock.Controller) *config.MockAppConfig {
	mockConfig := config.NewMockAppConfig(ctrl)
	mockConfig.EXPECT().GetWalletPrivateKey().Return(GetTestPrivateKey()).AnyTimes()
	mockConfig.EXPECT().GetDiscordWebhookID().Return("").AnyTimes()
	mockConfig.EXPECT().GetDiscordWebhookAccessToken().Return("").AnyTimes()
	mockConfig.EXPECT().GetNetwork().Return(config.DevnetNetwork).AnyTimes()
	mockConfig.EXPECT().GetEnvironment().Return(config.ProductionEnv).AnyTimes()
	mockConfig.EXPECT().GetDripProgramID().Return("dripTrkvSyQKvkyWg7oi4jmeEGMA5scSYowHArJ9Vwk").AnyTimes()
	mockConfig.EXPECT().GetServerPort().Return(8080).AnyTimes()
	return mockConfig
}

func GetMockDevnetStagingConfig(ctrl *gomock.Controller) *config.MockAppConfig {
	mockConfig := config.NewMockAppConfig(ctrl)
	mockConfig.EXPECT().GetWalletPrivateKey().Return(GetTestPrivateKey()).AnyTimes()
	mockConfig.EXPECT().GetDiscordWebhookID().Return("").AnyTimes()
	mockConfig.EXPECT().GetDiscordWebhookAccessToken().Return("").AnyTimes()
	mockConfig.EXPECT().GetNetwork().Return(config.DevnetNetwork).AnyTimes()
	mockConfig.EXPECT().GetEnvironment().Return(config.StagingEnv).AnyTimes()
	mockConfig.EXPECT().GetDripProgramID().Return("F1NyoZsUhJzcpGyoEqpDNbUMKVvCnSXcCki1nN3ycAeo").AnyTimes()
	mockConfig.EXPECT().GetServerPort().Return(8080).AnyTimes()
	return mockConfig
}

func GetHTTPRecorderClientProvider(fixturePath string) (recorderProvider func() api.RetryableHTTPClientProvider, recorderTeardown func()) {
	r, err := recorder.New(fixturePath)
	if err != nil {
		logrus.WithError(err).Error("could not get recorder")
		os.Exit(1)
	}
	teardown := func() {
		if err := r.Stop(); err != nil {
			logrus.WithError(err).Error("could stop recorder")
			os.Exit(1)
		}
	}
	if r.Mode() != recorder.ModeRecordOnce {
		logrus.Error("recorder should be in ModeRecordOnce")
		os.Exit(1)
	}
	r.SetReplayableInteractions(true)
	r.SetMatcher(requestMatcher)
	recorderHTTPClient := r.GetDefaultClient()
	return func() api.RetryableHTTPClientProvider {
		return func(options api.RateLimitHTTPClientOptions) api.RetryableHTTPClient {
			options.HttpClient = *recorderHTTPClient
			options.CallsPerSecond = 10.0
			return api.GetDefaultClientProvider()(options)
		}
	}, teardown
}

func requestMatcher(r *http.Request, i cassette.Request) bool {
	if r.Body == nil || r.Body == http.NoBody {
		return cassette.DefaultMatcher(r, i)
	}

	var reqBody []byte
	var err error
	reqBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Errorf("failed to read request body")
		os.Exit(1)
	}
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	return r.Method == i.Method && r.URL.String() == i.URL && string(reqBody) == i.Body
}
