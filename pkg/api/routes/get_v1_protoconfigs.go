package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"

	"github.com/dcaf-labs/drip/pkg/service/repository"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetV1Protoconfigs(c echo.Context, params apispec.GetV1ProtoconfigsParams) error {
	protoConfigModels, err := h.repo.GetProtoConfigs(c.Request().Context(), repository.ProtoConfigParams{
		TokenA: (*string)(params.TokenA),
		TokenB: (*string)(params.TokenB),
	})
	if err != nil {
		logrus.WithError(err).Errorf("failed to get proto config")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}
	return c.JSON(http.StatusOK, protoConfigModelsToAPI(protoConfigModels))
}
