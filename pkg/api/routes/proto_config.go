package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetProtoconfigs(c echo.Context) error {
	res := apispec.ListProtoConfigs{}
	protoConfigModels, err := h.repo.GetProtoConfigs(c.Request().Context())
	if err != nil {
		logrus.WithError(err).Errorf("failed to get proto configs")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}
	res = vaultProtoConfigDataBaseModelToAPIModel(protoConfigModels)
	return c.JSON(http.StatusOK, res)
}
