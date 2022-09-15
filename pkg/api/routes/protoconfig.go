package controller

import (
	"net/http"
	"strconv"

	Swagger "github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetProtoconfigs(c echo.Context) error {
	res := Swagger.ListProtoConfigs{}

	protoConfigModels, err := h.repo.GetProtoConfigs(c.Request().Context())
	if err != nil {
		logrus.WithError(err).Errorf("failed to get proto configs")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal api error"})
	}

	for i := range protoConfigModels {
		protoConfig := protoConfigModels[i]
		res = append(res, Swagger.ProtoConfig{
			Pubkey:                  protoConfig.Pubkey,
			Admin:                   protoConfig.Admin,
			Granularity:             strconv.FormatUint(protoConfig.Granularity, 10),
			TokenADripTriggerSpread: int(protoConfig.TokenADripTriggerSpread),
			TokenBWithdrawalSpread:  int(protoConfig.TokenBWithdrawalSpread),
			TokenBReferralSpread:    int(protoConfig.TokenBReferralSpread),
		},
		)
	}
	return c.JSON(http.StatusOK, res)
}
