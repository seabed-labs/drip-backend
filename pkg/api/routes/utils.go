package controller

import (
	"strconv"

	"github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/dcaf-labs/drip/pkg/repository"
	"github.com/dcaf-labs/drip/pkg/repository/model"
)

func vaultProtoConfigDataBaseModelToAPIModel(vaultProtoConfigs []*model.ProtoConfig) []apispec.ProtoConfig {
	res := []apispec.ProtoConfig{}
	for _, protoConfig := range vaultProtoConfigs {
		res = append(res, apispec.ProtoConfig{
			Pubkey:               protoConfig.Pubkey,
			BaseWithdrawalSpread: int(protoConfig.BaseWithdrawalSpread),
			Granularity:          strconv.FormatUint(protoConfig.Granularity, 10),
			TriggerDcaSpread:     int(protoConfig.TriggerDcaSpread),
		})
	}
	return res
}

func vaultPeriodDatabaseModelToAPIModel(vaultPeriods []*model.VaultPeriod) []apispec.VaultPeriod {
	res := []apispec.VaultPeriod{}
	for _, vaultPeriodDBModel := range vaultPeriods {
		res = append(res, apispec.VaultPeriod{
			Pubkey:   vaultPeriodDBModel.Pubkey,
			Vault:    vaultPeriodDBModel.Vault,
			PeriodId: strconv.FormatUint(vaultPeriodDBModel.PeriodID, 10),
			Twap:     vaultPeriodDBModel.Twap.String(),
			Dar:      strconv.FormatUint(vaultPeriodDBModel.Dar, 10),
		})
	}
	return res
}

func vaultWithTokenPairDatabaseModelToAPIModel(vaults []*repository.VaultWithTokenPair) []apispec.Vault {
	res := []apispec.Vault{}
	for _, vault := range vaults {
		res = append(res, apispec.Vault{
			DcaActivationTimestamp: strconv.FormatInt(vault.DcaActivationTimestamp.Unix(), 10),
			DripAmount:             strconv.FormatUint(vault.DripAmount, 10),
			LastDcaPeriod:          strconv.FormatUint(vault.LastDcaPeriod, 10),
			ProtoConfig:            vault.ProtoConfig,
			Pubkey:                 vault.Pubkey,
			TokenAAccount:          vault.TokenAAccount,
			TokenAMint:             vault.TokenA,
			TokenBAccount:          vault.TokenBAccount,
			TokenBMint:             vault.TokenB,
			TreasuryTokenBAccount:  vault.TreasuryTokenBAccount,
			Enabled:                vault.Enabled,
		})
	}
	return res
}

func tokenPairDatabaseModelToAPIModel(tokenPairs []*model.TokenPair) []apispec.TokenPair {
	res := []apispec.TokenPair{}
	for _, tokenPair := range tokenPairs {
		res = append(res, apispec.TokenPair{
			Id:     tokenPair.ID,
			TokenA: tokenPair.TokenA,
			TokenB: tokenPair.TokenB,
		})
	}
	return res
}
