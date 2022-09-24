package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetV1VaultTokens(c echo.Context, params apispec.GetV1VaultTokensParams) error {
	if params.TokenA != nil && params.TokenB != nil {
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "both tokenA and tokenB cannot be set"})
	}
	var tokens []*model.Token
	var err error
	if params.TokenA != nil {
		tokens, err = h.repo.GetSupportedTokenBs(c.Request().Context(), string(*params.TokenA))
	} else {
		tokens, err = h.repo.GetSupportedTokenAs(c.Request().Context(), (*string)(params.TokenB))
	}
	if err != nil {
		logrus.WithError(err).Errorf("failed to get tokens")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}
	return c.JSON(http.StatusOK, tokenModelsToAPI(tokens))
}
