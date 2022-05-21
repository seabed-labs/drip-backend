package scripts

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dcaf-protocol/drip/internal/pkg/repository/query"

	"github.com/shopspring/decimal"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana/dca_vault"

	"github.com/google/uuid"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/token"
	"gorm.io/gorm/clause"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/dcaf-protocol/drip/internal/pkg/repository/model"

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

// TODO(mocha): sql dump the output of this
func Backfill(
	config *configs.AppConfig,
	repo *query.Query,
) error {
	client := rpc.NewWithCustomRPCClient(rpc.NewWithRateLimit(rpc.DevNet_RPC, 10))

	if !configs.IsDev(config.Environment) {
		logrus.WithField("environment", config.Environment).Infof("skipping event")
		return nil
	}
	logrus.Infof("backfilling devnet vaults")
	configFileName := "./internal/scripts/devnet.yaml"
	configFileName = fmt.Sprintf("%s/%s", configs.GetProjectRoot(), configFileName)
	var vaultConfigs Config
	if err := cleanenv.ReadConfig(configFileName, &vaultConfigs); err != nil {
		return err
	}

	backfillTokens(repo, client, vaultConfigs)
	tokenPairMap := backfillTokenPairs(repo, vaultConfigs)
	backfillTokenSwaps(repo, vaultConfigs, tokenPairMap)
	backfillProtoConfigs(repo, client, vaultConfigs)
	vaultMap := backfillVaults(repo, client, vaultConfigs, tokenPairMap)
	backfillVaultPeriods(repo, client, vaultConfigs, vaultMap)
	logrus.Infof("done backfilling")
	return nil
}

func backfillVaultPeriods(repo *query.Query, client *rpc.Client, vaultConfigs Config, vaultSet map[string]*model.Vault) {
	for _, vaultConfig := range vaultConfigs.TriggerDCAConfigs {
		vault := vaultSet[vaultConfig.Vault]
		logrus.
			WithField("len", vault.LastDcaPeriod+1).
			WithField("vault", vaultConfig.Vault).
			Infof("fetching vaultPeriods")
		var addresses []string

		for i := int64(vault.LastDcaPeriod); i >= 0; i-- {
			vaultPeriodAddr, _, _ := solana.FindProgramAddress([][]byte{
				[]byte("vault_period"),
				solana.MustPrivateKeyFromBase58(vault.Pubkey)[:],
				[]byte(strconv.FormatInt(i, 10)),
			}, dca_vault.ProgramID)
			addresses = append(addresses, vaultPeriodAddr.String())
			if len(addresses) >= 99 || (i == 0) {
				var vaultPeriods []*model.VaultPeriod
				vaultPeriodAccounts, err := getVaultPeriods(client, addresses...)
				addresses = []string{}
				if err != nil {
					logrus.
						WithError(err).
						WithField("vault", vault.Pubkey).
						WithField("i", i).
						Errorf("failed to fetch vault periods")
					continue
				}
				for _, vaultPeriod := range vaultPeriodAccounts {
					pubkey, _, _ := solana.FindProgramAddress([][]byte{
						[]byte("vault_period"),
						solana.MustPrivateKeyFromBase58(vault.Pubkey)[:],
						[]byte(strconv.FormatInt(int64(vaultPeriod.PeriodId), 10)),
					}, dca_vault.ProgramID)
					twap, _ := decimal.NewFromString(vaultPeriod.Twap.String())
					vaultPeriodModel := &model.VaultPeriod{
						Pubkey:   pubkey.String(),
						Vault:    vaultPeriod.Vault.String(),
						PeriodID: vaultPeriod.PeriodId,
						Dar:      vaultPeriod.Dar,
						Twap:     twap,
					}
					vaultPeriods = append(vaultPeriods, vaultPeriodModel)
				}
				if len(vaultPeriods) > 0 {
					logrus.
						WithField("len(vaultPeriods)", len(vaultPeriods)).
						WithField("first_period_in_batch", vaultPeriods[0].PeriodID).
						Infof("inserting vaultPeriods")
					if err := repo.VaultPeriod.
						WithContext(context.Background()).
						Clauses(clause.OnConflict{UpdateAll: true}).
						Create(vaultPeriods...); err != nil {
						logrus.WithError(err).Error("failed to upsert vaultPeriods")
					}
				}
				logrus.
					WithField("vault", vaultConfig.Vault).
					WithField("len", len(vaultPeriods)).
					Infof("inserted vaultPeriods")
			}
		}
	}
}

