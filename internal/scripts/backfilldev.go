package scripts

//
//import (
//	"context"
//	"fmt"
//	"time"
//
//	"github.com/dcaf-protocol/drip/internal/pkg/repository/models"
//
//	"github.com/volatiletech/sqlboiler/v4/types"
//
//	"github.com/sirupsen/logrus"
//
//	"github.com/dcaf-protocol/drip/internal/configs"
//	"github.com/dcaf-protocol/drip/internal/pkg/repository"
//)
//
//type vaultConfig struct {
//	vaults []struct {
//		Vault                       string `yaml:"vault"`
//		VaultProtoConfig            string `yaml:"vaultProtoConfig"`
//		VaultProtoConfigGranularity uint32 `yaml:"vaultProtoConfigGranularity"`
//		VaultTokenAAccount          string `yaml:"vaultTokenAAccount"`
//		VaultTokenBAccount          string `yaml:"vaultTokenBAccount"`
//		VaultTreasuryTokenBAccount  string `yaml:"vaultTreasuryTokenBAccount"`
//		TokenAMint                  string `yaml:"tokenAMint"`
//		TokenASymbol                string `yaml:"tokenASymbol"`
//		TokenBMint                  string `yaml:"tokenBMint"`
//		TokenBSymbol                string `yaml:"tokenBSymbol"`
//		SwapTokenMint               string `yaml:"swapTokenMint"`
//		SwapTokenAAccount           string `yaml:"swapTokenAAccount"`
//		SwapTokenBAccount           string `yaml:"swapTokenBAccount"`
//		SwapFeeAccount              string `yaml:"swapFeeAccount"`
//		SwapAuthority               string `yaml:"swapAuthority"`
//		Swap                        string `yaml:"swap"`
//	} `yaml:"vaults"`
//}
//
//func Backfill(
//	config *configs.AppConfig,
//	vaultRepo repository.VaultRepository,
//) error {
//	if !configs.IsDev(config.Environment) {
//		return nil
//	}
//	logrus.Infof("backfilling devnet vaults")
//	configFileName := "./internal/scripts/devnet.yaml"
//	configFileName = fmt.Sprintf("%s/%s", configs.GetProjectRoot(), configFileName)
//	var vaultConfigs vaultConfig
//	err := configs.ParseToConfig(&vaultConfigs, configFileName)
//	if err != nil {
//		logrus.WithError(err).Error("failed to load devnet config for backfill")
//		return err
//	}
//	for _, vaultConfig := range vaultConfigs.vaults {
//		vault := models.Vault{
//			Pubkey:                 vaultConfig.Vault,
//			ProtoConfig:            vaultConfig.VaultProtoConfig,
//			TokenAMint:             vaultConfig.TokenAMint,
//			TokenBMint:             vaultConfig.TokenBMint,
//			TokenAAccount:          vaultConfig.VaultTokenAAccount,
//			TokenBAccount:          vaultConfig.VaultTokenBAccount,
//			TreasuryTokenBAccount:  vaultConfig.VaultTreasuryTokenBAccount,
//			LastDcaPeriod:          types.Decimal{}, //  TODO: get from onchain in worker thread
//			DripAmount:             types.Decimal{}, // TODO: get from onchain in worker thread
//			DcaActivationTimestamp: time.Time{},     // TODO: get from onchain in worker thread
//		}
//		err := vaultRepo.UpsertVault(context.Background(), vault)
//		logrus.WithError(err).Error("failed to upsert vault")
//	}
//	return nil
//}
