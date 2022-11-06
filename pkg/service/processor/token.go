package processor

import (
	"context"
	"fmt"

	"github.com/dcaf-labs/drip/pkg/service/clients/solana"

	"github.com/dcaf-labs/drip/pkg/service/utils"

	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/sirupsen/logrus"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/google/uuid"
)

func (p impl) UpsertTokenPair(ctx context.Context, tokenAAMint string, tokenBMint string) error {
	var tokenA token.Mint
	if err := p.solanaClient.GetAccount(ctx, tokenAAMint, &tokenA); err != nil {
		return err
	}
	var tokenB token.Mint
	if err := p.solanaClient.GetAccount(ctx, tokenBMint, &tokenB); err != nil {
		return err
	}
	if err := p.UpsertTokenByAddress(ctx, tokenAAMint); err != nil {
		return err
	}
	if err := p.UpsertTokenByAddress(ctx, tokenBMint); err != nil {
		return err
	}
	return p.repo.InsertTokenPairs(ctx, &model.TokenPair{
		ID:     uuid.New().String(),
		TokenA: tokenAAMint,
		TokenB: tokenBMint,
	})
}

func (p impl) UpsertTokenByAddress(ctx context.Context, mintAddress string) error {
	tokenMint, err := p.solanaClient.GetTokenMint(ctx, mintAddress)
	if err != nil {
		return fmt.Errorf("failed to GetTokenMint %s, err: %w", mintAddress, err)
	}
	existingToken, err := p.repo.GetTokenByAddress(ctx, mintAddress)
	if err != nil && err.Error() != repository.ErrRecordNotFound {
		return err
	}
	symbol, name, iconURL, coinGeckoID, marketCapRank, marketPrice := p.getTokenMetadata(ctx, mintAddress, existingToken)
	tokenModel := model.Token{
		Pubkey:        mintAddress,
		Symbol:        symbol,
		Name:          name,
		Decimals:      int16(tokenMint.Decimals),
		IconURL:       iconURL,
		CoinGeckoID:   coinGeckoID,
		MarketCapRank: marketCapRank,
		UIMarketPrice: marketPrice,
	}
	return p.repo.UpsertTokens(ctx, &tokenModel)
}

func (p impl) getTokenMetadata(
	ctx context.Context, mint string, existingToken *model.Token,
) (symbol *string, name *string, iconURL *string, coinGeckoID *string, marketCapRank *int32, marketPrice *float64) {
	if existingToken != nil {
		symbol = existingToken.Symbol
		name = existingToken.Name
		iconURL = existingToken.IconURL
		coinGeckoID = existingToken.CoinGeckoID
	}
	// try and populate symbol with token metadata if it exists
	if symbol == nil {
		if tokenMetadataAccount, err := p.solanaClient.GetTokenMetadataAccount(ctx, mint); err != nil {
			if err.Error() != solana.ErrNotFound {
				logrus.WithError(err).Error("failed to GetTokenMetadataAccount")
			}
		} else {
			symbol = &tokenMetadataAccount.Data.Symbol
		}
	}
	// try and backfill the rest of the data with coingecko if it doesn't exist
	shouldFetchCoinGeckoMetadata := symbol == nil || name == nil || iconURL == nil || coinGeckoID == nil || marketPrice == nil
	if !shouldFetchCoinGeckoMetadata {
		return
	}
	// fill coingecko meta data
	if coinGeckoMeta, err := p.coinGeckoClient.GetCoinGeckoMetadata(ctx, mint); err != nil {
		logrus.WithError(err).Error("failed to GetCoinGeckoMetadata")
	} else {
		coinGeckoID = utils.GetStringPtr(coinGeckoMeta.ID)
		marketCapRank = coinGeckoMeta.MarketCapRank
		if usdPrice, ok := coinGeckoMeta.MarketData.CurrentPrice["usd"]; ok {
			marketPrice = utils.GetFloat64Ptr(usdPrice)
		}
		if symbol == nil {
			symbol = utils.GetStringPtr(coinGeckoMeta.Symbol)
		}
		if name == nil {
			name = utils.GetStringPtr(coinGeckoMeta.Name)
		}
		if iconURL == nil {
			iconURL = coinGeckoMeta.Image.Small
		}
	}
	return
}

