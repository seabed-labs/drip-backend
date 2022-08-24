package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/repository"

	"github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetTokenpairs(c echo.Context, params apispec.GetTokenpairsParams) error {
	res := apispec.ListTokenPairs{}
	tokenPairs, err := h.repo.GetTokenPairs(c.Request().Context(),
		repository.TokenPairFilterParams{
			TokenA: (*string)(params.TokenA),
			TokenB: (*string)(params.TokenB),
		})
	if err != nil {
		logrus.WithError(err).Errorf("failed to get tokens")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}
	res = tokenPairDatabaseModelToAPIModel(tokenPairs)
	return c.JSON(http.StatusOK, res)
}
