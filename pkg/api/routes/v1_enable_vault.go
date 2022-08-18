package controller

import (
	"errors"
	"net/http"

	apispec "github.com/dcaf-labs/drip/pkg/swagger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (h Handler) PutAdminVaultPubkeyPathEnable(
	c echo.Context, pubkeyPath apispec.PubkeyPathParam, params apispec.PutAdminVaultPubkeyPathEnableParams,
) error {
	return h.PutV1AdminVaultPubkeyPathEnable(c, pubkeyPath, apispec.PutV1AdminVaultPubkeyPathEnableParams(params))
}

func (h Handler) PutV1AdminVaultPubkeyPathEnable(
	c echo.Context,
	pubkeyPath apispec.PubkeyPathParam,
	_ apispec.PutV1AdminVaultPubkeyPathEnableParams,
) error {
	vault, err := h.repo.AdminGetVaultByAddress(c.Request().Context(), string(pubkeyPath))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusBadRequest, apispec.ErrorResponse{Error: "invalid vault pubkey"})
	}
	updatedVault, err := h.repo.AdminSetVaultEnabled(c.Request().Context(), string(pubkeyPath), !vault.Enabled)
	if err != nil {
		logrus.
			WithError(err).
			WithField("vault", string(pubkeyPath)).
			Error("failed to enable vault")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "something went wrong"})
	}
	return c.JSON(http.StatusOK, updatedVault)
}
