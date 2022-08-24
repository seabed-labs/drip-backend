package controller

import (
	"fmt"
	"net/http"

	"github.com/dcaf-labs/drip/pkg/repository"

	"github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/dcaf-labs/drip/pkg/configs"
	"github.com/dcaf-labs/drip/pkg/repository/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetOrcawhirlpoolconfigs(c echo.Context, params apispec.GetOrcawhirlpoolconfigsParams) error {
	res := apispec.ListOrcaWhirlpoolConfigs{}

	// TODO(Mocha): Refactor this and a the token swap config controller
	var vaults []*repository.VaultWithTokenPair
	if params.Vault != nil {
		vault, err := h.repo.GetVaultByAddress(c.Request().Context(), string(*params.Vault))
		if err != nil {
			logrus.WithError(err).WithField("vault", *params.Vault).Errorf("failed to get vault by address")
			return c.JSON(http.StatusBadRequest, apispec.ErrorResponse{Error: "invalid vault address"})
		}
		vaults = []*repository.VaultWithTokenPair{vault}
	} else {
		var err error
		vaults, err = h.repo.GetVaultsWithFilter(c.Request().Context(), repository.VaultFilterParams{})
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
	vaultWhitelistsByVaultPubkey := make(map[string][]*model.VaultWhitelist)
	for i := range vaultWhitelists {
		vaultWhitelist := vaultWhitelists[i]
		if _, ok := vaultWhitelistsByVaultPubkey[vaultWhitelist.VaultPubkey]; !ok {
			vaultWhitelistsByVaultPubkey[vaultWhitelist.VaultPubkey] = []*model.VaultWhitelist{}
		}
		vaultWhitelistsByVaultPubkey[vaultWhitelist.VaultPubkey] = append(vaultWhitelistsByVaultPubkey[vaultWhitelist.VaultPubkey], vaultWhitelist)
	}

	orcaWhirlpools, err := h.repo.GetOrcaWhirlpoolsByTokenPairIDs(c.Request().Context(), tokenPairIDS)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get orca whirlpools")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal api error"})
	}

	orcaWhirlpoolsByTokenPairID := make(map[string][]*model.OrcaWhirlpool)
	for i := range orcaWhirlpools {
		orcaWhirlpool := orcaWhirlpools[i]
		if _, ok := orcaWhirlpoolsByTokenPairID[orcaWhirlpool.TokenPairID]; !ok {
			orcaWhirlpoolsByTokenPairID[orcaWhirlpool.TokenPairID] = []*model.OrcaWhirlpool{}
		}
		orcaWhirlpoolsByTokenPairID[orcaWhirlpool.TokenPairID] = append(orcaWhirlpoolsByTokenPairID[orcaWhirlpool.TokenPairID], orcaWhirlpool)
	}

	for i := range vaults {
		vault := vaults[i]
		orcaWhirlpool, err := findOrcaWhirlpoolForVault(vault, vaultWhitelistsByVaultPubkey, orcaWhirlpoolsByTokenPairID, h.env)
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
	vault *repository.VaultWithTokenPair,
	vaultWhitelistsByVaultPubkey map[string][]*model.VaultWhitelist,
	orcaWhirlpoolsByTokenPairID map[string][]*model.OrcaWhirlpool,
	env configs.Environment,
) (*model.OrcaWhirlpool, error) {
	orcaWhirlpools, ok := orcaWhirlpoolsByTokenPairID[vault.TokenPairID]
	if !ok {
		logrus.
			WithField("vault", vault.Pubkey).
			WithField("TokenPairID", vault.TokenPairID).
			Infof("skipping vault swap config, missing swap")
	}
	// TODO: Remove
	if env == configs.MainnetEnv {
		var elgibleOrcaWhirlpools []*model.OrcaWhirlpool
		for _, orcaWhirlpool := range orcaWhirlpools {
			if orcaWhirlpool.Pubkey == "HJPjoWUrhoZzkNfRpHuieeFk9WcZWjwy6PBjZ81ngndJ" ||
				orcaWhirlpool.Pubkey == "ErSQss3jrqDpQoLEYvo6onzjsi6zm4Sjpoz1pjqz2o6D" ||
				orcaWhirlpool.Pubkey == "E5KuHFnU2VuuZFKeghbTLazgxeni4dhQ7URE4oBtJju2" {
				elgibleOrcaWhirlpools = append(elgibleOrcaWhirlpools, orcaWhirlpool)
			}
		}
		return elgibleOrcaWhirlpools[0], nil
	}

	var elgibleOrcaWhirlpools []*model.OrcaWhirlpool
	vaultWhitelists, ok := vaultWhitelistsByVaultPubkey[vault.Pubkey]
	if !ok || len(vaultWhitelists) == 0 {
		elgibleOrcaWhirlpools = orcaWhirlpools
	} else {
		for _, tokenSwap := range orcaWhirlpools {
			for _, vaultWhitelist := range vaultWhitelists {
				if vaultWhitelist.TokenSwapPubkey == tokenSwap.Pubkey {
					elgibleOrcaWhirlpools = append(elgibleOrcaWhirlpools, tokenSwap)
				}
			}
		}
	}

	if len(elgibleOrcaWhirlpools) == 0 {
		return nil, fmt.Errorf("failed to get orcaWhirlpool")
	}

	return elgibleOrcaWhirlpools[0], nil
}

//// TODO(Mocha): Figure how to get deltaB for an orca whirlpool

//bestSwap := eligibleSwaps[0]
//bestSwapDeltaB := evaluateTokenSwap(vault.DripAmount, bestSwap.TokenABalanceAmount, bestSwap.TokenBBalanceAmount)
//for _, eligibleSwap := range eligibleSwaps {
//	swapDeltaB := evaluateTokenSwap(vault.DripAmount, eligibleSwap.TokenABalanceAmount, eligibleSwap.TokenBBalanceAmount)
//	if swapDeltaB > bestSwapDeltaB {
//		bestSwap = eligibleSwap
//		bestSwapDeltaB = swapDeltaB
//	}
//}
// Calculates DeltaB from (reserveA + deltaA) * (reserveB - deltaB) = reserveA*reserveB =  k
// deltaB = reserveB - ((reserveA * reserveB) / (reservaA + deltaA))
// to be used to MAXIMIZE delta b across all swaps
//func evaluateOrcaWhirlpool(deltaA, reserveA, reserveB uint64) uint64 {
//	return reserveB - ((reserveA * reserveB) / (reserveA + deltaA))
//}
