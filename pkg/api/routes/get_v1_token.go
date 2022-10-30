package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetV1TokenPubkeyPath(c echo.Context, pubkeyPath apispec.PubkeyPathParam) error {
	log := logrus.WithField("handler", "GetV1TokenPubkeyPath").WithField("pubkey", string(pubkeyPath))
	tokens, err := h.repo.GetTokensByAddresses(c.Request().Context(), []string{string(pubkeyPath)})
	if err != nil {
		log.WithError(err).Errorf("failed to get tokens")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}
	if len(tokens) != 1 {
		log.
			WithField("len(tokens", len(tokens)).
			WithError(err).Errorf("unexpected number of tokens returned")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}
	return c.JSON(http.StatusOK, tokenModelToApi(tokens[0]))
}
