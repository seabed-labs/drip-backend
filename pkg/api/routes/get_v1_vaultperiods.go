package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetV1Vaultperiods(c echo.Context, params apispec.GetV1VaultperiodsParams) error {
	vaultPeriodModels, err := h.repo.GetVaultPeriods(
		c.Request().Context(),
		(string)(params.Vault),
		(*string)(params.VaultPeriod),
		getPaginationParamsFromAPI(params.Offset, params.Limit),
	)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get vault periods")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}

	return c.JSON(http.StatusOK, vaultPeriodModelsToAPI(vaultPeriodModels))
}
