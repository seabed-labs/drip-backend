package scripts

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/processor"

	configs2 "github.com/dcaf-labs/drip/pkg/configs"
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
	//if !configs2.IsDev(config.Environment) {
	//	logrus.WithField("environment", config.Environment).Infof("skipping backfill")
	//	return nil
	//}
	//logrus.Infof("backfilling devnet vaults")
	//configFileName := "./internal/scripts/devnet.yaml"
	//configFileName = fmt.Sprintf("%s/%s", configs2.GetProjectRoot(), configFileName)
	//var vaultConfigs Config
	//if err := cleanenv.ReadConfig(configFileName, &vaultConfigs); err != nil {
	//	return err
	//}
	//backfillTokenPairs(vaultConfigs, processor)
	//backfillTokenSwaps(vaultConfigs, processor)
	//backfillProtoConfigs(vaultConfigs, processor)
	//backfillVaults(vaultConfigs, processor)
	//backfillTokens(repo, client, vaultConfigs)
	//backfillVaultPeriods(repo, client, vaultConfigs, vaultMap)
	for _, address := range []string{
		"2w9DNJRFGmvN2wuVi3CtLT49cM8DWdKfzAGas1XDw3Ve",
		"Dr75jCuqpkYCGs4ATp2v31sU6bzDKoCdxJNvxXUUgb4S",
	} {
		log := logrus.WithField("address", address)
		if err := processor.UpsertWhirlpoolByAddress(context.Background(), address); err != nil {
			log.WithError(err).Error("failed to backfill whirlpool")
		}
	}

	for _, address := range []string{
		"FcbTUmkEuizRjnr17iggHon9rAc2FeQ8uMdjxWBgFb58",
		"Evy3P7kFy6epfbY3QZmwWiqWj1veZ87LQoCEWYzb3vW7",
	} {
		log := logrus.WithField("address", address)
		if err := processor.UpsertTokenSwapByAddress(context.Background(), address); err != nil {
			log.WithError(err).Error("failed to backfill spl token swap")
		}
	}

	for _, address := range []string{
		"GHyHWLdLCRkYrXwxaxVrG5LEnkPEP9xVHqXPv3QMYMVZ",
		"HrQ5YfEWtsc4zNxQqv6dKT8136nsSfDpXWeXnn76xBGs",
		"8RVA7nF7xmieAV5mJSM7CHGegdm68ZqghvZBGfN9K4mY",
		"A6yjGfL4qC2y9KXpJGd6eEwHuXoyUsr2QSfBvPNFQDb6",
	} {
		log := logrus.WithField("address", address)
		if err := processor.UpsertVaultByAddress(context.Background(), address); err != nil {
			log.WithError(err).Error("failed to backfill vault")
		}
	}

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
