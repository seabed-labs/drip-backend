package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (h Handler) GetV1Tokens(c echo.Context) error {
	res := apispec.ListTokens{}
	tokens, err := h.repo.GetAllSupportTokens(c.Request().Context())
	if err != nil {
		log.WithError(err).Errorf("failed to get tokens")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}
	for i := range tokens {
		token := tokens[i]
		res = append(res, apispec.Token{
			Decimals: int(token.Decimals),
			Pubkey:   token.Pubkey,
			Symbol:   token.Symbol,
			IconUrl:  token.IconURL,
		})
	}
	return c.JSON(http.StatusOK, res)
}
