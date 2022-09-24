package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/service/utils"

	"github.com/dcaf-labs/drip/pkg/api/apispec"

	"github.com/dcaf-labs/drip/pkg/service/repository"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetV1Positions(c echo.Context, params apispec.GetV1PositionsParams) error {
	positions, err := h.repo.GetAdminPositions(
		c.Request().Context(),
		utils.GetBoolPtr(true),
		repository.PositionFilterParams{
			IsClosed: (*bool)(params.IsClosed),
			Wallet:   utils.GetStringPtr(string(params.Wallet)),
		},
		getPaginationParamsFromAPI(params.Offset, params.Limit),
	)
	if err != nil {
		logrus.WithError(err).Error("failed to GetAdminPositions")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	return c.JSON(http.StatusOK, positionModelsToAPI(positions))
}
