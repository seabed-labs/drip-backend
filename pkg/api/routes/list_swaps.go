package controller

import (
	"net/http"

	Swagger "github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetSwaps(c echo.Context, params Swagger.GetSwapsParams) error {
	res := Swagger.ListTokenSwaps{}

	tokenSwaps, err := h.repo.GetTokenSwaps(c.Request().Context(), []string{string(*params.TokenPair)})
	if err != nil {
		logrus.WithError(err).Errorf("failed to get tokens")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal api error"})
	}

	for i := range tokenSwaps {
		tokenSwap := tokenSwaps[i]
		res = append(res, struct {
			Authority  string `json:"authority"`
			FeeAccount string `json:"feeAccount"`
			Mint       string `json:"mint"`

			// token pair reference identifier
			Pair          string `json:"pair"`
			Pubkey        string `json:"pubkey"`
			TokenAAccount string `json:"tokenAAccount"`
			TokenBAccount string `json:"tokenBAccount"`
		}{
			Authority:     tokenSwap.Authority,
			FeeAccount:    tokenSwap.FeeAccount,
			Mint:          tokenSwap.Mint,
			Pair:          tokenSwap.TokenPairID,
			Pubkey:        tokenSwap.Pubkey,
			TokenAAccount: tokenSwap.TokenAAccount,
			TokenBAccount: tokenSwap.TokenBAccount,
		})
	}
	return c.JSON(http.StatusOK, res)
}
