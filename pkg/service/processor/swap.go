package processor

import (
	"context"
	"strconv"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	"github.com/dcaf-labs/solana-go-clients/pkg/tokenswap"
	"github.com/dcaf-labs/solana-go-clients/pkg/whirlpool"
	solana2 "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (p impl) UpsertTokenSwapByAddress(ctx context.Context, address string) error {
	var tokenSwap tokenswap.TokenSwap
	if err := p.solanaClient.GetAccount(ctx, address, &tokenSwap); err != nil {
		return err
	}
	var tokenLPMint token.Mint
	if err := p.solanaClient.GetAccount(ctx, tokenSwap.TokenPool.String(), &tokenLPMint); err != nil {
		return err
	}

	// Add swap A -> B
	tokenPair, err := p.ensureTokenPair(ctx, tokenSwap.MintA.String(), tokenSwap.MintB.String())
	if err != nil {
		return err
	}
	if err := p.repo.UpsertTokenSwaps(ctx, &model.TokenSwap{
		ID:            uuid.New().String(),
		Pubkey:        address,
		Mint:          tokenSwap.TokenPool.String(),
		Authority:     tokenLPMint.MintAuthority.String(),
		FeeAccount:    tokenSwap.FeeAccount.String(),
		TokenAMint:    tokenSwap.MintA.String(),
		TokenAAccount: tokenSwap.TokenAccountA.String(),
		TokenBMint:    tokenSwap.MintB.String(),
		TokenBAccount: tokenSwap.TokenAccountB.String(),
		TokenPairID:   tokenPair.ID,
	}); err != nil {
		return err
	}
	// Add swap B -> A
	tokenPairInverse, err := p.ensureTokenPair(ctx, tokenSwap.MintB.String(), tokenSwap.MintA.String())
	if err != nil {
		return err
	}
	if err := p.repo.UpsertTokenSwaps(ctx, &model.TokenSwap{
		ID:            uuid.New().String(),
		Pubkey:        address,
		Mint:          tokenSwap.TokenPool.String(),
		Authority:     tokenLPMint.MintAuthority.String(),
		FeeAccount:    tokenSwap.FeeAccount.String(),
		TokenAMint:    tokenSwap.MintB.String(),
		TokenAAccount: tokenSwap.TokenAccountB.String(),
		TokenBMint:    tokenSwap.MintA.String(),
		TokenBAccount: tokenSwap.TokenAccountA.String(),
		TokenPairID:   tokenPairInverse.ID,
	}); err != nil {
		return err
	}

	// Upsert balances
	if err := p.UpsertTokenAccountBalanceByAddress(ctx, tokenSwap.TokenAccountA.String()); err != nil {
		return err
	}
	if err := p.UpsertTokenAccountBalanceByAddress(ctx, tokenSwap.TokenAccountB.String()); err != nil {
		return err
	}
	return nil
}