func backfillVaults(
	repo *query.Query, client *rpc.Client,
	vaultConfigs Config, mintToTokenPair map[string]*model.TokenPair,
) map[string]*model.Vault {
	var vaults []*model.Vault
	vaultSet := make(map[string]*model.Vault)
	for _, vaultConfig := range vaultConfigs.TriggerDCAConfigs {
		if _, ok := vaultSet[vaultConfig.Vault]; !ok {
			var vault dca_vault.Vault
			if err := getAccount(client, vaultConfig.Vault, &vault); err != nil {
				continue
			}
			tokenPair, ok := mintToTokenPair[vaultConfig.TokenAMint+vaultConfig.TokenBMint]
			if !ok {
				logrus.
					WithField("tokenA", vaultConfig.TokenAMint).
					WithField("tokenB", vaultConfig.TokenBMint).
					Warning("missing token pair")
				continue
			}
			vaultModel := &model.Vault{
				Pubkey:                 vaultConfig.Vault,
				ProtoConfig:            vault.ProtoConfig.String(),
				TokenPairID:            tokenPair.ID,
				TokenAAccount:          vault.TokenAAccount.String(),
				TokenBAccount:          vault.TokenBAccount.String(),
				TreasuryTokenBAccount:  vault.TreasuryTokenBAccount.String(),
				LastDcaPeriod:          vault.LastDcaPeriod,
				DripAmount:             vault.DripAmount,
				DcaActivationTimestamp: time.Unix(vault.DcaActivationTimestamp, 0),
				Enabled:                true,
			}

			vaults = append(vaults, vaultModel)
			vaultSet[vaultConfig.Vault] = vaultModel
		}
	}
	logrus.WithField("len", len(vaults)).Infof("inserting vaults")
	if err := repo.Vault.
		WithContext(context.Background()).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(vaults...); err != nil {
		logrus.WithError(err).Error("failed to upsert vaults")
	}
	return vaultSet
}

func backfillTokens(repo *query.Query, client *rpc.Client, vaultConfigs Config) map[string]*model.Token {
	var tokens []*model.Token
	tokenSet := make(map[string]*model.Token)
	for _, vaultConfig := range vaultConfigs.TriggerDCAConfigs {
		if _, ok := tokenSet[vaultConfig.TokenAMint]; !ok {
			var mint token.Mint
			err := getAccount(client, vaultConfig.TokenAMint, &mint)
			if err != nil {
				continue
			}
			tokenModel := &model.Token{
				Pubkey:   vaultConfig.TokenAMint,
				Symbol:   &vaultConfig.TokenASymbol,
				Decimals: int16(mint.Decimals),
				IconURL:  nil, // TODO: hardcode a mint -> iconurl somewhere
			}
			tokens = append(tokens, tokenModel)
			tokenSet[vaultConfig.TokenAMint] = tokenModel
		}
		if _, ok := tokenSet[vaultConfig.TokenBMint]; !ok {
			var mint token.Mint
			err := getAccount(client, vaultConfig.TokenAMint, &mint)
			if err != nil {
				continue
			}
			tokenModel := &model.Token{
				Pubkey:   vaultConfig.TokenBMint,
				Symbol:   &vaultConfig.TokenBSymbol,
				Decimals: int16(mint.Decimals),
				IconURL:  nil, // TODO: hardcode a mint -> iconurl somewhere
			}
			tokens = append(tokens, tokenModel)
			tokenSet[vaultConfig.TokenBMint] = tokenModel
		}
	}
	logrus.WithField("len", len(tokens)).Infof("inserting tokens")
	if err := repo.Token.
		WithContext(context.Background()).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(tokens...); err != nil {
		logrus.WithError(err).Error("failed to upsert tokens")
	}
	return tokenSet
}

func backfillTokenPairs(repo *query.Query, vaultConfigs Config) map[string]*model.TokenPair {
	var tokenPairs []*model.TokenPair
	tokenPairSet := make(map[string]*model.TokenPair)
	for _, vaultConfig := range vaultConfigs.TriggerDCAConfigs {
		tokenPair, err := repo.TokenPair.
			WithContext(context.Background()).
			Where(repo.TokenPair.TokenA.Eq(vaultConfig.TokenAMint)).
			Where(repo.TokenPair.TokenB.Eq(vaultConfig.TokenBMint)).First()
		if err != nil {
			logrus.WithField("len", len(tokenPairs)).Infof("failed to find token pair")
		}
		if tokenPair != nil {
			tokenPairSet[vaultConfig.TokenAMint+vaultConfig.TokenBMint] = tokenPair
		} else {
			if _, ok := tokenPairSet[vaultConfig.TokenAMint+vaultConfig.TokenBMint]; !ok {
				tokenPair := &model.TokenPair{
					ID:     uuid.New().String(),
					TokenA: vaultConfig.TokenAMint,
					TokenB: vaultConfig.TokenBMint,
				}
				tokenPairs = append(tokenPairs, tokenPair)
				tokenPairSet[vaultConfig.TokenAMint+vaultConfig.TokenBMint] = tokenPair
			}
		}
	}
	logrus.WithField("len", len(tokenPairs)).Infof("inserting tokenPairs")
	if err := repo.TokenPair.
		WithContext(context.Background()).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(tokenPairs...); err != nil {
		logrus.WithError(err).Error("failed to upsert tokenPairs")
	}
	return tokenPairSet
}

