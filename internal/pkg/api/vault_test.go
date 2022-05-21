package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dcaf-protocol/drip/internal/pkg/repository"

	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/internal/pkg/repository/model"
	Swagger "github.com/dcaf-protocol/drip/pkg/swagger"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/test-go/testify/assert"
)

func TestHandler_GetVaults(t *testing.T) {
	privKey := "[95,189,40,215,74,154,138,123,245,115,184,90,2,187,104,25,241,164,79,247,14,69,207,235,40,245,13,157,149,20,13,227,252,155,201,43,89,96,76,119,162,241,148,53,80,193,126,159,80,213,140,166,144,139,205,143,160,238,11,34,192,249,59,31]"
	tokenPairs := []*model.TokenPair{
		{
			ID:     "96b8b0ed-79a9-4972-bf5e-4ac8ab9e7fda",
			TokenA: "BfqATYbPZJFdEdYWkEbFRBnhv1LB6wtLn299HjMmE4uq",
			TokenB: "ASuqwxvC4FXxJGT9XqZMXbCKDQBaRTApEhN2oN3VL3A8",
		},
		{
			ID:     "96b8b0ed-79a9-4972-bf5e-4ac8ab9e7fda",
			TokenA: "ASuqwxvC4FXxJGT9XqZMXbCKDQBaRTApEhN2oN3VL3A8",
			TokenB: "BfqATYbPZJFdEdYWkEbFRBnhv1LB6wtLn299HjMmE4uq",
		},
	}
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
	e := echo.New()

	t.Run("should return an error when providing invalid amount", func(t *testing.T) {
		m := repository.NewMockRepository(ctrl)
		h := NewHandler(&configs.AppConfig{
			Environment: configs.DevnetEnv,
			Wallet:      privKey,
		}, solana.NewMockSolana(ctrl), m)

		params := Swagger.GetVaultsParams{
			TokenA:      nil,
			TokenB:      nil,
			ProtoConfig: nil,
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		m.
			EXPECT().
			GetVaults(gomock.Any(), nil, nil, nil).
			Return(nil, fmt.Errorf("some error")).
			AnyTimes()

		err := h.GetVaults(c, params)
		assert.NoError(t, err)
		assert.Equal(t, rec.Code, http.StatusInternalServerError)
		assert.Equal(t, "{\"error\":\"internal server error\"}\n", rec.Body.String())
	})

	t.Run("should return vaults without filter", func(t *testing.T) {
		m := repository.NewMockRepository(ctrl)
		h := NewHandler(&configs.AppConfig{
			Environment: configs.DevnetEnv,
			Wallet:      privKey,
		}, solana.NewMockSolana(ctrl), m)

		params := Swagger.GetVaultsParams{
			TokenA:      nil,
			TokenB:      nil,
			ProtoConfig: nil,
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		m.
			EXPECT().
			GetTokenPair(gomock.Any(), "96b8b0ed-79a9-4972-bf5e-4ac8ab9e7fda").
			Return(tokenPairs[0], nil).
			AnyTimes()
		m.
			EXPECT().
			GetTokenPair(gomock.Any(), "96b8b0ed-79a9-4972-bf5e-4ac8ab9e7fda").
			Return(tokenPairs[1], nil).
			AnyTimes()
		m.
			EXPECT().
			GetVaults(gomock.Any(), nil, nil, nil).
			Return(vaults, nil).
			AnyTimes()

		err := h.GetVaults(c, params)
		assert.NoError(t, err)
		assert.Equal(t, rec.Code, http.StatusOK)
		vaults, err := Swagger.ParseGetVaultsResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, vaults.JSON200)
		assert.Equal(t, len(*vaults.JSON200), 2)
	})
}
