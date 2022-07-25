package controller

import (
	"fmt"
	"net/http"

	"github.com/dcaf-labs/drip/pkg/repository"
	"github.com/dcaf-labs/drip/pkg/repository/model"
	Swagger "github.com/dcaf-labs/drip/pkg/swagger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetSpltokenswapconfigs(c echo.Context, params Swagger.GetSpltokenswapconfigsParams) error {
	var res Swagger.ListSplTokenSwapConfigs

	var vaults []*model.Vault
	if params.Vault != nil {
		vault, err := h.repo.GetVaultByAddress(c.Request().Context(), string(*params.Vault))
		if err != nil {
			logrus.WithError(err).WithField("vault", *params.Vault).Errorf("failed to get vault by address")
			return c.JSON(http.StatusBadRequest, Swagger.ErrorResponse{Error: "invalid vault address"})
		}
		vaults = []*model.Vault{vault}
	} else {
		var err error
		vaults, err = h.repo.GetVaultsWithFilter(c.Request().Context(), nil, nil, nil)
		if err != nil {
			logrus.WithError(err).WithField("vault", *params.Vault).Errorf("failed to get vaults")
			return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "failed to get vaults"})
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
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal api error"})
	}
	vaultWhitelistsByVaultPubkey := make(map[string][]*model.VaultWhitelist)
	for i := range vaultWhitelists {
		vaultWhitelist := vaultWhitelists[i]
		if _, ok := vaultWhitelistsByVaultPubkey[vaultWhitelist.VaultPubkey]; !ok {
			vaultWhitelistsByVaultPubkey[vaultWhitelist.VaultPubkey] = []*model.VaultWhitelist{}
		}
		vaultWhitelistsByVaultPubkey[vaultWhitelist.VaultPubkey] = append(vaultWhitelistsByVaultPubkey[vaultWhitelist.VaultPubkey], vaultWhitelist)
	}
	tokenSwaps, err := h.repo.GetTokenSwapsSortedByLiquidity(c.Request().Context(), tokenPairIDS)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get token swaps")
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: "internal api error"})
	}

	// TODO(Mocha): Return token swap with the most liquidity
	tokenSwapsByTokenPairID := make(map[string][]*repository.TokenSwapWithLiquidityRatio)
	for i := range tokenSwaps {
		tokenSwap := tokenSwaps[i]
		if _, ok := tokenSwapsByTokenPairID[tokenSwap.TokenPairID]; !ok {
			tokenSwapsByTokenPairID[tokenSwap.TokenPairID] = []*repository.TokenSwapWithLiquidityRatio{}
		}
		tokenSwapsByTokenPairID[tokenSwap.TokenPairID] = append(tokenSwapsByTokenPairID[tokenSwap.TokenPairID], &tokenSwap)
	}

	for i := range vaults {
		vault := vaults[i]
		tokenSwap, err := findTokenSwapForVault(vault, vaultWhitelistsByVaultPubkey, tokenSwapsByTokenPairID)
		if err != nil {
			logrus.WithError(err).Errorf("failed to get token swap for vault")
			continue
		}
		res = append(res, Swagger.SplTokenSwapConfig{
			Swap:              tokenSwap.Pubkey,
			SwapAuthority:     tokenSwap.Authority,
			SwapFeeAccount:    tokenSwap.FeeAccount,
			SwapTokenAAccount: tokenSwap.TokenAAccount,
			SwapTokenBAccount: tokenSwap.TokenBAccount,
			SwapTokenMint:     tokenSwap.Mint,
			DripCommon: Swagger.DripCommon{
				TokenAMint:         tokenSwap.TokenAMint,
				TokenBMint:         tokenSwap.TokenBMint,
				Vault:              vault.Pubkey,
				VaultProtoConfig:   vault.ProtoConfig,
				VaultTokenAAccount: vault.TokenAAccount,
				VaultTokenBAccount: vault.TokenBAccount,
			},
		})
	}
	return c.JSON(http.StatusOK, res)
}

func findTokenSwapForVault(vault *model.Vault, vaultWhitelistsByVaultPubkey map[string][]*model.VaultWhitelist, tokenSwapsByTokenPairID map[string][]*repository.TokenSwapWithLiquidityRatio) (*repository.TokenSwapWithLiquidityRatio, error) {
	tokenSwaps, ok := tokenSwapsByTokenPairID[vault.TokenPairID]
	if !ok {
		logrus.
			WithField("vault", vault.Pubkey).
			WithField("TokenPairID", vault.TokenPairID).
			Infof("skipping vault swap config, missing swap")
	}
	var eligibleSwaps []*repository.TokenSwapWithLiquidityRatio
	vaultWhitelists, ok := vaultWhitelistsByVaultPubkey[vault.Pubkey]
	if !ok || len(vaultWhitelists) == 0 {
		eligibleSwaps = tokenSwaps
	} else {
		for _, tokenSwap := range tokenSwaps {
			for _, vaultWhitelist := range vaultWhitelists {
				if vaultWhitelist.TokenSwapPubkey == tokenSwap.Pubkey {
					eligibleSwaps = append(eligibleSwaps, tokenSwap)
				}
			}
		}
	}

	// TODO(Mocha): Take into account swap fees
	if len(eligibleSwaps) == 0 {
		return nil, fmt.Errorf("failed to get token_swap")
	}
	bestSwap := eligibleSwaps[0]
	bestSwapDeltaB := evaluateTokenSwap(vault.DripAmount, bestSwap.TokenABalanceAmount, bestSwap.TokenBBalanceAmount)
	for _, eligibleSwap := range eligibleSwaps {
		swapDeltaB := evaluateTokenSwap(vault.DripAmount, eligibleSwap.TokenABalanceAmount, eligibleSwap.TokenBBalanceAmount)
		if swapDeltaB > bestSwapDeltaB {
			bestSwap = eligibleSwap
			bestSwapDeltaB = swapDeltaB
		}
	}

	return bestSwap, nil
}

// Calculates DeltaB from (reserveA + deltaA) * (reserveB - deltaB) = reserveA*reserveB =  k
// deltaB = reserveB - ((reserveA * reserveB) / (reservaA + deltaA))
// to be used to MAXIMIZE delta b across all swaps
func evaluateTokenSwap(deltaA, reserveA, reserveB uint64) uint64 {
	return reserveB - ((reserveA * reserveB) / (reserveA + deltaA))
}
