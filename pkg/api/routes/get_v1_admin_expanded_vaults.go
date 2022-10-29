package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type GetAdminVaultsExpandParams string

const (
	expandAll                             = GetAdminVaultsExpandParams("all")
	expandVaultProtoConfigValue           = GetAdminVaultsExpandParams("protoConfigValue")
	expandVaultTokenAMintValue            = GetAdminVaultsExpandParams("tokenAMintValue")
	expandVaultTokenBMintValue            = GetAdminVaultsExpandParams("tokenBMintValue")
	expandVaultTokenAAccountValue         = GetAdminVaultsExpandParams("tokenAAccountValue")
	expandVaultTokenBAccountValue         = GetAdminVaultsExpandParams("tokenBAccountValue")
	expandVaultTreasuryTokenBAccountValue = GetAdminVaultsExpandParams("treasuryTokenBAccountValue")
)

func (h Handler) GetV1AdminVaults(c echo.Context, params apispec.GetV1AdminVaultsParams) error {
	res := apispec.ListExpandedAdminVaults{}
	// Get  Vaults
	vaults, err := h.repo.AdminGetVaults(c.Request().Context(),
		repository.VaultFilterParams{
			IsEnabled:        (*bool)(params.Enabled),
			TokenA:           (*string)(params.TokenA),
			TokenB:           (*string)(params.TokenB),
			Vault:            (*string)(params.Vault),
			VaultProtoConfig: (*string)(params.VaultProtoConfig),
		},
		getPaginationParamsFromAPI(params.Offset, params.Limit),
	)
	if err != nil {
		logrus.WithError(err).Error("failed to get vaults")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: err.Error()})
	}

	// Populate Base Result
	for i := range vaults {
		res = append(res, apispec.ExpandedAdminVault{
			Vault: vaultModelToAPI(vaults[i]),
		})
	}
	if params.Expand == nil {
		return c.JSON(http.StatusOK, res)
	}
	if hasValue(*params.Expand, string(expandAll)) {
		newParams := apispec.ExpandAdminVaultsQueryParam{string(expandVaultProtoConfigValue), string(expandVaultTokenAMintValue), string(expandVaultTokenBMintValue), string(expandVaultTokenAAccountValue), string(expandVaultTokenBAccountValue), string(expandVaultTreasuryTokenBAccountValue)}
		params.Expand = &newParams
	}
	shouldExpandProtoConfig := hasValue(*params.Expand, string(expandVaultProtoConfigValue))
	shouldExpandTokenAMint := hasValue(*params.Expand, string(expandVaultTokenAMintValue))
	shouldExpandTokenBMint := hasValue(*params.Expand, string(expandVaultTokenBMintValue))
	shouldExpandTokenAAccount := hasValue(*params.Expand, string(expandVaultTokenAAccountValue))
	shouldExpandTokenBAccount := hasValue(*params.Expand, string(expandVaultTokenBAccountValue))
	shouldExpandTreasuryTokenBAccount := hasValue(*params.Expand, string(expandVaultTreasuryTokenBAccountValue))

	var protoConfigPubkeys []string
	var protoConfigs []*model.ProtoConfig
	protoConfigsByPubkey := make(map[string]*model.ProtoConfig)

	var tokenPubkeys []string
	var tokens []*model.Token
	tokensByPubkey := make(map[string]*model.Token)

	var tokenAccountPubkeys []string
	var tokenAccounts []*model.TokenAccount
	tokenAccountsByPubkey := make(map[string]*model.TokenAccount)

	for _, vault := range vaults {
		protoConfigPubkeys = append(protoConfigPubkeys, vault.ProtoConfig)
		tokenPubkeys = append(tokenPubkeys, vault.TokenAMint, vault.TokenBMint)
		tokenAccountPubkeys = append(tokenAccountPubkeys, vault.TokenAAccount, vault.TokenBAccount, vault.TreasuryTokenBAccount)
	}

	// fetch expanded data
	if shouldExpandProtoConfig {
		protoConfigs, err = h.repo.GetProtoConfigsByAddresses(c.Request().Context(), protoConfigPubkeys)
		if err != nil {
			logrus.WithError(err).Error("failed to GetProtoConfigsByAddresses")
			return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: err.Error()})
		}
		for i := range protoConfigs {
			protoConfigsByPubkey[protoConfigs[i].Pubkey] = protoConfigs[i]
		}
	}
	if shouldExpandTokenAMint || shouldExpandTokenBMint {
		tokens, err = h.repo.GetTokensByMints(c.Request().Context(), tokenPubkeys)
		if err != nil {
			logrus.WithError(err).Error("failed to GetTokensByMints")
			return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: err.Error()})
		}
		for i := range tokens {
			tokensByPubkey[tokens[i].Pubkey] = tokens[i]
		}
	}
	if shouldExpandTokenAAccount || shouldExpandTokenBAccount || shouldExpandTreasuryTokenBAccount {
		tokenAccounts, err = h.repo.GetTokenAccountsByAddresses(c.Request().Context(), tokenAccountPubkeys...)
		if err != nil {
			logrus.WithError(err).Error("failed to GetTokenAccountsByAddresses")
			return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: err.Error()})
		}
		for i := range tokenAccounts {
			tokenAccountsByPubkey[tokenAccounts[i].Pubkey] = tokenAccounts[i]
		}
	}

	// Populate expanded values
	for i := range res {
		vault := res[i]
		if protoConfig, ok := protoConfigsByPubkey[vault.ProtoConfig]; ok && shouldExpandProtoConfig {
			apiProtoConfig := protoConfigModelToAPI(protoConfig)
			res[i].ProtoConfigValue = &apiProtoConfig
		}
		if tokenAMint, ok := tokensByPubkey[vault.TokenAMint]; ok && shouldExpandTokenAMint {
			apiToken := tokenModelToApi(tokenAMint)
			res[i].TokenAMintValue = &apiToken
		}
		if tokenBMint, ok := tokensByPubkey[vault.TokenBMint]; ok && shouldExpandTokenBMint {
			apiToken := tokenModelToApi(tokenBMint)
			res[i].TokenBMintValue = &apiToken
		}
		if tokenAAccount, ok := tokenAccountsByPubkey[vault.TokenAAccount]; ok && shouldExpandTokenAAccount {
			apiTokenAccount := tokenAccountModelToAPI(tokenAAccount)
			res[i].TokenAAccountValue = &apiTokenAccount
		}
		if tokenBAccount, ok := tokenAccountsByPubkey[vault.TokenBAccount]; ok && shouldExpandTokenBAccount {
			apiTokenAccount := tokenAccountModelToAPI(tokenBAccount)
			res[i].TokenBAccountValue = &apiTokenAccount
		}
		if treasuryTokenBAccount, ok := tokenAccountsByPubkey[vault.TreasuryTokenBAccount]; ok && shouldExpandTreasuryTokenBAccount {
			apiTokenAccount := tokenAccountModelToAPI(treasuryTokenBAccount)
			res[i].TreasuryTokenBAccountValue = &apiTokenAccount
		}

	}

	return c.JSON(http.StatusOK, res)
}
