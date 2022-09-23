package controller

import (
	"fmt"
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/configs"

	model2 "github.com/dcaf-labs/drip/pkg/service/repository/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetV1DripOrcawhirlpoolconfigs(c echo.Context, params apispec.GetV1DripOrcawhirlpoolconfigsParams) error {
	res := apispec.ListOrcaWhirlpoolConfigs{}

	// TODO(Mocha): Refactor this and a the token swap config controller
	var vaults []*model2.Vault
	if params.Vault != nil {
		vault, err := h.repo.GetVaultByAddress(c.Request().Context(), string(*params.Vault))
		if err != nil {
			logrus.WithError(err).WithField("vault", *params.Vault).Errorf("failed to get vault by address")
			return c.JSON(http.StatusBadRequest, apispec.ErrorResponse{Error: "invalid vault address"})
		}
		vaults = []*model2.Vault{vault}
	} else {
		var err error
		vaults, err = h.repo.GetVaultsWithFilter(c.Request().Context(), nil, nil, nil)
		if err != nil {
			logrus.WithError(err).WithField("vault", *params.Vault).Errorf("failed to get vaults")
			return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "failed to get vaults"})
		}
	}

	var tokenPairIDS []string
	var vaultPubkeys []string
	for i := range vaults {
		vault := vaults[i]
		tokenPairIDS = append(tokenPairIDS, vault.TokenPairID)
		vaultPubkeys = append(vaultPubkeys, vault.Pubkey)
	}
	vaultWhitelists, err := h.repo.GetVaultWhitelistsByVaultAddress(c.Request().Context(), vaultPubkeys)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get vault whitelists")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}
	vaultWhitelistsByVaultPubkey := make(map[string][]*model2.VaultWhitelist)
	for i := range vaultWhitelists {
		vaultWhitelist := vaultWhitelists[i]
		if _, ok := vaultWhitelistsByVaultPubkey[vaultWhitelist.VaultPubkey]; !ok {
			vaultWhitelistsByVaultPubkey[vaultWhitelist.VaultPubkey] = []*model2.VaultWhitelist{}
		}
		vaultWhitelistsByVaultPubkey[vaultWhitelist.VaultPubkey] = append(vaultWhitelistsByVaultPubkey[vaultWhitelist.VaultPubkey], vaultWhitelist)
	}

	orcaWhirlpools, err := h.repo.GetOrcaWhirlpoolsByTokenPairIDs(c.Request().Context(), tokenPairIDS)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get orca whirlpools")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}

	orcaWhirlpoolsByTokenPairID := make(map[string][]*model2.OrcaWhirlpool)
	for i := range orcaWhirlpools {
		orcaWhirlpool := orcaWhirlpools[i]
		if _, ok := orcaWhirlpoolsByTokenPairID[orcaWhirlpool.TokenPairID]; !ok {
			orcaWhirlpoolsByTokenPairID[orcaWhirlpool.TokenPairID] = []*model2.OrcaWhirlpool{}
		}
		orcaWhirlpoolsByTokenPairID[orcaWhirlpool.TokenPairID] = append(orcaWhirlpoolsByTokenPairID[orcaWhirlpool.TokenPairID], orcaWhirlpool)
	}

	for i := range vaults {
		vault := vaults[i]
		orcaWhirlpool, err := findOrcaWhirlpoolForVault(vault, vaultWhitelistsByVaultPubkey, orcaWhirlpoolsByTokenPairID, h.network)
		if err != nil {
			logrus.WithError(err).Errorf("failed to get orca whirlpool for vault")
			continue
		}
		res = append(res, apispec.OrcaWhirlpoolConfig{
			Oracle:      orcaWhirlpool.Oracle,
			TokenVaultA: orcaWhirlpool.TokenVaultA,
			TokenVaultB: orcaWhirlpool.TokenVaultB,
			Whirlpool:   orcaWhirlpool.Pubkey,
			DripCommon: apispec.DripCommon{
				TokenAMint:         orcaWhirlpool.TokenMintA,
				TokenBMint:         orcaWhirlpool.TokenMintB,
				Vault:              vault.Pubkey,
				VaultProtoConfig:   vault.ProtoConfig,
				VaultTokenAAccount: vault.TokenAAccount,
				VaultTokenBAccount: vault.TokenBAccount,
			},
		})
	}
	return c.JSON(http.StatusOK, res)
}

func findOrcaWhirlpoolForVault(
	vault *model2.Vault,
	vaultWhitelistsByVaultPubkey map[string][]*model2.VaultWhitelist,
	orcaWhirlpoolsByTokenPairID map[string][]*model2.OrcaWhirlpool,
	network configs.Network,
) (*model2.OrcaWhirlpool, error) {
	orcaWhirlpools, ok := orcaWhirlpoolsByTokenPairID[vault.TokenPairID]
	if !ok {
		logrus.
			WithField("vault", vault.Pubkey).
			WithField("TokenPairID", vault.TokenPairID).
			Infof("skipping vault swap config, missing swap")
	}

	var elgibleOrcaWhirlpools []*model2.OrcaWhirlpool
	vaultWhitelists, ok := vaultWhitelistsByVaultPubkey[vault.Pubkey]
	if !ok || len(vaultWhitelists) == 0 {
		elgibleOrcaWhirlpools = orcaWhirlpools
	} else {
		for _, swap := range orcaWhirlpools {
			for _, vaultWhitelist := range vaultWhitelists {
				if vaultWhitelist.TokenSwapPubkey == swap.Pubkey {
					elgibleOrcaWhirlpools = append(elgibleOrcaWhirlpools, swap)
				}
			}
		}
	}

	// TODO: Remove
	if network == configs.MainnetNetwork {
		var tempElgibleOrcaWhirlpools []*model2.OrcaWhirlpool
		for _, orcaWhirlpool := range elgibleOrcaWhirlpools {
			if _, ok := mainnetOrcaWhirlpoolsMap[orcaWhirlpool.Pubkey]; ok {
				tempElgibleOrcaWhirlpools = append(tempElgibleOrcaWhirlpools, orcaWhirlpool)
			}
		}
		elgibleOrcaWhirlpools = tempElgibleOrcaWhirlpools
	}

	if len(elgibleOrcaWhirlpools) == 0 {
		return nil, fmt.Errorf("failed to get orcaWhirlpool")
	}

	return elgibleOrcaWhirlpools[0], nil
}
