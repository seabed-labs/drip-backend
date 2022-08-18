package controller

import (
	"net/http"
	"strconv"

	"github.com/dcaf-labs/drip/pkg/repository"

	apispec "github.com/dcaf-labs/drip/pkg/swagger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetV1AdminPositions(c echo.Context, params apispec.GetV1AdminPositionsParams) error {
	var res apispec.ListAdminPositions

	positions, err := h.repo.GetAdminPositions(
		c.Request().Context(),
		(*bool)(params.Enabled),
		repository.PositionFilterParams{
			IsClosed: (*bool)(params.IsClosed),
			Wallet:   nil,
		},
		repository.PaginationParams{
			Limit:  (*int)(params.Limit),
			Offset: (*int)(params.Offset),
		},
	)
	if err != nil {
		logrus.WithError(err).Error("failed to GetAdminPositions")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	for _, position := range positions {
		res = append(res, apispec.Position{
			Authority:                position.Authority,
			DcaPeriodIdBeforeDeposit: strconv.FormatUint(position.DcaPeriodIDBeforeDeposit, 10),
			DepositTimestamp:         strconv.FormatInt(position.DepositTimestamp.Unix(), 10),
			DepositedTokenAAmount:    strconv.FormatUint(position.DepositedTokenAAmount, 10),
			IsClosed:                 false,
			NumberOfSwaps:            strconv.FormatUint(position.NumberOfSwaps, 10),
			PeriodicDripAmount:       strconv.FormatUint(position.PeriodicDripAmount, 10),
			Pubkey:                   position.Pubkey,
			Vault:                    position.Vault,
			WithdrawnTokenBAmount:    strconv.FormatUint(position.WithdrawnTokenBAmount, 10),
		})
	}
	return c.JSON(http.StatusOK, res)
}
