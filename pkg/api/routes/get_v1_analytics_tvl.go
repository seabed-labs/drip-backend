package controller

import (
	"fmt"
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetV1AnalyticsTvl(c echo.Context) error {
	tvl, err := h.repo.GetCurrentTVL(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	return c.JSON(http.StatusOK, apispec.CurrentTVLResponse{TotalUsdValue: fmt.Sprintf("%f", tvl.TotalUSDValue)})
}
