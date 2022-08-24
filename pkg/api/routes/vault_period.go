package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/repository"

	apispec "github.com/dcaf-labs/drip/pkg/swagger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetVaultperiods(c echo.Context, params apispec.GetVaultperiodsParams) error {
	res := apispec.ListVaultPeriods{}
	limit := defaultLimit
	if params.Limit != nil {
		limit = int(*params.Limit)
	}
	var offset int
	if params.Offset != nil {
		offset = int(*params.Offset)
	}
	vaultPeriodModels, err := h.repo.GetVaultPeriods(
		c.Request().Context(),
		(string)(params.Vault),
		(*string)(params.VaultPeriod),
		repository.PaginationParams{
			Limit:  &limit,
			Offset: &offset,
		},
	)
	if err != nil {
		logrus.WithError(err).Errorf("failed to GetVaultPeriods")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	res = vaultPeriodDatabaseModelToAPIModel(vaultPeriodModels)
	return c.JSON(http.StatusOK, res)
}
