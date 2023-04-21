package controller

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	analyticsrepository "github.com/dcaf-labs/drip/pkg/service/repository/analytics"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/base"
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/unittest"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/test-go/testify/assert"
)

func TestHandler_GetVaults(t *testing.T) {
	vaults := []*model.Vault{
		{
			Pubkey:                 "3iz6nZVjiGZtdEffAUDrVh4A5BnwN6ZoHj3nPPZtKJfV",
			ProtoConfig:            "mRcJ27ztTCFntbUvv7V2PSxqL9fJfg1KH4fzZSYVP4L",
			TokenAAccount:          "6PmcdLzbELLxaPc3Fq6FjiSj7wtjA4MEt1UCZBnHh6tw",
			TokenBAccount:          "5q7HLgfvxmkqAK6QaEFYrNmvKvzQZjWJzjwRu4toi9Sw",
			TreasuryTokenBAccount:  "CrVdqMmYCbBs8zG2rmwdWgmsSArKTbMUv3qvTz8J6YWC",
			TokenPairID:            "96b8b0ed-79a9-4972-bf5e-4ac8ab9e7fda",
			LastDcaPeriod:          10,
			DripAmount:             1000,
			DcaActivationTimestamp: time.Unix(1652748331, 0),
		},
		{
			Pubkey:                 "7rERMPGFMFi5k7jsQj3yXbW9uPDQj8FBx6vjwVXpknbh",
			ProtoConfig:            "mRcJ27ztTCFntbUvv7V2PSxqL9fJfg1KH4fzZSYVP4L",
			TokenAAccount:          "5y95dsjKJPaf94Kv8K6NbhyDYWswZycfcHNbwXMo6Xdk",
			TokenBAccount:          "9B6zNC4ijKfgGjReicqzv5dEPhJMt2oDH97b4AjUAs9a",
			TreasuryTokenBAccount:  "ugY3pNYSKmMo4msf9VRVnf7SxFQXDqhNcjbNMKSW9gL",
			TokenPairID:            "96b8b0ed-79a9-4972-bf5e-4ac8ab9e7fda",
			LastDcaPeriod:          99,
			DripAmount:             50,
			DcaActivationTimestamp: time.Unix(1652748331, 0),
		},
	}
	ctrl := gomock.NewController(t)
	mockConfig := unittest.GetMockDevnetStagingConfig(ctrl)
	e := echo.New()

	t.Run("should return internal server error if `GetVaultsWithFilter` returns an error", func(t *testing.T) {
		c, rec := unittest.GetTestRequestRecorder(e, nil)
		m := repository.NewMockRepository(ctrl)
		m2 := analyticsrepository.NewMockAnalyticsRepository(ctrl)
		h := NewHandler(mockConfig, solana.NewMockSolana(ctrl), base.NewMockBase(ctrl), m, m2)

		params := apispec.GetV1VaultsParams{
			TokenA:      nil,
			TokenB:      nil,
			ProtoConfig: nil,
		}

		m.
			EXPECT().
			GetVaultsWithFilter(gomock.Any(), nil, nil, nil).
			Return(nil, fmt.Errorf("some error")).
			AnyTimes()

		err := h.GetV1Vaults(c, params)
		assert.NoError(t, err)
		assert.Equal(t, rec.Code, http.StatusInternalServerError)
		assert.Equal(t, "{\"error\":\"internal api error\"}\n", rec.Body.String())
	})

	t.Run("should return vaults without filter", func(t *testing.T) {
		c, rec := unittest.GetTestRequestRecorder(e, nil)
		m := repository.NewMockRepository(ctrl)
		m2 := analyticsrepository.NewMockAnalyticsRepository(ctrl)
		h := NewHandler(mockConfig, solana.NewMockSolana(ctrl), base.NewMockBase(ctrl), m, m2)

		m.
			EXPECT().
			GetVaultsWithFilter(gomock.Any(), nil, nil, nil).
			Return(vaults, nil).
			AnyTimes()

		params := apispec.GetV1VaultsParams{
			TokenA:      nil,
			TokenB:      nil,
			ProtoConfig: nil,
		}

		err := h.GetV1Vaults(c, params)
		assert.NoError(t, err)
		assert.Equal(t, rec.Code, http.StatusOK)
		vaults, err := apispec.ParseGetV1VaultsResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, vaults.JSON200)
		assert.Equal(t, len(*vaults.JSON200), 2)
	})
}
