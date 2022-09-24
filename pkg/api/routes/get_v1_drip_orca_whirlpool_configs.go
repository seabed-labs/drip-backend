package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetV1DripOrcawhirlpoolconfigs(c echo.Context, params apispec.GetV1DripOrcawhirlpoolconfigsParams) error {
	res := apispec.ListOrcaWhirlpoolConfigs{}

	var vaults []*model.Vault
	if params.Vault != nil {
		vault, err := h.repo.GetVaultByAddress(c.Request().Context(), string(*params.Vault))
		if err != nil {
			logrus.WithError(err).WithField("vault", *params.Vault).Errorf("failed to get vault by address")
			return c.JSON(http.StatusBadRequest, apispec.ErrorResponse{Error: "invalid vault address"})
		}
		vaults = []*model.Vault{vault}
	} else {
		var err error
		vaults, err = h.repo.GetVaultsWithFilter(c.Request().Context(), nil, nil, nil)
		if err != nil {
			logrus.WithError(err).Errorf("failed to get vaults")
			return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "failed to get vaults"})
		}
	}

	whirlpoolsByVault, err := h.base.GetBestOrcaWhirlpoolForVaults(c.Request().Context(), vaults)
	if err != nil {
		logrus.WithError(err).Errorf("failed to GetBestOrcaWhirlpoolForVaults")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "Internal Server Error"})
	}
	vaultsByPubkey := model.GetVaultsByPubkey(vaults)
	for vaultPubkey, orcaWhirlpool := range whirlpoolsByVault {
		vault, ok := vaultsByPubkey[vaultPubkey]
		if !ok {
			logrus.
				WithError(err).
				WithField("vaultPubkey", vaultPubkey).
				Errorf("invalid vaultPubkey in whirlpoolsByVault")
			return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "Internal Server Error"})
		}
		res = append(res, vaultWhirlpoolToAPI(vault, orcaWhirlpool))
	}
	return c.JSON(http.StatusOK, res)
}
