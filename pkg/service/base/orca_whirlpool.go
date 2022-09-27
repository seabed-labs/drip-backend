package base

import (
	"context"
	"fmt"

	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func (i impl) GetBestOrcaWhirlpoolForVaults(
	ctx context.Context, vaults []*model.Vault,
) (map[string]*model.OrcaWhirlpool, error) {
	vaults = filterVaultsWithZeroDripAmount(vaults)
	vaultPubkeys := model.GetVaultPubkeys(vaults)
	vaultWhitelists, err := i.repo.GetVaultWhitelistsByVaultAddress(ctx, vaultPubkeys)
	if err != nil {
		return nil, err
	}
	whirlpoolDeltaBQuotes, err := i.repo.GetOrcaWhirlpoolDeltaBQuoteByVaultAddresses(ctx, vaultPubkeys...)
	if err != nil && err.Error() != repository.ErrRecordNotFound {
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
		model.GetOrcaWhirlpoolDeltaBQuoteByCompositeKey(whirlpoolDeltaBQuotes)), nil
}

func getBestOrcaWhirlpoolForVaults(
	vaults []*model.Vault,
	vaultWhitelists map[string][]*model.VaultWhitelist,
	orcaWhirlpoolByTokenPairID map[string][]*model.OrcaWhirlpool,
	orcaWhirlpoolDeltaBQuoteByCompositeKey map[string]*model.OrcaWhirlpoolDeltaBQuote,
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
		bestWhirlpool, err := getBestOrcaWhirlpoolForVault(vaults[i], whirlpools, orcaWhirlpoolDeltaBQuoteByCompositeKey)
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
	orcaWhirlpoolDeltaBQuoteByCompositeKey map[string]*model.OrcaWhirlpoolDeltaBQuote,
) (*model.OrcaWhirlpool, error) {
	if len(whirlpools) == 0 {
		return nil, fmt.Errorf("failed to get whirlpool")
	}
	bestSwapDeltaB := uint64(0)
	var bestSwap *model.OrcaWhirlpool
	for _, eligibleSwap := range whirlpools {
		compositeKey := vault.Pubkey + eligibleSwap.Pubkey
		whirlpoolDeltaBQuote, ok := orcaWhirlpoolDeltaBQuoteByCompositeKey[compositeKey]
		if !ok || whirlpoolDeltaBQuote == nil {
			logrus.
				WithField("vault", vault.Pubkey).
				WithField("whirlpool", eligibleSwap.Pubkey).
				Error("missing orca whirlpool deltaB estimate ")
			whirlpoolDeltaBQuote = &model.OrcaWhirlpoolDeltaBQuote{
				DeltaB:          uint64(0),
				VaultPubkey:     vault.Pubkey,
				WhirlpoolPubkey: eligibleSwap.Pubkey,
				TokenPairID:     vault.TokenPairID,
			}
		}
		if whirlpoolDeltaBQuote.DeltaB >= bestSwapDeltaB {
			bestSwap = eligibleSwap
			bestSwapDeltaB = whirlpoolDeltaBQuote.DeltaB
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
