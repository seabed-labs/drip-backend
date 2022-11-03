package controller

import (
	"net/http"
	"testing"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/base"
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	"github.com/dcaf-labs/drip/pkg/unittest"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/test-go/testify/assert"
)

func TestHandler_GetToken(t *testing.T) {
	token := model.Token{
		Pubkey:      "So11111111111111111111111111111111111111112",
		Symbol:      utils.GetStringPtr("WSOL"),
		Decimals:    0,
		IconURL:     utils.GetStringPtr("url"),
		CoinGeckoID: utils.GetStringPtr("wrapped-solana"),
	}
	ctrl := gomock.NewController(t)
	mockConfig := config.NewMockAppConfig(ctrl)
	mockConfig.EXPECT().GetWalletPrivateKey().Return(unittest.GetTestPrivateKey()).AnyTimes()
	mockConfig.EXPECT().GetNetwork().Return(config.DevnetNetwork).AnyTimes()
	mockConfig.EXPECT().GetEnvironment().Return(config.StagingEnv).AnyTimes()
	mockConfig.EXPECT().GetServerPort().Return(8080).AnyTimes()
	e := echo.New()

	t.Run("should return a token with expected fields", func(t *testing.T) {
		c, rec := unittest.GetTestRequestRecorder(e, nil)
		m := repository.NewMockRepository(ctrl)
		h := NewHandler(mockConfig, solana.NewMockSolana(ctrl), base.NewMockBase(ctrl), m)
		m.
			EXPECT().
			GetTokenByAddress(gomock.Any(), "So11111111111111111111111111111111111111112").
			Return(&token, nil).
			AnyTimes()

		assert.NoError(t, h.GetV1TokenPubkeyPath(c, "So11111111111111111111111111111111111111112"))
		assert.Equal(t, rec.Code, http.StatusOK)
		token, err := apispec.ParseGetV1TokenPubkeyPathResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, token.JSON200)
		assert.Equal(t, token.JSON200.Pubkey, "So11111111111111111111111111111111111111112")
		assert.Equal(t, token.JSON200.Decimals, 0)
		assert.Equal(t, token.JSON200.IconUrl, utils.GetStringPtr("url"))
		assert.Equal(t, token.JSON200.Symbol, utils.GetStringPtr("WSOL"))
		assert.Equal(t, token.JSON200.CoinGeckoId, utils.GetStringPtr("wrapped-solana"))
	})
}
