package controller

import (
	"fmt"
	"net/http"

	"github.com/AlekSi/pointer"
	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetV1AnalyticsOverview(c echo.Context) error {
	ctx := c.Request().Context()
	tvl, err := h.analyticsRepo.GetCurrentTVL(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	lifeTimeDeposit, err := h.analyticsRepo.GetLifeTimeDepositNormalizedToCurrentPrice(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	lifeTimeVolume, err := h.analyticsRepo.GetLifeTimeVolumeNormalizedToCurrentPrice(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	lifeTimeWithdrawal, err := h.analyticsRepo.GetLifeTimeWithdrawalNormalizedToCurrentPrice(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	uniqueDepositorCount, err := h.analyticsRepo.GetUniqueDepositorCount(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	activeWallets, err := h.repo.GetActiveWallets(ctx, repository.GetActiveWalletParams{
		PositionIsClosed: pointer.ToBool(false),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	return c.JSON(http.StatusOK, apispec.AnalyticsOverviewResponse{
		UsdTvl:                fmt.Sprintf("%f", tvl.TotalUSDValue),
		LifeTimeUsdDeposit:    fmt.Sprintf("%f", lifeTimeDeposit.TotalUSDDeposit),
		LifeTimeUsdVolume:     fmt.Sprintf("%f", lifeTimeVolume.TotalUSDVolume),
		LifeTimeUsdWithdrawal: fmt.Sprintf("%f", lifeTimeWithdrawal.TotalUSDWithdrawal),
		ActiveWallets:         float32(len(activeWallets)),
		UniqueDepositors:      float32(uniqueDepositorCount),
	})
}