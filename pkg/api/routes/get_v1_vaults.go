package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetV1Vaults(c echo.Context, params apispec.GetV1VaultsParams) error {
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

	return c.JSON(http.StatusOK, vaultModelsToAPI(vaultModels))
}
