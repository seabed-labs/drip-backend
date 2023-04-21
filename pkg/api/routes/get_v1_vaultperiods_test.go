package controller

import (
	"net/http"
	"testing"

	analyticsrepository "github.com/dcaf-labs/drip/pkg/service/repository/analytics"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/base"
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	"github.com/dcaf-labs/drip/pkg/unittest"
	solana2 "github.com/gagliardetto/solana-go"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/test-go/testify/assert"
	"gorm.io/gorm"
)

func TestHandler_GetV1Vaultperiods(t *testing.T) {
	vaultPeriods := []*model.VaultPeriod{
		{
			Pubkey:      solana2.NewWallet().PublicKey().String(),
			Vault:       "",
			PeriodID:    0,
			Twap:        decimal.Decimal{},
			Dar:         0,
			PriceBOverA: decimal.Decimal{},
		},
		{
			Pubkey:      solana2.NewWallet().PublicKey().String(),
			Vault:       "",
			PeriodID:    0,
			Twap:        decimal.Decimal{},
			Dar:         0,
			PriceBOverA: decimal.Decimal{},
		},
	}

	ctrl := gomock.NewController(t)
	mockConfig := unittest.GetMockDevnetStagingConfig(ctrl)
	e := echo.New()

	t.Run("should return internal server error if `GetVaultPeriods` returns an error", func(t *testing.T) {
		c, rec := unittest.GetTestRequestRecorder(e, nil)
		m := repository.NewMockRepository(ctrl)
		m2 := analyticsrepository.NewMockAnalyticsRepository(ctrl)
		h := NewHandler(mockConfig, solana.NewMockSolana(ctrl), base.NewMockBase(ctrl), m, m2)

		params := apispec.GetV1VaultperiodsParams{}

		m.
			EXPECT().
			GetVaultPeriods(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, gorm.ErrRecordNotFound).
			Times(1)

		assert.NoError(t, h.GetV1Vaultperiods(c, params))
		assert.Equal(t, rec.Code, http.StatusInternalServerError)
		res, err := apispec.ParseGetV1VaultperiodsResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, res.JSON500)
		assert.Equal(t, res.JSON500.Error, "internal server error")
	})

	t.Run("should call GetV1Vaultperiods with the correct parameters", func(t *testing.T) {
		c, rec := unittest.GetTestRequestRecorder(e, nil)
		m := repository.NewMockRepository(ctrl)
		m2 := analyticsrepository.NewMockAnalyticsRepository(ctrl)
		h := NewHandler(mockConfig, solana.NewMockSolana(ctrl), base.NewMockBase(ctrl), m, m2)

		limit := apispec.LimitQueryParam(5)
		offset := apispec.OffsetQueryParam(10)
		vaultPeriod := apispec.VaultPeriodQueryParam(solana2.NewWallet().PublicKey().String())
		params := apispec.GetV1VaultperiodsParams{
			Vault:       apispec.RequiredVaultQueryParam(solana2.NewWallet().PublicKey().String()),
			VaultPeriod: &vaultPeriod,
			Limit:       &limit,
			Offset:      &offset,
		}

		m.
			EXPECT().
			GetVaultPeriods(gomock.Any(), (string)(params.Vault), (*string)(params.VaultPeriod), repository.PaginationParams{
				Limit:  utils.GetIntPtr(5),
				Offset: utils.GetIntPtr(10),
			}).
			Return([]*model.VaultPeriod{}, nil).
			Times(1)

		assert.NoError(t, h.GetV1Vaultperiods(c, params))
		assert.Equal(t, rec.Code, http.StatusOK)
		res, err := apispec.ParseGetV1VaultperiodsResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, res.JSON200)
	})

	t.Run("should return 2 vaultPeriods", func(t *testing.T) {
		c, rec := unittest.GetTestRequestRecorder(e, nil)
		m := repository.NewMockRepository(ctrl)
		m2 := analyticsrepository.NewMockAnalyticsRepository(ctrl)
		h := NewHandler(mockConfig, solana.NewMockSolana(ctrl), base.NewMockBase(ctrl), m, m2)

		params := apispec.GetV1VaultperiodsParams{}

		m.
			EXPECT().
			GetVaultPeriods(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(vaultPeriods, nil).
			Times(1)

		assert.NoError(t, h.GetV1Vaultperiods(c, params))
		assert.Equal(t, rec.Code, http.StatusOK)
		res, err := apispec.ParseGetV1VaultperiodsResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, res.JSON200)
		assert.Len(t, *res.JSON200, 2)
	})
}
