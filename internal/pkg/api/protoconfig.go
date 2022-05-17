package api

import (
	"net/http"

	Swagger "github.com/dcaf-protocol/drip/pkg/swagger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetProtoconfigs(c echo.Context, params Swagger.GetProtoconfigsParams) error {
	var res Swagger.ListProtoConfigs

	protoConfigModels, err := h.drip.GetProtoConfigs(
		c.Request().Context(),
		(*string)(params.TokenA),
		(*string)(params.TokenB),
	)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get proto configs")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal server error"})
	}

	for i := range protoConfigModels {
		protoConfig := protoConfigModels[i]
		res = append(res, struct {
			BaseWithdrawalSpread float32 `json:"baseWithdrawalSpread"`
			Granularity          float32 `json:"granularity"`
			Pubkey               string  `json:"pubkey"`
			TriggerDcaSpread     float32 `json:"triggerDcaSpread"`
		}{
			Pubkey:               protoConfig.Pubkey,
			BaseWithdrawalSpread: float32(protoConfig.BaseWithdrawalSpread),
			Granularity:          float32(protoConfig.Granularity),
			TriggerDcaSpread:     float32(protoConfig.TriggerDcaSpread),
		},
		)
	}
	return c.JSON(http.StatusOK, res)
}
