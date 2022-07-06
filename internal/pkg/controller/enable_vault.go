package controller

import (
	"database/sql"
	"net/http"

	Swagger "github.com/dcaf-protocol/drip/pkg/swagger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) PutAdminVaultPubkeyPathEnable(
	c echo.Context, pubkeyPath Swagger.PubkeyPath, params Swagger.PutAdminVaultPubkeyPathEnableParams,
) error {
	_, err := h.repo.InternalGetVaultByAddress(c.Request().Context(), string(pubkeyPath))
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusBadRequest, Swagger.ErrorResponse{Error: "invalid vault pubkey"})
	}
	vault, err := h.repo.EnableVault(c.Request().Context(), string(pubkeyPath))
	if err != nil {
		logrus.
			WithError(err).
			WithField("vault", string(pubkeyPath)).
			Error("failed to enable vault")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "something went wrong"})
	}
	return c.JSON(http.StatusOK, vault)
}
