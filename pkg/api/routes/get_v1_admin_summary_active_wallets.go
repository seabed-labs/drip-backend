package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/dcaf-labs/drip/pkg/repository"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetV1AdminSummaryActivewallets(
	c echo.Context, params apispec.GetV1AdminSummaryActivewalletsParams,
) error {
	res := []apispec.ActiveWallet{}
	activeWallets, err := h.repo.GetActiveWallets(c.Request().Context(), repository.GetActiveWalletParams{
		PositionIsClosed: (*bool)(params.IsClosed),
		Owner:            (*string)(params.Owner),
		Vault:            (*string)(params.Vault),
	})
	if err != nil {
		logrus.WithError(err).Error("failed to GetActiveWallets")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	for _, activeWallet := range activeWallets {
		res = append(res, apispec.ActiveWallet{
			Owner:         activeWallet.Owner,
			PositionCount: activeWallet.PositionCount,
		})
	}
	return c.JSON(http.StatusOK, res)
}
