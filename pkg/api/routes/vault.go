package controller

import (
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"

	Swagger "github.com/dcaf-labs/drip/pkg/apispec"
)

func (h Handler) GetVaults(c echo.Context, params Swagger.GetVaultsParams) error {
	res := Swagger.ListVaults{}

	vaultModels, err := h.repo.GetVaultsWithFilter(
		c.Request().Context(),
		(*string)(params.TokenA),
		(*string)(params.TokenB),
		(*string)(params.ProtoConfig),
	)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get vaults")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal api error"})
	}

	for i := range vaultModels {
		vault := vaultModels[i]
		// TODO(mocha): this can be done in the same query as the vault
		tokenPair, err := h.repo.GetTokenPairByID(c.Request().Context(), vault.TokenPairID)
		if err != nil {
			logrus.WithError(err).WithField("tokenPairID", tokenPair.ID).Errorf("could not find token pair")
			return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal api error"})
		}
		res = append(res,
			Swagger.Vault{
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