func (p impl) UpsertWhirlpoolByAddress(ctx context.Context, address string) error {
	var orcaWhirlpool whirlpool.Whirlpool
	whirlpoolPubkey := solana2.MustPublicKeyFromBase58(address)
	if err := p.solanaClient.GetAccount(ctx, address, &orcaWhirlpool); err != nil {
		return err
	}

	tokenPair, err := p.ensureTokenPair(ctx, orcaWhirlpool.TokenMintA.String(), orcaWhirlpool.TokenMintB.String())
	if err != nil {
		return err
	}
	inverseTokenPair, err := p.ensureTokenPair(ctx, orcaWhirlpool.TokenMintB.String(), orcaWhirlpool.TokenMintA.String())
	if err != nil {
		return err
	}

	protocolFeeOwedA, err := decimal.NewFromString(strconv.FormatUint(orcaWhirlpool.ProtocolFeeOwedA, 10))
	if err != nil {
		return err
	}
	protocolFeeOwedB, err := decimal.NewFromString(strconv.FormatUint(orcaWhirlpool.ProtocolFeeOwedB, 10))
	if err != nil {
		return err
	}
	rewardLastUpdatedTimestamp, err := decimal.NewFromString(strconv.FormatUint(orcaWhirlpool.RewardLastUpdatedTimestamp, 10))
	if err != nil {
		return err
	}
	liquidity, err := decimal.NewFromString(orcaWhirlpool.Liquidity.String())
	if err != nil {
		return err
	}
	sqrtPrice, err := decimal.NewFromString(orcaWhirlpool.SqrtPrice.String())
	if err != nil {
		return err
	}
	feeGrowthGlobalA, err := decimal.NewFromString(orcaWhirlpool.FeeGrowthGlobalA.String())
	if err != nil {
		return err
	}
	feeGrowthGlobalB, err := decimal.NewFromString(orcaWhirlpool.FeeGrowthGlobalB.String())
	if err != nil {
		return err
	}
	oracle, err := utils.GetWhirlpoolPDA(whirlpoolPubkey.String())
	if err != nil {
		return err
	}

	if err := p.repo.UpsertOrcaWhirlpools(ctx,
		// insert with token pair ID
		&model.OrcaWhirlpool{
			ID:                         uuid.New().String(),
			TokenPairID:                tokenPair.ID,
			Pubkey:                     whirlpoolPubkey.String(),
			WhirlpoolsConfig:           orcaWhirlpool.WhirlpoolsConfig.String(),
			TokenMintA:                 orcaWhirlpool.TokenMintA.String(),
			TokenVaultA:                orcaWhirlpool.TokenVaultA.String(),
			TokenMintB:                 orcaWhirlpool.TokenMintB.String(),
			TokenVaultB:                orcaWhirlpool.TokenVaultB.String(),
			TickSpacing:                int32(orcaWhirlpool.TickSpacing),
			FeeRate:                    int32(orcaWhirlpool.FeeRate),
			ProtocolFeeRate:            int32(orcaWhirlpool.ProtocolFeeRate),
			ProtocolFeeOwedA:           protocolFeeOwedA,
			ProtocolFeeOwedB:           protocolFeeOwedB,
			RewardLastUpdatedTimestamp: rewardLastUpdatedTimestamp,
			TickCurrentIndex:           orcaWhirlpool.TickCurrentIndex,
			Liquidity:                  liquidity,
			SqrtPrice:                  sqrtPrice,
			FeeGrowthGlobalA:           feeGrowthGlobalA,
			FeeGrowthGlobalB:           feeGrowthGlobalB,
			Oracle:                     oracle,
		},
		// insert with inverse token pair ID
		&model.OrcaWhirlpool{
			ID:                         uuid.New().String(),
			TokenPairID:                inverseTokenPair.ID,
			Pubkey:                     whirlpoolPubkey.String(),
			WhirlpoolsConfig:           orcaWhirlpool.WhirlpoolsConfig.String(),
			TokenMintA:                 orcaWhirlpool.TokenMintA.String(),
			TokenVaultA:                orcaWhirlpool.TokenVaultA.String(),
			TokenMintB:                 orcaWhirlpool.TokenMintB.String(),
			TokenVaultB:                orcaWhirlpool.TokenVaultB.String(),
			TickSpacing:                int32(orcaWhirlpool.TickSpacing),
			FeeRate:                    int32(orcaWhirlpool.FeeRate),
			ProtocolFeeRate:            int32(orcaWhirlpool.ProtocolFeeRate),
			ProtocolFeeOwedA:           protocolFeeOwedA,
			ProtocolFeeOwedB:           protocolFeeOwedB,
			RewardLastUpdatedTimestamp: rewardLastUpdatedTimestamp,
			TickCurrentIndex:           orcaWhirlpool.TickCurrentIndex,
			Liquidity:                  liquidity,
			SqrtPrice:                  sqrtPrice,
			FeeGrowthGlobalA:           feeGrowthGlobalA,
			FeeGrowthGlobalB:           feeGrowthGlobalB,
			Oracle:                     oracle,
		},
	); err != nil {
		return err
	}

	if err := p.UpsertTokenAccountBalanceByAddress(ctx, orcaWhirlpool.TokenVaultA.String()); err != nil {
		return err
	}
	if err := p.UpsertTokenAccountBalanceByAddress(ctx, orcaWhirlpool.TokenVaultB.String()); err != nil {
		return err
	}
	// this is hella slow
	// todo: expose this to the backfill dev script so that the script isn't constrained by ctx timeout
	p.maybeUpdateOrcaWhirlPoolDeltaBCache(ctx, whirlpoolPubkey.String(), tokenPair.ID, inverseTokenPair.ID)
	return nil
}

