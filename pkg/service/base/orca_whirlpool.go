package base

import (
	"context"
	"fmt"

	"github.com/dcaf-labs/drip/pkg/service/configs"
	"github.com/dcaf-labs/drip/pkg/service/orcawhirlpool"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	log "github.com/sirupsen/logrus"
)

func (i impl) GetBestOrcaWhirlpoolForVaults(
	ctx context.Context, vaults []*model.Vault,
) (map[string]*model.OrcaWhirlpool, error) {
	vaults = filterVaultsWithZeroDripAmount(vaults)
	vaultWhitelists, err := i.repo.GetVaultWhitelistsByVaultAddress(ctx, model.GetVaultPubkeys(vaults))
	if err != nil {
		return nil, err
	}
	whirlpools, err := i.repo.GetOrcaWhirlpoolsByTokenPairIDs(ctx, model.GetTokenPairIDsForVaults(vaults)...)
	if err != nil {
		return nil, err
	}
	return getBestOrcaWhirlpoolForVaults(
		vaults,
		model.GetVaultWhitelistsByVault(vaultWhitelists),
		model.GetOrcaWhirlpoolsByTokenPairID(whirlpools),
		i.network), nil
}

func getBestOrcaWhirlpoolForVaults(
	vaults []*model.Vault,
	vaultWhitelists map[string][]*model.VaultWhitelist,
	orcaWhirlpoolByTokenPairID map[string][]*model.OrcaWhirlpool,
	network configs.Network,
) map[string]*model.OrcaWhirlpool {
	res := make(map[string]*model.OrcaWhirlpool)
	for i := range vaults {
		log.WithField("vault", vaults[i].Pubkey).WithField("tokenPairID", vaults[i].TokenPairID)
		whirlpools := filterNonWhitelistedOrcaWhirlpools(
			orcaWhirlpoolByTokenPairID[vaults[i].TokenPairID],
			vaultWhitelists[vaults[i].Pubkey])
		if len(whirlpools) == 0 {
			log.Errorf("no whirlpools found")
			continue
		}
		bestWhirlpool, err := getBestOrcaWhirlpoolForVault(vaults[i], whirlpools, network)
		if err != nil {
			log.WithError(err).Errorf("whirlpools found, but error in choosing the best one")
			continue
		}
		res[vaults[i].Pubkey] = bestWhirlpool
	}
	return res
}

func getBestOrcaWhirlpoolForVault(
	vault *model.Vault,
	whirlpools []*model.OrcaWhirlpool,
	network configs.Network,
) (*model.OrcaWhirlpool, error) {
	if len(whirlpools) == 0 {
		return nil, fmt.Errorf("failed to get token_swap")
	}
	bestSwapDeltaB := uint64(0)
	var bestSwap *model.OrcaWhirlpool
	for _, eligibleSwap := range whirlpools {
		// todo: this is inefficient, should we just calculate and cache this value on whirlpool state changes?
		swapDeltaB, err := orcawhirlpool.EvaluateOrcaWhirlpool(eligibleSwap.Pubkey, vault, network)
		if err != nil {
			return nil, err
		}

		if swapDeltaB > bestSwapDeltaB {
			bestSwap = eligibleSwap
			bestSwapDeltaB = swapDeltaB
		}
	}
	if bestSwap == nil {
		return nil, fmt.Errorf("failed to get bestSwap from list of %d whirlpools", len(whirlpools))
	}
	return bestSwap, nil
}

// filterNonWhitelistedOrcaWhirlpools returns whitelisted swaps for a given vaultWhitelist
// assumption: `swaps` are valid tokenA/tokenB swaps for the vaults referenced in `vaultWhitelists`
func filterNonWhitelistedOrcaWhirlpools(
	swaps []*model.OrcaWhirlpool,
	vaultWhitelists []*model.VaultWhitelist,
) []*model.OrcaWhirlpool {
	res := []*model.OrcaWhirlpool{}
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
