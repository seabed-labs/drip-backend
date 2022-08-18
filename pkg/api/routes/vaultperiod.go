package controller

import (
	"net/http"
	"strconv"

	"github.com/dcaf-labs/drip/pkg/repository"

	Swagger "github.com/dcaf-labs/drip/pkg/swagger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetVaultperiods(c echo.Context, params Swagger.GetVaultperiodsParams) error {
	var res Swagger.ListVaultPeriods
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
		logrus.WithError(err).Errorf("failed to get vault periods")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal server error"})
	}

	for i := range vaultPeriodModels {
		vaultPeriod := vaultPeriodModels[i]
		res = append(res, struct {
			Dar      string `json:"dar"`
			PeriodId string `json:"periodId"`
			Pubkey   string `json:"pubkey"`
			Twap     string `json:"twap"`
			Vault    string `json:"vault"`
		}{
			Pubkey:   vaultPeriod.Pubkey,
			Vault:    vaultPeriod.Vault,
			PeriodId: strconv.FormatUint(vaultPeriod.PeriodID, 10),
			Twap:     vaultPeriod.Twap.String(),
			Dar:      strconv.FormatUint(vaultPeriod.Dar, 10),
		},
		)
	}
	return c.JSON(http.StatusOK, res)
}
