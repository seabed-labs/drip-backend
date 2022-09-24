package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetV1DripSpltokenswapconfigs(c echo.Context, params apispec.GetV1DripSpltokenswapconfigsParams) error {
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
	tokenSwapsByVault, err := h.base.GetBestTokenSwapsForVaults(c.Request().Context(), vaults)
	if err != nil {
		logrus.WithError(err).WithField("vault", *params.Vault).Errorf("failed to GetBestTokenSwapsForVaults")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "Internal Server Error"})
	}

	res := apispec.ListSplTokenSwapConfigs{}
	vaultsByPubkey := model.GetVaultsByPubkey(vaults)
	for vaultPubkey, tokenSwap := range tokenSwapsByVault {
		vault, ok := vaultsByPubkey[vaultPubkey]
		if !ok {
			logrus.
				WithError(err).
				WithField("vaultPubkey", vaultPubkey).
				Errorf("invalid vaultPubkey in tokenSwapsByVault")
			return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "Internal Server Error"})
		}
		res = append(res, vaultTokenSwapToAPI(vault, tokenSwap))
	}
	return c.JSON(http.StatusOK, res)
}
