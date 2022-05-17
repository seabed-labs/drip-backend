package api

import (
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"

	Swagger "github.com/dcaf-protocol/drip/pkg/swagger"
)

func (h Handler) GetVaults(c echo.Context, params Swagger.GetVaultsParams) error {
	var res Swagger.ListVaults

	vaultModels, err := h.drip.GetVaults(c.Request().Context())
	if err != nil {
		logrus.WithError(err).Errorf("failed to get vaults")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal server error"})
	}

	for _, vault := range vaultModels {
		res = append(res, struct {
			DcaActivationTimestamp string `json:"dcaActivationTimestamp"`
			DripAmount             string `json:"dripAmount"`
			LastDcaPeriod          string `json:"lastDcaPeriod"`
			ProtoConfig            string `json:"protoConfig"`
			Pubkey                 string `json:"pubkey"`
			TokenAAccount          string `json:"tokenAAccount"`
			TokenAMint             string `json:"tokenAMint"`
			TokenBAccount          string `json:"tokenBAccount"`
			TokenBMint             string `json:"tokenBMint"`
			TreasuryTokenBAccount  string `json:"treasuryTokenBAccount"`
		}{
			DcaActivationTimestamp: strconv.FormatInt(vault.DcaActivationTimestamp.Unix(), 10),
			DripAmount:             strconv.FormatUint(vault.DripAmount, 10),
			LastDcaPeriod:          strconv.FormatUint(vault.LastDcaPeriod, 10),
			ProtoConfig:            vault.ProtoConfig,
			Pubkey:                 vault.Pubkey,
			TokenAAccount:          vault.TokenAAccount,
			TokenAMint:             vault.TokenAMint,
			TokenBAccount:          vault.TokenBAccount,
			TokenBMint:             vault.TokenBMint,
			TreasuryTokenBAccount:  vault.TreasuryTokenBAccount},
		)
	}
	return c.JSON(http.StatusOK, res)
}
