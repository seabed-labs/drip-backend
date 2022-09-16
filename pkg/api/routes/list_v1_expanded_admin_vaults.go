package controller

import (
	"net/http"
	"strconv"

	apispec "github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/dcaf-labs/drip/pkg/repository"
	"github.com/dcaf-labs/drip/pkg/repository/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type GetAdminVaultsExpandParams string

const (
	expandAll                  = GetAdminVaultsExpandParams("all")
	protoConfigValue           = GetAdminVaultsExpandParams("protoConfigValue")
	tokenAMintValue            = GetAdminVaultsExpandParams("tokenAMintValue")
	tokenBMintValue            = GetAdminVaultsExpandParams("tokenBMintValue")
	tokenAAccountValue         = GetAdminVaultsExpandParams("tokenAAccountValue")
	tokenBAccountValue         = GetAdminVaultsExpandParams("tokenBAccountValue")
	treasuryTokenBAccountValue = GetAdminVaultsExpandParams("treasuryTokenBAccountValue")
)

func (h Handler) GetV1AdminVaults(c echo.Context, params apispec.GetV1AdminVaultsParams) error {
	res := apispec.ListExpandedAdminVaults{}

	// Get all Vaults
	vaults, err := h.repo.AdminGetVaults(c.Request().Context(),
		repository.VaultFilterParams{
			IsEnabled:  (*bool)(params.Enabled),
			LikeTokenA: (*string)(params.TokenA),
			LikeTokenB: (*string)(params.TokenB),
			LikeVault:  (*string)(params.Vault),
		},
		repository.PaginationParams{
			Limit:  (*int)(params.Limit),
			Offset: (*int)(params.Offset),
		})
	if err != nil {
		logrus.WithError(err).Error("failed to get vaults")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "failed to get vaults as admin"})
	}
	// Get and Map all TokenPairs
	var tokenPairIDS []string
	for i := range vaults {
		vault := vaults[i]
		tokenPairIDS = append(tokenPairIDS, vault.TokenPairID)
	}
	tokenPairs, err := h.repo.GetTokenPairsByIDS(c.Request().Context(), tokenPairIDS)
	if err != nil {
		logrus.WithError(err).Error("failed to get tokenPairs")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "failed to get token pairs"})
	}
	tokenPairsByID := make(map[string]*model.TokenPair)
	for i := range tokenPairs {
		tokenPair := tokenPairs[i]
		tokenPairsByID[tokenPair.ID] = tokenPair
	}

	// Populate Base Result
	for i := range vaults {
		vault := vaults[i]
		tokenPair, ok := tokenPairsByID[vault.TokenPairID]
		if !ok {
			logrus.
				WithField("tokenPairID", vault.TokenPairID).
				WithField("vault", vault.Pubkey).
				Errorf("could not find token pair")
			return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
		}
		res = append(res, apispec.ExpandedAdminVault{
			Vault: apispec.Vault{
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
			},
			TokenAAccountValue:         nil,
			TokenBAccountValue:         nil,
			TreasuryTokenBAccountValue: nil,
			ProtoConfigValue:           nil,
			TokenAMintValue:            nil,
			TokenBMintValue:            nil,
		})
	}
	if params.Expand == nil {
		return c.JSON(http.StatusOK, res)
	}
	if hasValue(*params.Expand, string(expandAll)) {
		newParams := apispec.ExpandAdminVaultsQueryParam{string(protoConfigValue), string(tokenAMintValue), string(tokenBMintValue), string(tokenAAccountValue), string(tokenBAccountValue), string(treasuryTokenBAccountValue)}
		params.Expand = &newParams
	}

	// Prefetch data to make populating expand fields easier
	var protoConfigPubkeys []string
	var tokenAccountPubkeys []string
	for i := range vaults {
		vault := vaults[i]
		protoConfigPubkeys = append(protoConfigPubkeys, vault.ProtoConfig)
		tokenAccountPubkeys = append(tokenAccountPubkeys, vault.TokenAAccount)
		tokenAccountPubkeys = append(tokenAccountPubkeys, vault.TokenBAccount)
		tokenAccountPubkeys = append(tokenAccountPubkeys, vault.TreasuryTokenBAccount)
	}

	tokenAccountBalances, err := h.repo.GetTokenAccountBalancesByIDS(c.Request().Context(), tokenAccountPubkeys)
	if err != nil {
		logrus.WithError(err).Error("failed to get tokenAccountBalances")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	tokenAccountBalancesByPubkey := make(map[string]*model.TokenAccountBalance)
	for i := range tokenAccountBalances {
		tokeAccountBalance := tokenAccountBalances[i]
		tokenAccountBalancesByPubkey[tokeAccountBalance.Pubkey] = tokeAccountBalance
	}

	var tokenPubkeys []string
	for i := range tokenPairs {
		tokenPair := tokenPairs[i]
		tokenPubkeys = append(tokenPubkeys, tokenPair.TokenA)
		tokenPubkeys = append(tokenPubkeys, tokenPair.TokenB)
	}

	tokens, err := h.repo.GetTokensByMints(c.Request().Context(), tokenPubkeys)
	if err != nil {
		logrus.WithError(err).Error("failed to get tokenAccountBalances")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	tokensByPubkey := make(map[string]*model.Token)
	for i := range tokens {
		token := tokens[i]
		tokensByPubkey[token.Pubkey] = token
	}

	protoConfigs, err := h.repo.GetProtoConfigsByAddresses(c.Request().Context(), protoConfigPubkeys)
	if err != nil {
		logrus.WithError(err).Error("failed to get protoConfigs")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	protoConfigsByPubkey := make(map[string]*model.ProtoConfig)
	for i := range protoConfigs {
		protoConfig := protoConfigs[i]
		protoConfigsByPubkey[protoConfig.Pubkey] = protoConfig
	}

	for _, expandParam := range *params.Expand {
		switch expandParam {
		case string(protoConfigValue):
			for i := range res {
				protoConfig, ok := protoConfigsByPubkey[res[i].ProtoConfig]
				if !ok {
					logrus.
						WithField("vault", res[i].Vault).
						WithField("pubkey", res[i].ProtoConfig).
						Error("missing ProtoConfig")
					continue
				}
				// TODO(Mocha): unsafe cast
				res[i].ProtoConfigValue = &apispec.ProtoConfig{
					Pubkey:                  protoConfig.Pubkey,
					Admin:                   protoConfig.Admin,
					Granularity:             strconv.FormatUint(protoConfig.Granularity, 10),
					TokenADripTriggerSpread: int(protoConfig.TokenADripTriggerSpread),
					TokenBWithdrawalSpread:  int(protoConfig.TokenBWithdrawalSpread),
					TokenBReferralSpread:    int(protoConfig.TokenBReferralSpread),
				}
			}
		case string(tokenAMintValue):
			for i := range res {
				token, ok := tokensByPubkey[res[i].TokenAMint]
				if !ok {
					logrus.
						WithField("vault", res[i].Vault).
						WithField("pubkey", res[i].TokenAMint).
						Error("missing TokenAMint")
					continue
				}
				//TODO(Mocha): unsafe cast
				res[i].TokenAMintValue = &apispec.Token{
					Decimals: int(token.Decimals),
					Pubkey:   token.Pubkey,
					Symbol:   token.Symbol,
				}
			}
		case string(tokenBMintValue):
			for i := range res {
				token, ok := tokensByPubkey[res[i].TokenBMint]
				if !ok {
					logrus.
						WithField("vault", res[i].Vault).
						WithField("pubkey", res[i].TokenBMint).
						Error("missing TokenBMint")
					continue
				}
				res[i].TokenBMintValue = &apispec.Token{
					//TODO(Mocha): unsafe cast
					Decimals: int(token.Decimals),
					Pubkey:   token.Pubkey,
					Symbol:   token.Symbol,
				}
			}
		case string(tokenAAccountValue):
			for i := range res {
				tokenAccountBalance, ok := tokenAccountBalancesByPubkey[res[i].TokenAAccount]
				if !ok {
					logrus.
						WithField("vault", res[i].Vault).
						WithField("pubkey", res[i].TokenAAccount).
						Error("missing TokenAAccount")
					continue
				}
				// TODO(Mocha): Unsafe cast
				res[i].TokenAAccountValue = &apispec.TokenAccountBalance{
					Amount: strconv.FormatUint(tokenAccountBalance.Amount, 10),
					Mint:   tokenAccountBalance.Mint,
					Owner:  tokenAccountBalance.Owner,
					Pubkey: tokenAccountBalance.Pubkey,
					State:  tokenAccountBalance.State,
				}
			}
		case string(tokenBAccountValue):
			for i := range res {
				tokenAccountBalance, ok := tokenAccountBalancesByPubkey[res[i].TokenBAccount]
				if !ok {
					logrus.
						WithField("vault", res[i].Vault).
						WithField("pubkey", res[i].TokenBAccount).
						Error("missing TokenBAccount")
					continue
				}
				// TODO(Mocha): Unsafe cast
				res[i].TokenBAccountValue = &apispec.TokenAccountBalance{
					Amount: strconv.FormatUint(tokenAccountBalance.Amount, 10),
					Mint:   tokenAccountBalance.Mint,
					Owner:  tokenAccountBalance.Owner,
					Pubkey: tokenAccountBalance.Pubkey,
					State:  tokenAccountBalance.State,
				}
			}
		case string(treasuryTokenBAccountValue):
			for i := range res {
				tokenAccountBalance, ok := tokenAccountBalancesByPubkey[res[i].TreasuryTokenBAccount]
				if !ok {
					logrus.
						WithField("vault", res[i].Vault).
						WithField("pubkey", res[i].TreasuryTokenBAccount).
						Error("missing TreasuryTokenBAccount")
					continue
				}
				// TODO(Mocha): Unsafe cast
				res[i].TreasuryTokenBAccountValue = &apispec.TokenAccountBalance{
					Amount: strconv.FormatUint(tokenAccountBalance.Amount, 10),
					Mint:   tokenAccountBalance.Mint,
					Owner:  tokenAccountBalance.Owner,
					Pubkey: tokenAccountBalance.Pubkey,
					State:  tokenAccountBalance.State,
				}
			}
		}
	}

	return c.JSON(http.StatusOK, res)
}
