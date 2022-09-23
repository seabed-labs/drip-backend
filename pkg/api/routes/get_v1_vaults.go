package controller

import (
	"net/http"
	"strconv"

	"github.com/dcaf-labs/drip/pkg/api/apispec"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetV1Vaults(c echo.Context, params apispec.GetV1VaultsParams) error {
	res := apispec.ListVaults{}

	vaultModels, err := h.repo.GetVaultsWithFilter(
		c.Request().Context(),
		(*string)(params.TokenA),
		(*string)(params.TokenB),
		(*string)(params.ProtoConfig),
	)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get vaults")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}

	for i := range vaultModels {
		vault := vaultModels[i]
		// TODO(mocha): this can be done in the same query as the vault
		tokenPair, err := h.repo.GetTokenPairByID(c.Request().Context(), vault.TokenPairID)
		if err != nil {
			logrus.WithError(err).WithField("tokenPairID", tokenPair.ID).Errorf("could not find token pair")
			return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
		}
		res = append(res,
			apispec.Vault{
				DcaActivationTimestamp: strconv.FormatInt(vault.DcaActivationTimestamp.Unix(), 10),
				DripAmount:             strconv.FormatUint(vault.DripAmount, 10),
				LastDcaPeriod:          strconv.FormatUint(vault.LastDcaPeriod, 10),
				ProtoConfig:            vault.ProtoConfig,
				Pubkey:                 vault.Pubkey,
				TokenAAccount:          vault.TokenAAccount,
				TokenAMint:             tokenPair.TokenA,
				TokenBAccount:          vault.TokenBAccount,
				TokenBMint:             tokenPair.TokenB,
				TreasuryTokenBAccount:  vault.TreasuryTokenBAccount,
				Enabled:                vault.Enabled,
			},
		)
	}
	return c.JSON(http.StatusOK, res)
}
