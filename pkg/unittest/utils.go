package unittest

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/golang/mock/gomock"

	"github.com/labstack/echo/v4"
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
