package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (h Handler) GetV1Tokens(c echo.Context) error {
	tokens, err := h.repo.GetAllSupportedTokens(c.Request().Context())
	if err != nil {
		log.WithError(err).Errorf("failed to get tokens")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}
	return c.JSON(http.StatusOK, tokenModelsToAPI(tokens))
}
