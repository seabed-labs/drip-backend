package controller

import (
	"net/http"

	Swagger "github.com/dcaf-labs/drip/pkg/swagger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetProtoconfigs(c echo.Context, params Swagger.GetProtoconfigsParams) error {
	var res Swagger.ListProtoConfigs

	protoConfigModels, err := h.repo.GetProtoConfigs(
		c.Request().Context(),
		(*string)(params.TokenA),
		(*string)(params.TokenB),
	)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get proto configs")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal api error"})
	}

	for i := range protoConfigModels {
		protoConfig := protoConfigModels[i]
		res = append(res, struct {
			BaseWithdrawalSpread int    `json:"baseWithdrawalSpread"`
			Granularity          int    `json:"granularity"`
			Pubkey               string `json:"pubkey"`
			TriggerDcaSpread     int    `json:"triggerDcaSpread"`
		}{
			Pubkey:               protoConfig.Pubkey,
			BaseWithdrawalSpread: int(protoConfig.BaseWithdrawalSpread),
			Granularity:          int(protoConfig.Granularity),
			TriggerDcaSpread:     int(protoConfig.TriggerDcaSpread),
		},
		)
	}
	return c.JSON(http.StatusOK, res)
}
