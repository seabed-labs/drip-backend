package controller

import (
	"net/http"
	"strconv"

	"github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/dcaf-labs/drip/pkg/repository"
	"github.com/dcaf-labs/drip/pkg/repository/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type GetAdminPositionsExpandParams string

const (
	expandPositionAll         = GetAdminPositionsExpandParams("all")
	expandPositionVault       = GetAdminPositionsExpandParams("vault")
	expandPositionProtoConfig = GetAdminPositionsExpandParams("protoConfig")
	expandPositionTokenA      = GetAdminPositionsExpandParams("tokenA")
	expandPositionTokenB      = GetAdminPositionsExpandParams("tokenB")
)

func (h Handler) GetV1AdminPositions(c echo.Context, params apispec.GetV1AdminPositionsParams) error {
	res := apispec.ListAdminPositions{}
	positions, err := h.repo.GetAdminPositions(
		c.Request().Context(),
		(*bool)(params.Enabled),
		repository.PositionFilterParams{
			IsClosed: (*bool)(params.IsClosed),
			Wallet:   nil,
		},
		repository.PaginationParams{
			Limit:  (*int)(params.Limit),
			Offset: (*int)(params.Offset),
		},
	)
	if err != nil {
		logrus.WithError(err).Error("failed to GetAdminPositions")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	for _, position := range positions {
		res = append(res, apispec.ExpandedAdminPosition{
			Position: apispec.Position{
				Authority:                position.Authority,
				DcaPeriodIdBeforeDeposit: strconv.FormatUint(position.DcaPeriodIDBeforeDeposit, 10),
				DepositTimestamp:         strconv.FormatInt(position.DepositTimestamp.Unix(), 10),
				DepositedTokenAAmount:    strconv.FormatUint(position.DepositedTokenAAmount, 10),
				IsClosed:                 false,
				NumberOfSwaps:            strconv.FormatUint(position.NumberOfSwaps, 10),
				PeriodicDripAmount:       strconv.FormatUint(position.PeriodicDripAmount, 10),
				Pubkey:                   position.Pubkey,
				Vault:                    position.Vault,
				WithdrawnTokenBAmount:    strconv.FormatUint(position.WithdrawnTokenBAmount, 10),
			},
		})
	}
	if params.Expand == nil {
		return c.JSON(http.StatusOK, res)
	}
	if hasValue(*params.Expand, string(expandPositionAll)) {
		newParams := apispec.ExpandAdminPositionsQueryParam{string(expandPositionVault), string(expandPositionProtoConfig), string(expandPositionTokenA), string(expandPositionTokenB)}
		params.Expand = &newParams
	}
	shouldExpandVault := hasValue(*params.Expand, string(expandPositionVault))
	shouldExpandProtoConfig := hasValue(*params.Expand, string(expandPositionProtoConfig))
	shouldExpandTokenA := hasValue(*params.Expand, string(expandPositionTokenA))
	shouldExpandTokenB := hasValue(*params.Expand, string(expandPositionTokenB))

	var vaultPubkeys []string
	for i := range positions {
		position := positions[i]
		vaultPubkeys = append(vaultPubkeys, position.Vault)
	}

	var vaults []*model.Vault
	vaultsByPubkey := make(map[string]*model.Vault)

	var protoConfigPubkeys []string
	var protoConfigs []*model.ProtoConfig
	protoConfigsByPubkey := make(map[string]*model.ProtoConfig)

	var tokenPairIDS []string
	var tokenPairs []*model.TokenPair
	tokenPairsByID := make(map[string]*model.TokenPair)

	var tokenPubkeys []string
	var tokens []*model.Token
	tokensByPubkey := make(map[string]*model.Token)

	vaults, err = h.repo.AdminGetVaultsByAddresses(c.Request().Context(), vaultPubkeys...)
	if err != nil {
		logrus.WithError(err).Error("failed to AdminGetVaultsByAddresses")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	for i := range vaults {
		vault := vaults[i]
		vaultsByPubkey[vault.Pubkey] = vault
		protoConfigPubkeys = append(protoConfigPubkeys, vault.ProtoConfig)
		tokenPairIDS = append(tokenPairIDS, vault.TokenPairID)
	}
	tokenPairs, err = h.repo.GetTokenPairsByIDS(c.Request().Context(), tokenPairIDS)
	if err != nil {
		logrus.WithError(err).Error("failed to GetTokenPairsByIDS")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	for i := range tokenPairs {
		tokenPairsByID[tokenPairs[i].ID] = tokenPairs[i]
		tokenPubkeys = append(tokenPubkeys, tokenPairs[i].TokenA, tokenPairs[i].TokenB)
	}

	if shouldExpandProtoConfig {
		protoConfigs, err = h.repo.GetProtoConfigsByAddresses(c.Request().Context(), protoConfigPubkeys)
		if err != nil {
			logrus.WithError(err).Error("failed to GetProtoConfigsByAddresses")
			return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
		}
		for i := range protoConfigs {
			protoConfigsByPubkey[protoConfigs[i].Pubkey] = protoConfigs[i]
		}
	}

	if shouldExpandTokenA || shouldExpandTokenB {
		tokens, err = h.repo.GetTokensByMints(c.Request().Context(), tokenPubkeys)
		if err != nil {
			logrus.WithError(err).Error("failed to GetTokensByTokenPairIDS")
			return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
		}
		for i := range tokens {
			tokensByPubkey[tokens[i].Pubkey] = tokens[i]
		}
	}
	for i := range res {
		var tokenPair *model.TokenPair
		vault := vaultsByPubkey[res[i].Position.Vault]
		if vault != nil {
			tokenPair = tokenPairsByID[vault.TokenPairID]
		}
		if shouldExpandVault && vault != nil && tokenPair != nil {
			res[i].Vault = &apispec.Vault{
				DcaActivationTimestamp: strconv.FormatInt(vault.DcaActivationTimestamp.Unix(), 10),
				DripAmount:             strconv.FormatUint(vault.DripAmount, 10),
				LastDcaPeriod:          strconv.FormatUint(vault.LastDcaPeriod, 10),
				ProtoConfig:            vault.ProtoConfig,
				Pubkey:                 vault.Pubkey,
				TokenAAccount:          vault.TokenAAccount,
				TokenAMint:             tokenPair.TokenA,
				TokenBAccount:          vault.TokenBAccount,
				TokenBMint:             tokenPair.TokenB,
				TreasuryTokenBAccount:  vault.TreasuryTokenBAccount,
				Enabled:                vault.Enabled,
			}
		}

		if shouldExpandTokenA && vault != nil && tokenPair != nil {
			tokenA := tokensByPubkey[tokenPair.TokenA]
			if tokenA != nil {
				res[i].TokenA = &apispec.Token{
					Decimals: int(tokenA.Decimals),
					Pubkey:   tokenA.Pubkey,
					Symbol:   tokenA.Symbol,
				}
			}
		}
		if shouldExpandTokenB && vault != nil && tokenPair != nil {
			tokenB := tokensByPubkey[tokenPair.TokenB]
			if tokenB != nil {
				res[i].TokenB = &apispec.Token{
					Decimals: int(tokenB.Decimals),
					Pubkey:   tokenB.Pubkey,
					Symbol:   tokenB.Symbol,
				}
			}
		}
		if shouldExpandProtoConfig && vault != nil {
			protoConfig := protoConfigsByPubkey[vault.ProtoConfig]
			res[i].ProtoConfig = &apispec.ProtoConfig{
				Admin:                   protoConfig.Admin,
				Granularity:             strconv.FormatUint(protoConfig.Granularity, 10),
				Pubkey:                  protoConfig.Pubkey,
				TokenADripTriggerSpread: int(protoConfig.TokenADripTriggerSpread),
				TokenBReferralSpread:    int(protoConfig.TokenBReferralSpread),
				TokenBWithdrawalSpread:  int(protoConfig.TokenBWithdrawalSpread),
			}
		}
	}

	return c.JSON(http.StatusOK, res)
}
