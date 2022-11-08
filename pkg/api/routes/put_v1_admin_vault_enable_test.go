package controller

import (
	"net/http"
	"testing"

	"gorm.io/gorm"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/base"
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/unittest"
	solana2 "github.com/gagliardetto/solana-go"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/test-go/testify/assert"
)

func TestHandler_PutV1AdminVaultPubkeyPathEnable(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockConfig := unittest.GetMockDevnetStagingConfig(ctrl)
	e := echo.New()

	t.Run("should return updated vault if repository methods do not throw an error", func(t *testing.T) {
		c, rec := unittest.GetTestRequestRecorder(e, nil)

		vaultPubkey := solana2.NewWallet().PublicKey().String()

		m := repository.NewMockRepository(ctrl)
		h := NewHandler(mockConfig, solana.NewMockSolana(ctrl), base.NewMockBase(ctrl), m)

		m.
			EXPECT().
			AdminGetVaultByAddress(gomock.Any(), vaultPubkey).
			Return(nil, gorm.ErrRecordNotFound).
			Times(1)

		m.
			EXPECT().
			AdminSetVaultEnabled(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil).
			Times(0)

		params := apispec.PutV1AdminVaultPubkeyPathEnableParams{}
		err := h.PutV1AdminVaultPubkeyPathEnable(c, apispec.PubkeyPathParam(vaultPubkey), params)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		updatedVault, err := apispec.ParsePutV1AdminVaultPubkeyPathEnableResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, updatedVault.JSON400)
		assert.Equal(t, updatedVault.JSON400.Error, gorm.ErrRecordNotFound.Error())
	})

	t.Run("should return updated vault if repository methods do not throw an error", func(t *testing.T) {
		c, rec := unittest.GetTestRequestRecorder(e, nil)

		mockDisabledVault := model.Vault{
			Pubkey:        solana2.NewWallet().PublicKey().String(),
			LastDcaPeriod: 0,
			Enabled:       false,
		}
		mockEnabledVault := model.Vault{
			Pubkey:        mockDisabledVault.Pubkey,
			LastDcaPeriod: 0,
			Enabled:       false,
		}

		m := repository.NewMockRepository(ctrl)
		h := NewHandler(mockConfig, solana.NewMockSolana(ctrl), base.NewMockBase(ctrl), m)

		m.
			EXPECT().
			AdminGetVaultByAddress(gomock.Any(), mockDisabledVault.Pubkey).
			Return(&mockDisabledVault, nil).
			Times(1)

		m.
			EXPECT().
			AdminSetVaultEnabled(gomock.Any(), mockDisabledVault.Pubkey, !mockDisabledVault.Enabled).
			Return(&mockEnabledVault, nil).
			Times(1)

		params := apispec.PutV1AdminVaultPubkeyPathEnableParams{}
		err := h.PutV1AdminVaultPubkeyPathEnable(c, apispec.PubkeyPathParam(mockDisabledVault.Pubkey), params)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		updatedVault, err := apispec.ParsePutV1AdminVaultPubkeyPathEnableResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, updatedVault.JSON200)
		assert.Equal(t, updatedVault.JSON200.Pubkey, mockEnabledVault.Pubkey)
		assert.Equal(t, updatedVault.JSON200.Enabled, mockEnabledVault.Enabled)
	})
}
