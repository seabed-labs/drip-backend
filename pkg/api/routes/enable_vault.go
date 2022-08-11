package controller

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	Swagger "github.com/dcaf-labs/drip/pkg/swagger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) PutAdminVaultPubkeyPathEnable(
	c echo.Context, pubkeyPath Swagger.PubkeyPathParam, _ Swagger.PutAdminVaultPubkeyPathEnableParams,
) error {
	_, err := h.repo.AdminGetVaultByAddress(c.Request().Context(), string(pubkeyPath))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusBadRequest, Swagger.ErrorResponse{Error: "invalid vault pubkey"})
	}
	vault, err := h.repo.GetVaultByAddress(c.Request().Context(), string(pubkeyPath))
	if err != nil {
		logrus.
			WithError(err).
			WithField("vault", string(pubkeyPath)).
			Error("failed to get vault")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "something went wrong"})
	}
	updatedVault, err := h.repo.AdminSetVaultEnabled(c.Request().Context(), string(pubkeyPath), !vault.Enabled)
	if err != nil {
		logrus.
			WithError(err).
			WithField("vault", string(pubkeyPath)).
			Error("failed to enable vault")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "something went wrong"})
	}
	return c.JSON(http.StatusOK, updatedVault)
}
