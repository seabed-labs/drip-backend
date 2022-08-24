package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/repository"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"

	Swagger "github.com/dcaf-labs/drip/pkg/swagger"
)

func (h Handler) GetVaults(c echo.Context, params Swagger.GetVaultsParams) error {
	res := Swagger.ListVaults{}
	vaultEnabledFilter := true
	vaultModels, err := h.repo.GetVaultsWithFilter(
		c.Request().Context(),
		repository.VaultFilterParams{
			IsEnabled:   &vaultEnabledFilter,
			TokenA:      (*string)(params.TokenA),
			TokenB:      (*string)(params.TokenB),
			Vault:       nil,
			ProtoConfig: (*string)(params.ProtoConfig),
		},
	)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get vaults")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal api error"})
	}
	res = vaultWithTokenPairDatabaseModelToAPIModel(vaultModels)
	return c.JSON(http.StatusOK, res)
}
