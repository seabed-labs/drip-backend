package scripts

import "github.com/dcaf-protocol/drip/internal/configs"

type VaultConfig struct {
	Vault                       string `yaml:"vault"`
	VaultProtoConfig            string `yaml:"vaultProtoConfig"`
	VaultProtoConfigGranularity uint32 `yaml:"vaultProtoConfigGranularity"`
	VaultTokenAAccount          string `yaml:"vaultTokenAAccount"`
	VaultTokenBAccount          string `yaml:"vaultTokenBAccount"`
	VaultTreasuryTokenBAccount  string `yaml:"vaultTreasuryTokenBAccount"`
	TokenAMint                  string `yaml:"tokenAMint"`
	TokenASymbol                string `yaml:"tokenASymbol"`
	TokenBMint                  string `yaml:"tokenBMint"`
	TokenBSymbol                string `yaml:"tokenBSymbol"`
	SwapTokenMint               string `yaml:"swapTokenMint"`
	SwapTokenAAccount           string `yaml:"swapTokenAAccount"`
	SwapTokenBAccount           string `yaml:"swapTokenBAccount"`
	SwapFeeAccount              string `yaml:"swapFeeAccount"`
	SwapAuthority               string `yaml:"swapAuthority"`
	Swap                        string `yaml:"swap"`
}

func backfill(config *configs.PSQLConfig) error {

	return nil
}
