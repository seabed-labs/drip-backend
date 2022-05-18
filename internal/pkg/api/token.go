package api

import (
	"net/http"

	Swagger "github.com/dcaf-protocol/drip/pkg/swagger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetTokens(c echo.Context, params Swagger.GetTokensParams) error {
	var res Swagger.ListTokens
	if params.TokenA != nil && params.TokenB != nil {
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "both tokenA and tokenB cannot be set"})
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
	tokens, err := h.drip.GetTokensWithSupportedTokenPair(c.Request().Context(), mintFilter, supportedTokenA)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get tokens")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal server error"})
	}
	for i := range tokens {
		token := tokens[i]
		symbol := ""
		if token.Symbol != nil {
			symbol = *token.Symbol
		}
		res = append(res, struct {
			Decimals int    `json:"decimals"`
			Pubkey   string `json:"pubkey"`
			Symbol   string `json:"symbol"`
		}{Decimals: int(token.Decimals), Pubkey: token.Pubkey, Symbol: symbol},
		)
	}
	return c.JSON(http.StatusOK, res)
}
