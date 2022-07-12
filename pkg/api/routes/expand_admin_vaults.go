package controller

import (
	"net/http"
	"strconv"

	"github.com/dcaf-labs/drip/pkg/repository/model"

	Swagger "github.com/dcaf-labs/drip/pkg/swagger"
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

func (h Handler) GetAdminVaults(c echo.Context, params Swagger.GetAdminVaultsParams) error {
	var res Swagger.ListExpandedAdminVaults

	// Get all Vaults
	vaults, err := h.repo.AdminGetVaults(c.Request().Context(), (*bool)(params.Enabled), (*int)(params.Limit), (*int)(params.Offset))
	if err != nil {
		logrus.WithError(err).Error("failed to get vaults")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "failed to get vaults as admin"})
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
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "failed to get token pairs"})
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
			return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal api error"})
		}
		res = append(res, Swagger.ExpandedAdminVault{
			Vault: Swagger.Vault{
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
		newParams := Swagger.ExpandAdminVaultsQueryParam{string(protoConfigValue), string(tokenAMintValue), string(tokenBMintValue), string(tokenAAccountValue), string(tokenBAccountValue), string(treasuryTokenBAccountValue)}
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
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal server error"})
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
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal server error"})
	}
	tokensByPubkey := make(map[string]*model.Token)
	for i := range tokens {
		token := tokens[i]
		tokensByPubkey[token.Pubkey] = token
	}

	protoConfigs, err := h.repo.GetProtoConfigsByPubkeys(c.Request().Context(), protoConfigPubkeys)
	if err != nil {
		logrus.WithError(err).Error("failed to get protoConfigs")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal server error"})
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
				res[i].ProtoConfigValue = &Swagger.ProtoConfig{
					BaseWithdrawalSpread: int(protoConfig.BaseWithdrawalSpread),
					Granularity:          int(protoConfig.Granularity),
					Pubkey:               protoConfig.Pubkey,
					TriggerDcaSpread:     int(protoConfig.TriggerDcaSpread),
				}
			}
			break
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
				res[i].TokenAMintValue = &Swagger.Token{
					Decimals: int(token.Decimals),
					Pubkey:   token.Pubkey,
					Symbol:   token.Symbol,
				}
			}
			break
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
				res[i].TokenBMintValue = &Swagger.Token{
					//TODO(Mocha): unsafe cast
					Decimals: int(token.Decimals),
					Pubkey:   token.Pubkey,
					Symbol:   token.Symbol,
				}
			}
			break
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
				res[i].TokenAAccountValue = &Swagger.TokenAccountBalance{
					Amount: int(tokenAccountBalance.Amount),
					Mint:   tokenAccountBalance.Mint,
					Owner:  tokenAccountBalance.Owner,
					Pubkey: tokenAccountBalance.Pubkey,
					State:  tokenAccountBalance.State,
				}
			}
			break
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
				res[i].TokenBAccountValue = &Swagger.TokenAccountBalance{
					Amount: int(tokenAccountBalance.Amount),
					Mint:   tokenAccountBalance.Mint,
					Owner:  tokenAccountBalance.Owner,
					Pubkey: tokenAccountBalance.Pubkey,
					State:  tokenAccountBalance.State,
				}
			}
			break
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
				res[i].TreasuryTokenBAccountValue = &Swagger.TokenAccountBalance{
					Amount: int(tokenAccountBalance.Amount),
					Mint:   tokenAccountBalance.Mint,
					Owner:  tokenAccountBalance.Owner,
					Pubkey: tokenAccountBalance.Pubkey,
					State:  tokenAccountBalance.State,
				}
			}
			break
		}
	}

	return c.JSON(http.StatusOK, res)
}

func hasValue(params Swagger.ExpandAdminVaultsQueryParam, value string) bool {
	for _, v := range params {
		if v == value {
			return true
		}
	}
	return false
}