func (p impl) maybeUpdateOrcaWhirlPoolDeltaBCache(ctx context.Context, whirlpoolPubkey string, tokenPairID string, inverseTokenPairID string) {
	log := logrus.WithField("whirlpool", whirlpoolPubkey).WithField("tokenPairID", tokenPairID).WithField("inverseTokenPairID", inverseTokenPairID)
	tokenPairVaults, err := p.repo.AdminGetVaultsByTokenPairID(ctx, tokenPairID)
	if err != nil && err.Error() != repository.ErrRecordNotFound {
		log.WithError(err).Error("failed to get vaults wile updating whirlpool cache")
	}
	vaults, err := p.repo.AdminGetVaultsByTokenPairID(ctx, inverseTokenPairID)
	if err != nil && err.Error() != repository.ErrRecordNotFound {
		log.WithError(err).Error("failed to get vaults wile updating whirlpool cache")
	}
	vaults = append(vaults, tokenPairVaults...)
	quotes := []*model.OrcaWhirlpoolDeltaBQuote{}

	for i := range vaults {
		if vaults[i].DripAmount == 0 {
			continue
		}
		existingQuote, err := p.repo.GetOrcaWhirlpoolDeltaBQuote(ctx, vaults[i].Pubkey, whirlpoolPubkey)
		if err != nil && err.Error() != repository.ErrRecordNotFound {
			log.WithError(err).Error("failed to GetOrcaWhirlpoolDeltaBQuote")
			continue
		}
		existingQuoteLastUpdateTime := utils.GetTimePtr(time.Time{})
		if existingQuote != nil {
			existingQuoteLastUpdateTime = existingQuote.LastUpdated
		}
		// No need to constantly update this, lets only update every 10 minutes
		if time.Now().Before(existingQuoteLastUpdateTime.Add(time.Minute * 10)) {
			continue
		}

		swapQuoteEstimate, err := p.orcaWhirlpoolClient.GetOrcaWhirlpoolQuoteEstimate(ctx, whirlpoolPubkey, vaults[i].TokenAMint, strconv.FormatUint(vaults[i].DripAmount, 10))
		if err != nil {
			log.WithError(err).Error("failed to fetch orcaWhirlpoolQuoteEstimate")
			continue
		}
		deltaB, err := strconv.ParseUint(swapQuoteEstimate.Amount, 10, 64)
		if err != nil {
			log.WithError(err).Error("failed to evaluate whirlpool")
			continue
		}

		quotes = append(quotes, &model.OrcaWhirlpoolDeltaBQuote{
			VaultPubkey:     vaults[i].Pubkey,
			WhirlpoolPubkey: whirlpoolPubkey,
			TokenPairID:     vaults[i].TokenPairID,
			DeltaB:          deltaB,
			LastUpdated:     utils.GetTimePtr(time.Now()),
		})
	}

	if err := p.repo.UpsertOrcaWhirlpoolDeltaBQuotes(ctx, quotes...); err != nil {
		log.WithError(err).Error("failed to UpsertOrcaWhirlpoolDeltaBQuotes")
	}
}
