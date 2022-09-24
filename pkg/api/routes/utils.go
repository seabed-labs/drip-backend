package controller

import (
	"strconv"

	"github.com/dcaf-labs/drip/pkg/service/repository"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
)

func hasValue(params []string, value string) bool {
	for _, v := range params {
		if v == value {
			return true
		}
	}
	return false
}

func vaultModelToAPI(vaultModel *model.Vault) apispec.Vault {
	return apispec.Vault{
		DcaActivationTimestamp: strconv.FormatInt(vaultModel.DcaActivationTimestamp.Unix(), 10),
		DripAmount:             strconv.FormatUint(vaultModel.DripAmount, 10),
		LastDcaPeriod:          strconv.FormatUint(vaultModel.LastDcaPeriod, 10),
		ProtoConfig:            vaultModel.ProtoConfig,
		Pubkey:                 vaultModel.Pubkey,
		TokenAAccount:          vaultModel.TokenAAccount,
		TokenAMint:             vaultModel.TokenAMint,
		TokenBAccount:          vaultModel.TokenBAccount,
		TokenBMint:             vaultModel.TokenBMint,
		TreasuryTokenBAccount:  vaultModel.TreasuryTokenBAccount,
		Enabled:                vaultModel.Enabled,
	}
}

func vaultModelsToAPI(vaultModels []*model.Vault) apispec.ListVaults {
	res := apispec.ListVaults{}
	for i := range vaultModels {
		res = append(res, vaultModelToAPI(vaultModels[i]))
	}
	return res
}

func vaultPeriodModelsToAPI(vaultPeriodModels []*model.VaultPeriod) apispec.ListVaultPeriods {
	res := apispec.ListVaultPeriods{}
	for i := range vaultPeriodModels {
		vaultPeriod := vaultPeriodModels[i]
		res = append(res, struct {
			Dar      string `json:"dar"`
			PeriodId string `json:"periodId"`
			Pubkey   string `json:"pubkey"`
			Twap     string `json:"twap"`
			Vault    string `json:"vault"`
		}{
			Pubkey:   vaultPeriod.Pubkey,
			Vault:    vaultPeriod.Vault,
			PeriodId: strconv.FormatUint(vaultPeriod.PeriodID, 10),
			Twap:     vaultPeriod.Twap.String(),
			Dar:      strconv.FormatUint(vaultPeriod.Dar, 10),
		},
		)
	}
	return res
}

func tokenModelToApi(tokenModel *model.Token) apispec.Token {
	return apispec.Token{
		Decimals: int(tokenModel.Decimals),
		Pubkey:   tokenModel.Pubkey,
		Symbol:   tokenModel.Symbol,
		IconUrl:  tokenModel.IconURL,
	}
}

func tokenModelsToAPI(tokenModels []*model.Token) apispec.ListTokens {
	res := apispec.ListTokens{}
	for i := range tokenModels {
		res = append(res, tokenModelToApi(tokenModels[i]))
	}
	return res
}

func protoConfigModelToAPI(protoConfigModel *model.ProtoConfig) apispec.ProtoConfig {
	return apispec.ProtoConfig{
		Pubkey:                  protoConfigModel.Pubkey,
		Admin:                   protoConfigModel.Admin,
		Granularity:             strconv.FormatUint(protoConfigModel.Granularity, 10),
		TokenADripTriggerSpread: int(protoConfigModel.TokenADripTriggerSpread),
		TokenBWithdrawalSpread:  int(protoConfigModel.TokenBWithdrawalSpread),
		TokenBReferralSpread:    int(protoConfigModel.TokenBReferralSpread),
	}
}

func protoConfigModelsToAPI(protoConfigModels []*model.ProtoConfig) apispec.ListProtoConfigs {
	res := apispec.ListProtoConfigs{}
	for i := range protoConfigModels {
		res = append(res, protoConfigModelToAPI(protoConfigModels[i]))
	}
	return res
}

func positionModelToAPI(positionModel *model.Position) apispec.Position {
	return apispec.Position{
		Authority:                positionModel.Authority,
		DcaPeriodIdBeforeDeposit: strconv.FormatUint(positionModel.DcaPeriodIDBeforeDeposit, 10),
		DepositTimestamp:         strconv.FormatInt(positionModel.DepositTimestamp.Unix(), 10),
		DepositedTokenAAmount:    strconv.FormatUint(positionModel.DepositedTokenAAmount, 10),
		IsClosed:                 positionModel.IsClosed,
		NumberOfSwaps:            strconv.FormatUint(positionModel.NumberOfSwaps, 10),
		PeriodicDripAmount:       strconv.FormatUint(positionModel.PeriodicDripAmount, 10),
		Pubkey:                   positionModel.Pubkey,
		Vault:                    positionModel.Vault,
		WithdrawnTokenBAmount:    strconv.FormatUint(positionModel.WithdrawnTokenBAmount, 10),
	}
}

func positionModelsToAPI(positionModels []*model.Position) apispec.ListPositions {
	res := apispec.ListPositions{}
	for i := range positionModels {
		res = append(res, positionModelToAPI(positionModels[i]))
	}
	return res
}

func activeWalletModelsToAPI(activeWallets []repository.ActiveWallet) apispec.ListActiveWallets {
	res := apispec.ListActiveWallets{}
	for _, activeWallet := range activeWallets {
		res = append(res, apispec.ActiveWallet{
			Owner:         activeWallet.Owner,
			PositionCount: activeWallet.PositionCount,
		})
	}
	return res
}

func getPaginationParamsFromAPI(offsetParam *apispec.OffsetQueryParam, limitParam *apispec.LimitQueryParam) repository.PaginationParams {
	limit := defaultLimit
	if limitParam != nil {
		limit = int(*limitParam)
	}
	var offset int
	if offsetParam != nil {
		offset = int(*offsetParam)
	}
	return repository.PaginationParams{
		Limit:  &limit,
		Offset: &offset,
	}
}
