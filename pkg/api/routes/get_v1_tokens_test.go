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

func TestHandler_GetTokens(t *testing.T) {
	tokenModels := []*model.Token{
		{
			Pubkey:      "So11111111111111111111111111111111111111112",
			Symbol:      utils.GetStringPtr("WSOL"),
			Decimals:    0,
			IconURL:     utils.GetStringPtr("url"),
			CoinGeckoID: utils.GetStringPtr("wrapped-solana"),
		},
		{
			Pubkey:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
			Symbol:      utils.GetStringPtr("USDC"),
			Decimals:    0,
			IconURL:     utils.GetStringPtr("url"),
			CoinGeckoID: utils.GetStringPtr("usd-coin"),
		},
	}
	ctrl := gomock.NewController(t)
	mockConfig := config.NewMockAppConfig(ctrl)
	mockConfig.EXPECT().GetWalletPrivateKey().Return(unittest.GetTestPrivateKey()).AnyTimes()
	mockConfig.EXPECT().GetNetwork().Return(config.DevnetNetwork).AnyTimes()
	mockConfig.EXPECT().GetEnvironment().Return(config.StagingEnv).AnyTimes()
	mockConfig.EXPECT().GetServerPort().Return(8080).AnyTimes()
	e := echo.New()

	t.Run("should return tokens with expected fields", func(t *testing.T) {
		c, rec := unittest.GetTestRequestRecorder(e, nil)
		m := repository.NewMockRepository(ctrl)
		h := NewHandler(mockConfig, solana.NewMockSolana(ctrl), base.NewMockBase(ctrl), m)
		m.
			EXPECT().
			GetAllSupportedTokens(gomock.Any()).
			Return(tokenModels, nil).
			AnyTimes()

		assert.NoError(t, h.GetV1Tokens(c))
		assert.Equal(t, rec.Code, http.StatusOK)
		tokens, err := apispec.ParseGetV1TokensResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, tokens.JSON200)
		assert.Equal(t, len(*tokens.JSON200), 2)
		assert.Equal(t, (*tokens.JSON200)[0].Pubkey, "So11111111111111111111111111111111111111112")
		assert.Equal(t, (*tokens.JSON200)[0].Decimals, 0)
		assert.Equal(t, (*tokens.JSON200)[0].IconUrl, utils.GetStringPtr("url"))
		assert.Equal(t, (*tokens.JSON200)[0].Symbol, utils.GetStringPtr("WSOL"))
		assert.Equal(t, (*tokens.JSON200)[0].CoinGeckoId, utils.GetStringPtr("wrapped-solana"))

		assert.Equal(t, (*tokens.JSON200)[1].Pubkey, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
		assert.Equal(t, (*tokens.JSON200)[1].Decimals, 0)
		assert.Equal(t, (*tokens.JSON200)[1].IconUrl, utils.GetStringPtr("url"))
		assert.Equal(t, (*tokens.JSON200)[1].Symbol, utils.GetStringPtr("USDC"))
		assert.Equal(t, (*tokens.JSON200)[1].CoinGeckoId, utils.GetStringPtr("usd-coin"))
	})
}
