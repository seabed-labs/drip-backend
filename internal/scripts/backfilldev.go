package scripts

import (
	"context"
	"fmt"

	"github.com/dcaf-protocol/drip/internal/pkg/processor"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/dcaf-protocol/drip/internal/configs"
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
	config *configs.AppConfig,
	processor processor.Processor,
) error {
	if !configs.IsDev(config.Environment) {
		logrus.WithField("environment", config.Environment).Infof("skipping backfill")
		return nil
	}
	logrus.Infof("backfilling devnet vaults")
	configFileName := "./internal/scripts/devnet.yaml"
	configFileName = fmt.Sprintf("%s/%s", configs.GetProjectRoot(), configFileName)
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

//func backfillVaultPeriods(repo *query.Query, client *rpc.Client, vaultConfigs Config, vaultSet map[string]*model.Vault) {
//	for _, vaultConfig := range vaultConfigs.TriggerDCAConfigs {
//		vault := vaultSet[vaultConfig.Vault]
//		logrus.
//			WithField("len", vault.LastDcaPeriod+1).
//			WithField("vault", vaultConfig.Vault).
//			Infof("fetching vaultPeriods")
//		var addresses []string
//
//		for i := int64(vault.LastDcaPeriod); i >= 0; i-- {
//			vaultPeriodAddr, _, _ := solana.FindProgramAddress([][]byte{
//				[]byte("vault_period"),
//				solana.MustPrivateKeyFromBase58(vault.Pubkey)[:],
//				[]byte(strconv.FormatInt(i, 10)),
//			}, dca_vault.ProgramID)
//			addresses = append(addresses, vaultPeriodAddr.String())
//			if len(addresses) >= 99 || (i == 0) {
//				var vaultPeriods []*model.VaultPeriod
//				vaultPeriodAccounts, err := getVaultPeriods(client, addresses...)
//				addresses = []string{}
//				if err != nil {
//					logrus.
//						WithError(err).
//						WithField("vault", vault.Pubkey).
//						WithField("i", i).
//						Errorf("failed to fetch vault periods")
//					continue
//				}
//				for _, vaultPeriod := range vaultPeriodAccounts {
//					pubkey, _, _ := solana.FindProgramAddress([][]byte{
//						[]byte("vault_period"),
//						solana.MustPrivateKeyFromBase58(vault.Pubkey)[:],
//						[]byte(strconv.FormatInt(int64(vaultPeriod.PeriodId), 10)),
//					}, dca_vault.ProgramID)
//					twap, _ := decimal.NewFromString(vaultPeriod.Twap.String())
//					vaultPeriodModel := &model.VaultPeriod{
//						Pubkey:   pubkey.String(),
//						Vault:    vaultPeriod.Vault.String(),
//						PeriodID: vaultPeriod.PeriodId,
//						Dar:      vaultPeriod.Dar,
//						Twap:     twap,
//					}
//					vaultPeriods = append(vaultPeriods, vaultPeriodModel)
//				}
//				if len(vaultPeriods) > 0 {
//					logrus.
//						WithField("len(vaultPeriods)", len(vaultPeriods)).
//						WithField("first_period_in_batch", vaultPeriods[0].PeriodID).
//						Infof("inserting vaultPeriods")
//					if err := repo.VaultPeriod.
//						WithContext(context.Background()).
//						Clauses(clause.OnConflict{UpdateAll: true}).
//						Create(vaultPeriods...); err != nil {
//						logrus.WithError(err).Error("failed to upsert vaultPeriods")
//					}
//				}
//				logrus.
//					WithField("vault", vaultConfig.Vault).
//					WithField("len", len(vaultPeriods)).
//					Infof("inserted vaultPeriods")
//			}
//		}
//	}
//}
//
//func getVaultPeriods(client *rpc.Client, addresses ...string) ([]dca_vault.VaultPeriod, error) {
//	var pubkeys []solana.PublicKey
//	for _, address := range addresses {
//		pubkeys = append(pubkeys, solana.MustPublicKeyFromBase58(address))
//	}
//	resp, err := client.GetMultipleAccountsWithOpts(context.Background(), pubkeys, &rpc.GetMultipleAccountsOpts{
//		Encoding:   solana.EncodingBase64,
//		Commitment: "confirmed",
//		DataSlice:  nil,
//	})
//	if err != nil {
//		logrus.
//			WithError(err).
//			Errorf("couldn't get multiple account infos")
//		return nil, err
//	}
//	var vaultPeriods []dca_vault.VaultPeriod
//	for _, val := range resp.Value {
//		var vaultPeriod dca_vault.VaultPeriod
//		if val == nil || val.Data == nil {
//			continue
//		}
//		if err := bin.NewBinDecoder(val.Data.GetBinary()).Decode(&vaultPeriod); err != nil {
//			logrus.
//				WithError(err).
//				Errorf("failed to decode an account, continueing with the rest...")
//			return nil, err
//		}
//		vaultPeriods = append(vaultPeriods, vaultPeriod)
//	}
//
//	return vaultPeriods, nil
//}
