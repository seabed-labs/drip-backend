package controller

import (
	"net/http"

	Swagger "github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetTokenpairs(c echo.Context, params Swagger.GetTokenpairsParams) error {
	res := Swagger.ListTokenPairs{}

	tokenPairs, err := h.repo.GetTokenPairs(c.Request().Context(), (*string)(params.TokenA), (*string)(params.TokenB))
	if err != nil {
		logrus.WithError(err).Errorf("failed to get tokens")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal api error"})
	}

	for i := range tokenPairs {
		tokenPair := tokenPairs[i]
		res = append(res, struct {
			Id     string `json:"id"`
			TokenA string `json:"tokenA"`
			TokenB string `json:"tokenB"`
		}{Id: tokenPair.ID, TokenA: tokenPair.TokenA, TokenB: tokenPair.TokenB},
		)
	}
	return c.JSON(http.StatusOK, res)
}
