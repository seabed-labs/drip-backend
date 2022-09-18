package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetTokens(c echo.Context, params apispec.GetTokensParams) error {
	res := apispec.ListTokens{}
	if params.TokenA != nil && params.TokenB != nil {
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "both tokenA and tokenB cannot be set"})
	}
	var mintFilter *string
	var supportedTokenA bool
	if params.TokenA != nil {
		mintFilter = (*string)(params.TokenA)
		supportedTokenA = true
	} else {
		mintFilter = (*string)(params.TokenB)
		supportedTokenA = false
	}
	tokens, err := h.repo.GetTokensWithSupportedTokenPair(c.Request().Context(), mintFilter, supportedTokenA)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get tokens")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}
	for i := range tokens {
		token := tokens[i]
		res = append(res, apispec.Token{
			Decimals: int(token.Decimals),
			Pubkey:   token.Pubkey,
			Symbol:   token.Symbol,
		})
	}
	return c.JSON(http.StatusOK, res)
}