func (p impl) shouldIngestTokenBalance(ctx context.Context, tokenAccountAddress string, tokenAccount token.Account) bool {
	if p.IsTokenSwapTokenAccount(ctx, tokenAccount.Owner.String()) ||
		p.isOrcaWhirlpoolTokenAccount(ctx, tokenAccount.Owner.String()) ||
		p.isVaultTokenAccount(ctx, tokenAccount.Owner.String()) ||
		p.isVaultTreasuryAccount(ctx, tokenAccountAddress) ||
		p.isUserPositionTokenAccount(ctx, tokenAccount.Mint.String()) {
		return true
	}
	return false
}

func (p impl) UpsertTokenAccountByAddress(ctx context.Context, address string) error {
	var tokenAccount token.Account
	if err := p.solanaClient.GetAccount(ctx, address, &tokenAccount); err != nil {
		return err
	}
	return p.UpsertTokenAccount(ctx, address, tokenAccount)
}

func (p impl) UpsertTokenAccount(ctx context.Context, address string, tokenAccount token.Account) error {
	if !p.shouldIngestTokenBalance(ctx, address, tokenAccount) {
		return nil
	}
	state := "initialized"
	if tokenAccount.State == token.Uninitialized {
		state = "uninitialized"
	} else if tokenAccount.State == token.Frozen {
		state = "frozen"
	}

	var tokenMint token.Mint
	if err := p.solanaClient.GetAccount(ctx, tokenAccount.Mint.String(), &tokenMint); err != nil {
		return fmt.Errorf("failed to GetAccount %s, err: %w", tokenAccount.Mint.String(), err)
	}
	if err := p.UpsertTokenByAddress(ctx, tokenAccount.Mint.String()); err != nil {
		return fmt.Errorf("failed to UpsertTokenByAddress %s, err: %w", tokenAccount.Mint.String(), err)
	}
	return p.repo.UpsertTokenAccounts(ctx, &model.TokenAccount{
		Pubkey: address,
		Mint:   tokenAccount.Mint.String(),
		Owner:  tokenAccount.Owner.String(),
		Amount: tokenAccount.Amount,
		State:  state,
	})
}

func (p impl) IsTokenSwapTokenAccount(ctx context.Context, tokenAccountOwner string) bool {
	_, err := p.repo.GetTokenSwapByAddress(ctx, tokenAccountOwner)
	if err != nil {
		return false
	}
	if err != nil {
		if err.Error() != repository.ErrRecordNotFound {
			logrus.WithError(err).Error("failed to query for token swap")
		}
		return false
	}
	return true
}

func (p impl) isOrcaWhirlpoolTokenAccount(ctx context.Context, tokenAccountOwner string) bool {
	_, err := p.repo.GetOrcaWhirlpoolByAddress(ctx, tokenAccountOwner)
	if err != nil {
		if err.Error() != repository.ErrRecordNotFound {
			logrus.WithError(err).Error("failed to query for whirlpool")
		}
		return false
	}
	return true
}

func (p impl) isVaultTokenAccount(ctx context.Context, tokenAccountOwner string) bool {
	_, err := p.repo.AdminGetVaultByAddress(ctx, tokenAccountOwner)
	if err != nil {
		if err.Error() != repository.ErrRecordNotFound {
			logrus.WithError(err).Error("failed to query for vault")
		}
		return false
	}
	return true
}

func (p impl) isVaultTreasuryAccount(ctx context.Context, tokenAccount string) bool {
	_, err := p.repo.AdminGetVaultByTreasuryTokenBAccount(ctx, tokenAccount)
	if err != nil {
		if err.Error() != repository.ErrRecordNotFound {
			logrus.WithError(err).Error("failed to query for vault")
		}
		return false
	}
	return true
}

func (p impl) isUserPositionTokenAccount(ctx context.Context, mint string) bool {
	_, err := p.repo.GetPositionByNFTMint(ctx, mint)
	if err != nil {
		if err.Error() != repository.ErrRecordNotFound {
			logrus.WithError(err).Error("failed to query for position")
		}
		return false
	}
	return true
}

// ensureTokenPair - if token pair exists return it, else upsert tokenPair and all needed tokenPair foreign keys
func (p impl) ensureTokenPair(ctx context.Context, tokenAAMint string, tokenBMint string) (*model.TokenPair, error) {
	tokenPair, err := p.repo.GetTokenPair(ctx, tokenAAMint, tokenBMint)
	if err != nil && err.Error() == repository.ErrRecordNotFound {
		if err := p.UpsertTokenPair(ctx, tokenAAMint, tokenBMint); err != nil {
			return nil, err
		}
		return p.repo.GetTokenPair(ctx, tokenAAMint, tokenBMint)
	}

	return tokenPair, err
}
