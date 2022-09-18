package controller

import (
	"net/http"
	"strconv"

	"github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetV1Protoconfigs(c echo.Context) error {
	res := apispec.ListProtoConfigs{}

	protoConfigModels, err := h.repo.GetProtoConfigs(c.Request().Context())
	if err != nil {
		logrus.WithError(err).Errorf("failed to get proto configs")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}

	for i := range protoConfigModels {
		protoConfig := protoConfigModels[i]
		res = append(res, apispec.ProtoConfig{
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
