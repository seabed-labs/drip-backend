package controller

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/repository"

	"github.com/dcaf-labs/drip/pkg/service/utils"

	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/test-go/testify/assert"
)

func TestUtils(t *testing.T) {

	t.Run("getServerURL should return correct api URL", func(t *testing.T) {
		assert.True(t, strings.Contains(getServerURL(config.NilNetwork, config.StagingEnv, 0), "localhost"))
		assert.True(t, strings.Contains(getServerURL(config.LocalNetwork, config.StagingEnv, 0), "localhost"))
		assert.True(t, strings.Contains(getServerURL(config.DevnetNetwork, config.StagingEnv, 0), "devnet"))
		assert.True(t, strings.Contains(getServerURL(config.MainnetNetwork, config.StagingEnv, 0), "mainnet"))
	})

	t.Run("hasValue should return true", func(t *testing.T) {
		assert.True(t, hasValue([]string{"1", "2"}, "1"))
		assert.True(t, hasValue([]string{"1", "2"}, "2"))
	})

	t.Run("hasValue should return false", func(t *testing.T) {
		assert.False(t, hasValue([]string{"1", "2"}, "3"))
		assert.False(t, hasValue([]string{}, ""))
	})

	t.Run("vaultModelToAPI should return correct apiSpec model", func(t *testing.T) {
		dbModel := model.Vault{
			DcaActivationTimestamp: time.Unix(1667076855, 0),
			Pubkey:                 solana.NewWallet().PublicKey().String(),
			ProtoConfig:            solana.NewWallet().PublicKey().String(),
			TokenAAccount:          solana.NewWallet().PublicKey().String(),
			TokenBAccount:          solana.NewWallet().PublicKey().String(),
			TreasuryTokenBAccount:  solana.NewWallet().PublicKey().String(),
			LastDcaPeriod:          10,
			DripAmount:             99,
			Enabled:                true,
			TokenPairID:            uuid.New().String(),
			MaxSlippageBps:         50,
			TokenAMint:             solana.NewWallet().PublicKey().String(),
			TokenBMint:             solana.NewWallet().PublicKey().String(),
		}
		apiModel := vaultModelToAPI(&dbModel)
		assert.Equal(t, strconv.FormatInt(dbModel.DcaActivationTimestamp.Unix(), 10), apiModel.DcaActivationTimestamp)
		assert.Equal(t, dbModel.Pubkey, apiModel.Pubkey)
		assert.Equal(t, dbModel.ProtoConfig, apiModel.ProtoConfig)
		assert.Equal(t, dbModel.TokenAAccount, apiModel.TokenAAccount)
		assert.Equal(t, dbModel.TokenBAccount, apiModel.TokenBAccount)
		assert.Equal(t, dbModel.TreasuryTokenBAccount, apiModel.TreasuryTokenBAccount)
		assert.Equal(t, dbModel.TokenAMint, apiModel.TokenAMint)
		assert.Equal(t, dbModel.TokenBMint, apiModel.TokenBMint)
		assert.Equal(t, strconv.FormatUint(dbModel.LastDcaPeriod, 10), apiModel.LastDcaPeriod)
		assert.Equal(t, strconv.FormatUint(dbModel.DripAmount, 10), apiModel.DripAmount)
		assert.Equal(t, dbModel.Enabled, apiModel.Enabled)
		assert.Equal(t, int(dbModel.MaxSlippageBps), apiModel.MaxSlippageBps)

		// if the line below needs to be updated, add the field assertion above
		assert.Equal(t, reflect.TypeOf(apiModel).NumField(), 12)
	})

	t.Run("vaultModelsToAPI should return correct apiSpec model", func(t *testing.T) {
		dbModels := []*model.Vault{
			{},
			{},
		}
		apiModel := vaultModelsToAPI(dbModels)
		for i := range apiModel {
			assert.Equal(t, reflect.TypeOf(apiModel[i]).NumField(), 12)
		}
	})

	t.Run("vaultPeriodModelToAPI should return correct apiSpec model", func(t *testing.T) {
		dbModel := model.VaultPeriod{
			Pubkey:      solana.NewWallet().PublicKey().String(),
			Vault:       solana.NewWallet().PublicKey().String(),
			PeriodID:    91,
			Twap:        decimal.NewFromFloat(1.0),
			Dar:         10,
			PriceBOverA: decimal.Decimal{},
		}
		apiModel := vaultPeriodModelToAPI(&dbModel)
		assert.Equal(t, dbModel.Pubkey, apiModel.Pubkey)
		assert.Equal(t, dbModel.Vault, apiModel.Vault)
		assert.Equal(t, strconv.FormatUint(dbModel.PeriodID, 10), apiModel.PeriodId)
		assert.Equal(t, dbModel.Twap.String(), apiModel.Twap)
		assert.Equal(t, strconv.FormatUint(dbModel.Dar, 10), apiModel.Dar)
		assert.Equal(t, dbModel.PriceBOverA.String(), *apiModel.PriceBOverA)

		// if the line below needs to be updated, add the field assertion above
		assert.Equal(t, reflect.TypeOf(apiModel).NumField(), 6)
	})

	t.Run("vaultPeriodModelsToAPI should return correct apiSpec model", func(t *testing.T) {
		dbModels := []*model.VaultPeriod{
			{},
			{},
		}
		apiModel := vaultPeriodModelsToAPI(dbModels)
		for i := range apiModel {
			assert.Equal(t, reflect.TypeOf(apiModel[i]).NumField(), 6)
		}
	})

	t.Run("tokenModelToApi should return correct apiSpec model", func(t *testing.T) {
		dbModel := model.Token{
			Pubkey:      solana.NewWallet().PublicKey().String(),
			Symbol:      utils.GetStringPtr("BTC"),
			Decimals:    0,
			IconURL:     utils.GetStringPtr("url"),
			CoinGeckoID: utils.GetStringPtr("wrapped-solana"),
		}
		apiModel := tokenModelToApi(&dbModel)
		assert.Equal(t, dbModel.Pubkey, apiModel.Pubkey)
		assert.Equal(t, dbModel.Symbol, apiModel.Symbol)
		assert.Equal(t, int(dbModel.Decimals), apiModel.Decimals)
		assert.Equal(t, dbModel.IconURL, apiModel.IconUrl)
		assert.Equal(t, dbModel.CoinGeckoID, apiModel.CoinGeckoId)

		// if the line below needs to be updated, add the field assertion above
		assert.Equal(t, reflect.TypeOf(apiModel).NumField(), 5)
	})

	t.Run("vaultPeriodModelsToAPI should return correct apiSpec model", func(t *testing.T) {
		dbModels := []*model.Token{
			{},
			{},
		}
		apiModel := tokenModelsToAPI(dbModels)
		for i := range apiModel {
			assert.Equal(t, reflect.TypeOf(apiModel[i]).NumField(), 5)
		}
	})

	t.Run("tokenAccountModelToAPI should return correct apiSpec model", func(t *testing.T) {
		dbModel := model.TokenAccount{
			Pubkey: solana.NewWallet().PublicKey().String(),
			Mint:   solana.NewWallet().PublicKey().String(),
			Owner:  solana.NewWallet().PublicKey().String(),
			Amount: 0,
			State:  "initialized",
		}
		apiModel := tokenAccountModelToAPI(&dbModel)
		assert.Equal(t, dbModel.Pubkey, apiModel.Pubkey)
		assert.Equal(t, dbModel.Mint, apiModel.Mint)
		assert.Equal(t, dbModel.Owner, apiModel.Owner)
		assert.Equal(t, strconv.FormatUint(dbModel.Amount, 10), apiModel.Amount)

		// if the line below needs to be updated, add the field assertion above
		assert.Equal(t, reflect.TypeOf(apiModel).NumField(), 5)
	})

	t.Run("protoConfigModelToAPI should return correct apiSpec model", func(t *testing.T) {
		dbModel := model.ProtoConfig{
			Pubkey:                  solana.NewWallet().PublicKey().String(),
			Granularity:             1,
			TokenADripTriggerSpread: 2,
			TokenBWithdrawalSpread:  3,
			Admin:                   solana.NewWallet().PublicKey().String(),
			TokenBReferralSpread:    4,
		}
		apiModel := protoConfigModelToAPI(&dbModel)
		assert.Equal(t, dbModel.Pubkey, apiModel.Pubkey)
		assert.Equal(t, dbModel.Admin, apiModel.Admin)
		assert.Equal(t, strconv.FormatUint(dbModel.Granularity, 10), apiModel.Granularity)
		assert.Equal(t, int(dbModel.TokenADripTriggerSpread), apiModel.TokenADripTriggerSpread)
		assert.Equal(t, int(dbModel.TokenBWithdrawalSpread), apiModel.TokenBWithdrawalSpread)
		assert.Equal(t, int(dbModel.TokenBReferralSpread), apiModel.TokenBReferralSpread)

		// if the line below needs to be updated, add the field assertion above
		assert.Equal(t, reflect.TypeOf(apiModel).NumField(), 6)
	})

	t.Run("protoConfigModelsToAPI should return correct apiSpec model", func(t *testing.T) {
		dbModels := []*model.ProtoConfig{
			{},
			{},
		}
		apiModel := protoConfigModelsToAPI(dbModels)
		for i := range apiModel {
			assert.Equal(t, reflect.TypeOf(apiModel[i]).NumField(), 6)
		}
	})

	t.Run("positionModelToAPI should return correct apiSpec model", func(t *testing.T) {
		dbModel := model.Position{
			Pubkey:                   solana.NewWallet().PublicKey().String(),
			Vault:                    solana.NewWallet().PublicKey().String(),
			Authority:                solana.NewWallet().PublicKey().String(),
			DepositedTokenAAmount:    1,
			WithdrawnTokenBAmount:    2,
			DepositTimestamp:         time.Unix(1667076855, 0),
			DcaPeriodIDBeforeDeposit: 3,
			NumberOfSwaps:            4,
			PeriodicDripAmount:       5,
			IsClosed:                 false,
		}
		apiModel := positionModelToAPI(&dbModel)
		assert.Equal(t, dbModel.Pubkey, apiModel.Pubkey)
		assert.Equal(t, dbModel.Vault, apiModel.Vault)
		assert.Equal(t, dbModel.Authority, apiModel.Authority)
		assert.Equal(t, strconv.FormatUint(dbModel.DepositedTokenAAmount, 10), apiModel.DepositedTokenAAmount)
		assert.Equal(t, strconv.FormatUint(dbModel.WithdrawnTokenBAmount, 10), apiModel.WithdrawnTokenBAmount)
		assert.Equal(t, strconv.FormatInt(dbModel.DepositTimestamp.Unix(), 10), apiModel.DepositTimestamp)
		assert.Equal(t, strconv.FormatUint(dbModel.DcaPeriodIDBeforeDeposit, 10), apiModel.DcaPeriodIdBeforeDeposit)
		assert.Equal(t, strconv.FormatUint(dbModel.NumberOfSwaps, 10), apiModel.NumberOfSwaps)
		assert.Equal(t, strconv.FormatUint(dbModel.PeriodicDripAmount, 10), apiModel.PeriodicDripAmount)
		assert.Equal(t, dbModel.IsClosed, apiModel.IsClosed)

		// if the line below needs to be updated, add the field assertion above
		assert.Equal(t, reflect.TypeOf(apiModel).NumField(), 10)
	})

	t.Run("positionModelsToAPI should return correct apiSpec model", func(t *testing.T) {
		dbModels := []*model.Position{
			{},
			{},
		}
		apiModel := positionModelsToAPI(dbModels)
		for i := range apiModel {
			assert.Equal(t, reflect.TypeOf(apiModel[i]).NumField(), 10)
		}
	})

	t.Run("activeWalletModelsToAPI should return correct apiSpec model", func(t *testing.T) {
		dbModels := []repository.ActiveWallet{
			{},
			{},
		}
		apiModel := activeWalletModelsToAPI(dbModels)
		for i := range apiModel {
			assert.Equal(t, reflect.TypeOf(apiModel[i]).NumField(), 2)
		}
	})

	t.Run("vaultTokenSwapToAPI should return correct apiSpec model", func(t *testing.T) {
		vault := model.Vault{
			DcaActivationTimestamp: time.Unix(1667076855, 0),
			Pubkey:                 solana.NewWallet().PublicKey().String(),
			ProtoConfig:            solana.NewWallet().PublicKey().String(),
			TokenAAccount:          solana.NewWallet().PublicKey().String(),
			TokenBAccount:          solana.NewWallet().PublicKey().String(),
			TreasuryTokenBAccount:  solana.NewWallet().PublicKey().String(),
			LastDcaPeriod:          10,
			DripAmount:             99,
			Enabled:                true,
			TokenPairID:            uuid.New().String(),
			MaxSlippageBps:         50,
			TokenAMint:             solana.NewWallet().PublicKey().String(),
			TokenBMint:             solana.NewWallet().PublicKey().String(),
		}
		tokenSwap := model.TokenSwap{
			Pubkey:        solana.NewWallet().PublicKey().String(),
			Mint:          solana.NewWallet().PublicKey().String(),
			Authority:     solana.NewWallet().PublicKey().String(),
			FeeAccount:    solana.NewWallet().PublicKey().String(),
			TokenAAccount: solana.NewWallet().PublicKey().String(),
			TokenBAccount: solana.NewWallet().PublicKey().String(),
			TokenPairID:   solana.NewWallet().PublicKey().String(),
			TokenAMint:    solana.NewWallet().PublicKey().String(),
			TokenBMint:    solana.NewWallet().PublicKey().String(),
			ID:            solana.NewWallet().PublicKey().String(),
		}

		apiModel := vaultTokenSwapToAPI(&vault, &tokenSwap)
		assert.Equal(t, tokenSwap.Pubkey, apiModel.Swap)
		assert.Equal(t, tokenSwap.Authority, apiModel.SwapAuthority)
		assert.Equal(t, tokenSwap.FeeAccount, apiModel.SwapFeeAccount)
		assert.Equal(t, tokenSwap.TokenAAccount, apiModel.SwapTokenAAccount)
		assert.Equal(t, tokenSwap.TokenBAccount, apiModel.SwapTokenBAccount)
		assert.Equal(t, tokenSwap.Mint, apiModel.SwapTokenMint)
		// if the line below needs to be updated, add the field assertion above
		assert.Equal(t, reflect.TypeOf(apiModel).NumField(), 7)

		assert.Equal(t, vault.TokenAMint, apiModel.TokenAMint)
		assert.Equal(t, vault.TokenBMint, apiModel.TokenBMint)
		assert.Equal(t, vault.Pubkey, apiModel.Vault)
		assert.Equal(t, vault.ProtoConfig, apiModel.VaultProtoConfig)
		assert.Equal(t, vault.TokenAAccount, apiModel.DripCommon.VaultTokenAAccount)
		assert.Equal(t, vault.TokenBAccount, apiModel.DripCommon.VaultTokenBAccount)
		// if the line below needs to be updated, add the field assertion above
		assert.Equal(t, reflect.TypeOf(apiModel.DripCommon).NumField(), 6)
	})

	t.Run("vaultWhirlpoolToAPI should return correct apiSpec model", func(t *testing.T) {
		vault := model.Vault{
			DcaActivationTimestamp: time.Unix(1667076855, 0),
			Pubkey:                 solana.NewWallet().PublicKey().String(),
			ProtoConfig:            solana.NewWallet().PublicKey().String(),
			TokenAAccount:          solana.NewWallet().PublicKey().String(),
			TokenBAccount:          solana.NewWallet().PublicKey().String(),
			TreasuryTokenBAccount:  solana.NewWallet().PublicKey().String(),
			LastDcaPeriod:          10,
			DripAmount:             99,
			Enabled:                true,
			TokenPairID:            uuid.New().String(),
			MaxSlippageBps:         50,
			TokenAMint:             solana.NewWallet().PublicKey().String(),
			TokenBMint:             solana.NewWallet().PublicKey().String(),
		}
		orcaWhirlpool := model.OrcaWhirlpool{
			Pubkey:                     solana.NewWallet().PublicKey().String(),
			WhirlpoolsConfig:           solana.NewWallet().PublicKey().String(),
			TokenMintA:                 solana.NewWallet().PublicKey().String(),
			TokenVaultA:                solana.NewWallet().PublicKey().String(),
			TokenMintB:                 solana.NewWallet().PublicKey().String(),
			TokenVaultB:                solana.NewWallet().PublicKey().String(),
			TickSpacing:                0,
			FeeRate:                    0,
			ProtocolFeeRate:            0,
			TickCurrentIndex:           0,
			ProtocolFeeOwedA:           decimal.Decimal{},
			ProtocolFeeOwedB:           decimal.Decimal{},
			RewardLastUpdatedTimestamp: decimal.Decimal{},
			Liquidity:                  decimal.Decimal{},
			SqrtPrice:                  decimal.Decimal{},
			FeeGrowthGlobalA:           decimal.Decimal{},
			FeeGrowthGlobalB:           decimal.Decimal{},
			TokenPairID:                "",
			ID:                         "",
		}

		apiModel := vaultWhirlpoolToAPI(&vault, &orcaWhirlpool)
		assert.Equal(t, orcaWhirlpool.Pubkey, apiModel.Whirlpool)
		assert.Equal(t, orcaWhirlpool.TokenVaultA, apiModel.TokenVaultA)
		assert.Equal(t, orcaWhirlpool.TokenVaultB, apiModel.TokenVaultB)
		assert.Equal(t, orcaWhirlpool.Oracle, apiModel.Oracle)
		// if the line below needs to be updated, add the field assertion above
		assert.Equal(t, reflect.TypeOf(apiModel).NumField(), 5)

		assert.Equal(t, vault.TokenAMint, apiModel.TokenAMint)
		assert.Equal(t, vault.TokenBMint, apiModel.TokenBMint)
		assert.Equal(t, vault.Pubkey, apiModel.Vault)
		assert.Equal(t, vault.ProtoConfig, apiModel.VaultProtoConfig)
		assert.Equal(t, vault.TokenAAccount, apiModel.DripCommon.VaultTokenAAccount)
		assert.Equal(t, vault.TokenBAccount, apiModel.DripCommon.VaultTokenBAccount)
		// if the line below needs to be updated, add the field assertion above
		assert.Equal(t, reflect.TypeOf(apiModel.DripCommon).NumField(), 6)
	})

}