func backfillTokenSwaps(repo *query.Query, vaultConfigs Config, tokenPairSet map[string]*model.TokenPair) map[string]*model.TokenSwap {
	var tokenSwaps []*model.TokenSwap
	tokenSwapSet := make(map[string]*model.TokenSwap)
	for _, vaultConfig := range vaultConfigs.TriggerDCAConfigs {
		if _, ok := tokenSwapSet[vaultConfig.Swap]; !ok {
			pair := tokenPairSet[vaultConfig.TokenAMint+vaultConfig.TokenBMint]
			tokenSwap := &model.TokenSwap{
				Pubkey:        vaultConfig.Swap,
				Mint:          vaultConfig.SwapTokenMint,
				Authority:     vaultConfig.SwapAuthority,
				FeeAccount:    vaultConfig.SwapFeeAccount,
				TokenAAccount: vaultConfig.SwapTokenAAccount,
				TokenBAccount: vaultConfig.SwapTokenBAccount,
				Pair:          pair.ID,
			}
			tokenSwaps = append(tokenSwaps, tokenSwap)
			tokenSwapSet[vaultConfig.Swap] = tokenSwap
		}
	}
	logrus.WithField("len", len(tokenSwaps)).Infof("inserting tokenSwaps")
	if err := repo.TokenSwap.
		WithContext(context.Background()).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(tokenSwaps...); err != nil {
		logrus.WithError(err).Error("failed to upsert tokenSwaps")
	}
	return tokenSwapSet
}

func backfillProtoConfigs(repo *query.Query, client *rpc.Client, vaultConfigs Config) map[string]*model.ProtoConfig {
	var protoConfigs []*model.ProtoConfig
	protoConfigSet := make(map[string]*model.ProtoConfig)
	for _, vaultConfig := range vaultConfigs.TriggerDCAConfigs {
		if _, ok := protoConfigSet[vaultConfig.VaultProtoConfig]; !ok {
			var protoConfig dca_vault.VaultProtoConfig
			if err := getAccount(client, vaultConfig.VaultProtoConfig, &protoConfig); err != nil {
				continue
			}
			protoConfigModel := &model.ProtoConfig{
				Pubkey:               vaultConfig.VaultProtoConfig,
				Granularity:          protoConfig.Granularity,
				TriggerDcaSpread:     protoConfig.TriggerDcaSpread,
				BaseWithdrawalSpread: protoConfig.BaseWithdrawalSpread,
			}
			protoConfigs = append(protoConfigs, protoConfigModel)
			protoConfigSet[vaultConfig.VaultProtoConfig] = protoConfigModel
		}
	}
	logrus.WithField("len", len(protoConfigs)).Infof("inserting proto configs")
	if err := repo.ProtoConfig.
		WithContext(context.Background()).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(protoConfigs...); err != nil {
		logrus.WithError(err).Error("failed to upsert protoConfigs")
	}
	return protoConfigSet
}

func getVaultPeriods(client *rpc.Client, addresses ...string) ([]dca_vault.VaultPeriod, error) {
	var pubkeys []solana.PublicKey
	for _, address := range addresses {
		pubkeys = append(pubkeys, solana.MustPublicKeyFromBase58(address))
	}
	resp, err := client.GetMultipleAccountsWithOpts(context.Background(), pubkeys, &rpc.GetMultipleAccountsOpts{
		Encoding:   solana.EncodingBase64,
		Commitment: "confirmed",
		DataSlice:  nil,
	})
	if err != nil {
		logrus.
			WithError(err).
			Errorf("couldn't get multiple account infos")
		return nil, err
	}
	var vaultPeriods []dca_vault.VaultPeriod
	for _, val := range resp.Value {
		var vaultPeriod dca_vault.VaultPeriod
		if val == nil || val.Data == nil {
			continue
		}
		if err := bin.NewBinDecoder(val.Data.GetBinary()).Decode(&vaultPeriod); err != nil {
			logrus.
				WithError(err).
				Errorf("failed to decode an account, continueing with the rest...")
			return nil, err
		}
		vaultPeriods = append(vaultPeriods, vaultPeriod)
	}

	return vaultPeriods, nil
}

func getAccount(client *rpc.Client, addr string, v interface{}) error {
	resp, err := client.GetAccountInfoWithOpts(context.Background(), solana.MustPublicKeyFromBase58(addr), &rpc.GetAccountInfoOpts{
		Encoding:   solana.EncodingBase64,
		Commitment: "confirmed",
		DataSlice:  nil,
	})
	if err != nil {
		logrus.
			WithError(err).
			WithField("addr", addr).
			Errorf("couldn't get acount info")
		return err
	}
	if err := bin.NewBinDecoder(resp.Value.Data.GetBinary()).Decode(v); err != nil {
		logrus.
			WithError(err).
			WithField("addr", addr).
			Errorf("failed to decode")
		return err
	}
	return nil
}
