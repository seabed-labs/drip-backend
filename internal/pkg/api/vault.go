package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	Swagger "github.com/dcaf-protocol/drip/pkg/swagger"
)

func (h Handler) GetVaults(c echo.Context, params Swagger.GetVaultsParams) error {
	var res Swagger.ListVaultsResponse
	for _, vault := range h.vaultConfigs {
		if params.TokenA != nil && vault.TokenAMint != string(*params.TokenA) {
			continue
		}
		if params.TokenB != nil && vault.TokenBMint != string(*params.TokenB) {
			continue
		}
		if params.Granularity != nil && float32(vault.VaultProtoConfigGranularity) != float32(*params.Granularity) {
			continue
		}
		res = append(res, struct {
			Swap                       string `json:"swap"`
			SwapAuthority              string `json:"swapAuthority"`
			SwapFeeAccount             string `json:"swapFeeAccount"`
			SwapTokenAAccount          string `json:"swapTokenAAccount"`
			SwapTokenBAccount          string `json:"swapTokenBAccount"`
			SwapTokenMint              string `json:"swapTokenMint"`
			TokenAMint                 string `json:"tokenAMint"`
			TokenASymbol               string `json:"tokenASymbol"`
			TokenBMint                 string `json:"tokenBMint"`
			TokenBSymbol               string `json:"tokenBSymbol"`
			Vault                      string `json:"vault"`
			VaultProtoConfig           string `json:"vaultProtoConfig"`
			VaultTokenAAcount          string `json:"vaultTokenAAcount"`
			VaultTokenBAccount         string `json:"vaultTokenBAccount"`
			VaultTreasuryTokenBAccount string `json:"vaultTreasuryTokenBAccount"`
		}{
			Swap:                       vault.Swap,
			SwapAuthority:              vault.SwapAuthority,
			SwapFeeAccount:             vault.SwapFeeAccount,
			SwapTokenAAccount:          vault.SwapTokenAAccount,
			SwapTokenBAccount:          vault.SwapTokenBAccount,
			SwapTokenMint:              vault.SwapTokenMint,
			TokenAMint:                 vault.TokenAMint,
			TokenASymbol:               vault.TokenASymbol,
			TokenBMint:                 vault.TokenBMint,
			Vault:                      vault.Vault,
			VaultProtoConfig:           vault.VaultProtoConfig,
			VaultTokenAAcount:          vault.VaultTokenAAccount,
			VaultTokenBAccount:         vault.VaultTokenBAccount,
			VaultTreasuryTokenBAccount: vault.VaultTreasuryTokenBAccount,
		})
	}
	return c.JSON(http.StatusOK, res)
}
