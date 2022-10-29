package base

import (
	"context"
	"fmt"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	log "github.com/sirupsen/logrus"
)

func (i impl) GetBestTokenSwapsForVaults(
	ctx context.Context, vaults []*model.Vault,
) (map[string]*model.TokenSwap, error) {
	vaults = filterVaultsWithZeroDripAmount(vaults)
	vaultWhitelists, err := i.repo.GetVaultWhitelistsByVaultAddress(ctx, model.GetVaultPubkeys(vaults))
	if err != nil {
		return nil, err
	}
	tokenSwaps, err := i.repo.GetSPLTokenSwapsByTokenPairIDs(ctx, model.GetTokenPairIDsForVaults(vaults)...)
	if err != nil {
		return nil, err
	}
	tokenSwapAccounts, err := i.repo.GetTokenAccountsByAddresses(
		ctx, model.GetTokenAccountPubkeysForTokenSwaps(tokenSwaps)...)
	if err != nil {
		return nil, err
	}
	return getBestTokenSwapForVaults(
		vaults,
		model.GetVaultWhitelistsByVault(vaultWhitelists),
		model.GetTokenSwapsByTokenPairID(tokenSwaps),
		model.GetTokenAccountsByPubkey(tokenSwapAccounts)), nil
}

func getBestTokenSwapForVaults(
	vaults []*model.Vault,
	vaultWhitelists map[string][]*model.VaultWhitelist,
	tokenSwapsByTokenPairId map[string][]*model.TokenSwap,
	tokenAccountsByPubkey map[string]*model.TokenAccount,
) map[string]*model.TokenSwap {
	res := make(map[string]*model.TokenSwap)
	for i := range vaults {
		log.WithField("vault", vaults[i].Pubkey).WithField("tokenPairID", vaults[i].TokenPairID)
		tokenSwaps := filterNonWhitelistedTokenSwaps(
			tokenSwapsByTokenPairId[vaults[i].TokenPairID],
			vaultWhitelists[vaults[i].Pubkey])
		if len(tokenSwaps) == 0 {
			log.Errorf("no tokenSwaps found")
			continue
		}
		bestTokenSwap, err := getBestTokenSwapForVault(vaults[i], tokenSwaps, tokenAccountsByPubkey)
		if err != nil {
			log.WithError(err).Errorf("tokenSwaps found, but error in choosing the best one")
			continue
		}
		res[vaults[i].Pubkey] = bestTokenSwap
	}
	return res
}

func getBestTokenSwapForVault(
	vault *model.Vault,
	tokenSwaps []*model.TokenSwap,
	tokenAccountsByPubkey map[string]*model.TokenAccount,
) (*model.TokenSwap, error) {
	if len(tokenSwaps) == 0 {
		return nil, fmt.Errorf("failed to get token_swap")
	}
	bestSwapDeltaB := uint64(0)
	var bestSwap *model.TokenSwap = nil
	for _, eligibleSwap := range tokenSwaps {
		tokenAAccount, ok := tokenAccountsByPubkey[eligibleSwap.TokenAAccount]
		if !ok {
			return nil, fmt.Errorf("missing tokenAAcountBalance %s for tokenSwap %s", eligibleSwap.TokenAAccount, eligibleSwap.Pubkey)
		}
		tokenBAccount, ok := tokenAccountsByPubkey[eligibleSwap.TokenBAccount]
		if !ok {
			return nil, fmt.Errorf("missing tokenBAccount %s for tokenSwap %s", eligibleSwap.TokenBAccount, eligibleSwap.Pubkey)
		}
		swapDeltaB := evaluateTokenSwap(vault.DripAmount, tokenAAccount.Amount, tokenBAccount.Amount)

		if swapDeltaB > bestSwapDeltaB {
			bestSwap = eligibleSwap
			bestSwapDeltaB = swapDeltaB
		}
	}
	if bestSwap == nil {
		return nil, fmt.Errorf("failed to get bestSwap from list of %d tokenSwaps", len(tokenSwaps))
	}
	return bestSwap, nil
}

// evaluateTokenSwap Calculates DeltaB for a token swap
// (reserveA + deltaA) * (reserveB - deltaB) = reserveA*reserveB =  k
// deltaB = reserveB - ((reserveA * reserveB) / (reservaA + deltaA))
func evaluateTokenSwap(deltaA, reserveA, reserveB uint64) uint64 {
	return reserveB - ((reserveA * reserveB) / (reserveA + deltaA))
}

// filterNonWhitelistedTokenSwaps returns whitelisted swaps for a given vaultWhitelist
// assumption: `swaps` are valid tokenA/tokenB swaps for the vaults referenced in `vaultWhitelists`
func filterNonWhitelistedTokenSwaps(
	swaps []*model.TokenSwap,
	vaultWhitelists []*model.VaultWhitelist,
) []*model.TokenSwap {

	res := []*model.TokenSwap{}
	if len(swaps) == 0 || len(vaultWhitelists) == 0 {
		return swaps
	}
	whitelistBySwap := model.GetVaultWhitelistsBySwap(vaultWhitelists)
	for i := range swaps {
		if _, ok := whitelistBySwap[swaps[i].Pubkey]; ok {
			res = append(res, swaps[i])
		}
	}
	return res
}
