package scripts

import (
	"context"
	"fmt"

	configs2 "github.com/dcaf-protocol/drip/pkg/configs"
	"github.com/dcaf-protocol/drip/pkg/processor"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Environment       string             `yaml:"environment" env:"ENV"`
	Wallet            string             `yaml:"wallet"      env:"KEEPER_BOT_WALLET"`
	TriggerDCAConfigs []TriggerDCAConfig `yaml:"triggerDCA"`
}

type TriggerDCAConfig struct {
	Vault              string `yaml:"vault"`
	VaultProtoConfig   string `yaml:"vaultProtoConfig"`
	VaultTokenAAccount string `yaml:"vaultTokenAAccount"`
	VaultTokenBAccount string `yaml:"vaultTokenBAccount"`
	TokenAMint         string `yaml:"tokenAMint"`
	TokenASymbol       string `yaml:"tokenASymbol"`
	TokenBMint         string `yaml:"tokenBMint"`
	TokenBSymbol       string `yaml:"tokenBSymbol"`
	SwapTokenMint      string `yaml:"swapTokenMint"`
	SwapTokenAAccount  string `yaml:"swapTokenAAccount"`
	SwapTokenBAccount  string `yaml:"swapTokenBAccount"`
	SwapFeeAccount     string `yaml:"swapFeeAccount"`
	SwapAuthority      string `yaml:"swapAuthority"`
	Swap               string `yaml:"swap"`
}

func Backfill(
	config *configs2.AppConfig,
	processor processor.Processor,
) error {
	if !configs2.IsDev(config.Environment) {
		logrus.WithField("environment", config.Environment).Infof("skipping backfill")
		return nil
	}
	logrus.Infof("backfilling devnet vaults")
	configFileName := "./internal/scripts/devnet.yaml"
	configFileName = fmt.Sprintf("%s/%s", configs2.GetProjectRoot(), configFileName)
	var vaultConfigs Config
	if err := cleanenv.ReadConfig(configFileName, &vaultConfigs); err != nil {
		return err
	}
	backfillTokenPairs(vaultConfigs, processor)
	backfillTokenSwaps(vaultConfigs, processor)
	backfillProtoConfigs(vaultConfigs, processor)
	backfillVaults(vaultConfigs, processor)
	//backfillTokens(repo, client, vaultConfigs)
	//backfillVaultPeriods(repo, client, vaultConfigs, vaultMap)
	logrus.Infof("done backfilling")
	return nil
}

func backfillVaults(
	vaultConfigs Config,
	processor processor.Processor,
) {
	for _, vaultConfig := range vaultConfigs.TriggerDCAConfigs {
		log := logrus.WithField("address", vaultConfig.Vault)
		if err := processor.UpsertVaultByAddress(context.Background(), vaultConfig.Vault); err != nil {
			log.WithError(err).Error("failed to backfill vault")
		}
		log.Info("backfilled vault")
	}
}

func backfillTokenPairs(
	vaultConfigs Config, processor processor.Processor,
) {
	ctx := context.Background()
	for _, vaultConfig := range vaultConfigs.TriggerDCAConfigs {
		log := logrus.WithField("address", vaultConfig.Vault)
		if err := processor.UpsertTokenPair(ctx, vaultConfig.TokenAMint, vaultConfig.TokenBMint); err != nil {
			log.WithError(err).Error("failed to backfill token pair")
		}
		log.Info("backfilled tokenPair")
	}
}

func backfillTokenSwaps(
	vaultConfigs Config, processor processor.Processor,
) {
	ctx := context.Background()
	for _, vaultConfig := range vaultConfigs.TriggerDCAConfigs {
		log := logrus.WithField("address", vaultConfig.Vault)
		if err := processor.UpsertTokenSwapByAddress(ctx, vaultConfig.Swap); err != nil {
			log.WithError(err).Error("failed to backfill vault")
		}
		log.Info("backfilled tokenSwap")
	}
}

func backfillProtoConfigs(
	vaultConfigs Config, processor processor.Processor,
) {
	ctx := context.Background()
	for _, vaultConfig := range vaultConfigs.TriggerDCAConfigs {
		log := logrus.WithField("address", vaultConfig.Vault)
		if err := processor.UpsertProtoConfigByAddress(ctx, vaultConfig.VaultProtoConfig); err != nil {
			log.WithError(err).Error("failed to backfill vault")
		}
		log.Info("backfilled protoConfig")
	}
}
