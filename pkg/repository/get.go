package repository

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/repository/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (d repositoryImpl) GetPositionByNFTMint(ctx context.Context, nftMint string) (*model.Position, error) {
	// The position_authority is the nft mint
	return d.repo.Position.
		WithContext(ctx).
		Where(d.repo.Position.Authority.Eq(nftMint)).
		First()
}

func (d repositoryImpl) GetTokenPair(ctx context.Context, tokenA string, tokenB string) (*model.TokenPair, error) {
	return d.repo.TokenPair.WithContext(ctx).
		Where(d.repo.TokenPair.TokenA.Eq(tokenA)).
		Where(d.repo.TokenPair.TokenB.Eq(tokenB)).
		First()
}

func (d repositoryImpl) GetTokenPairByID(ctx context.Context, id string) (*model.TokenPair, error) {
	return d.repo.TokenPair.WithContext(ctx).
		Where(d.repo.TokenPair.ID.Eq(id)).
		First()
}

func (d repositoryImpl) GetTokenPairs(ctx context.Context, tokenAMint *string, tokenBMint *string) ([]*model.TokenPair, error) {
	stmt := d.repo.TokenPair.WithContext(ctx)
	if tokenAMint != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenA.Eq(*tokenAMint))
	}
	if tokenBMint != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenB.Eq(*tokenBMint))
	}
	return stmt.Find()
}

func (d repositoryImpl) GetTokenSwaps(ctx context.Context, tokenPairID []string) ([]*model.TokenSwap, error) {
	stmt := d.repo.TokenSwap.WithContext(ctx)
	if len(tokenPairID) > 0 {
		stmt = stmt.Where(d.repo.TokenSwap.TokenPairID.In(tokenPairID...))
	}
	return stmt.Find()
}

func (d repositoryImpl) GetTokenSwapsWithBalance(ctx context.Context, tokenPairIDs []string) ([]TokenSwapWithBalance, error) {
	var tokenSwaps []TokenSwapWithBalance
	// TODO(Mocha): No clue how to do this in gorm-gen
	if len(tokenPairIDs) > 0 {
		var temp []uuid.UUID
		for _, tokenPairID := range tokenPairIDs {
			tokenPairUUID, _ := uuid.Parse(tokenPairID)
			temp = append(temp, tokenPairUUID)
		}
		// We should sort by liquidity ratio descending, so that the largest ratio is at the beginning of the list
		if err := d.db.SelectContext(ctx,
			&tokenSwaps,
			`SELECT token_swap.*, token_account_a_balance.amount as token_account_a_balance_amount, token_account_b_balance.amount as token_account_b_balance_amount
				FROM token_swap
				JOIN vault
				ON vault.token_pair_id = token_swap.token_pair_id
				JOIN token_account_balance token_account_a_balance
				ON token_account_a_balance.pubkey = token_swap.token_a_account
				JOIN token_account_balance token_account_b_balance
				ON token_account_b_balance.pubkey = token_swap.token_b_account
				WHERE token_account_a_balance.amount != 0
				AND token_account_b_balance.amount != 0
				AND vault.enabled = true
				AND token_swap.token_pair_id=ANY($1)
				ORDER BY token_swap.token_pair_id desc;`,
			pq.Array(temp),
		); err != nil {
			return nil, err
		}
	} else {
		if err := d.db.SelectContext(ctx,
			&tokenSwaps,
			`SELECT token_swap.*, token_account_a_balance.amount as token_account_a_balance_amount, token_account_b_balance.amount as token_account_b_balance_amount
				FROM token_swap
				JOIN vault
				ON vault.token_pair_id = token_swap.token_pair_id
				JOIN token_account_balance token_account_a_balance
				ON token_account_a_balance.pubkey = token_swap.token_a_account
				JOIN token_account_balance token_account_b_balance
				ON token_account_b_balance.pubkey = token_swap.token_b_account
				WHERE token_account_a_balance.amount != 0
				AND token_account_b_balance.amount != 0
				AND vault.enabled = true;`,
		); err != nil {
			return nil, err
		}
	}
	return tokenSwaps, nil
}

func (d repositoryImpl) GetTokensWithSupportedTokenPair(ctx context.Context, tokenMint *string, supportedTokenA bool) ([]*model.Token, error) {
	stmt := d.repo.Token.WithContext(ctx).Distinct(d.repo.Token.ALL)
	if tokenMint != nil {
		if supportedTokenA {
			stmt = stmt.
				Join(d.repo.TokenPair, d.repo.TokenPair.TokenB.EqCol(d.repo.Token.Pubkey)).
				Join(d.repo.Vault, d.repo.Vault.TokenPairID.EqCol(d.repo.TokenPair.ID)).
				Where(d.repo.Vault.Enabled.Is(true)).
				Where(d.repo.TokenPair.TokenA.Eq(*tokenMint))
		} else {
			stmt = stmt.
				Join(d.repo.TokenPair, d.repo.TokenPair.TokenA.EqCol(d.repo.Token.Pubkey)).
				Join(d.repo.Vault, d.repo.Vault.TokenPairID.EqCol(d.repo.TokenPair.ID)).
				Where(d.repo.Vault.Enabled.Is(true)).
				Where(d.repo.TokenPair.TokenB.Eq(*tokenMint))
		}
	}
	return stmt.Find()
}

func (d repositoryImpl) GetTokenSwapByAddress(ctx context.Context, address string) (*model.TokenSwap, error) {
	return d.repo.TokenSwap.WithContext(ctx).Where(d.repo.TokenSwap.Pubkey.Eq(address)).First()
}

func (d repositoryImpl) GetOrcaWhirlpoolByAddress(ctx context.Context, address string) (*model.OrcaWhirlpool, error) {
	return d.repo.OrcaWhirlpool.WithContext(ctx).Where(d.repo.OrcaWhirlpool.Pubkey.Eq(address)).First()
}

func (d repositoryImpl) GetOrcaWhirlpoolsByTokenPairIDs(ctx context.Context, tokenPairIDs []string) ([]*model.OrcaWhirlpool, error) {
	stmt := d.repo.OrcaWhirlpool.
		WithContext(ctx).
		Where(d.repo.OrcaWhirlpool.TokenPairID.In(tokenPairIDs...))
	return stmt.Find()
}

func (d repositoryImpl) GetVaultWhitelistsByVaultAddress(ctx context.Context, vaultPubkeys []string) ([]*model.VaultWhitelist, error) {
	if len(vaultPubkeys) == 0 {
		return nil, nil
	}
	return d.repo.VaultWhitelist.
		WithContext(ctx).
		Where(d.repo.VaultWhitelist.VaultPubkey.In(vaultPubkeys...)).
		Find()
}

func (d repositoryImpl) GetProtoConfigsByAddresses(ctx context.Context, pubkeys []string) ([]*model.ProtoConfig, error) {
	stmt := d.repo.ProtoConfig.
		WithContext(ctx)
	if len(pubkeys) > 0 {
		stmt = stmt.Where(d.repo.ProtoConfig.Pubkey.In(pubkeys...))
	}
	return stmt.Find()
}

func (d repositoryImpl) GetTokensByMints(ctx context.Context, mints []string) ([]*model.Token, error) {
	stmt := d.repo.Token.
		WithContext(ctx)
	if len(mints) > 0 {
		stmt = stmt.Where(d.repo.Token.Pubkey.In(mints...))
	}
	return stmt.Find()
}

func (d repositoryImpl) GetTokenAccountBalancesByIDS(ctx context.Context, tokenAccountPubkeys []string) ([]*model.TokenAccountBalance, error) {
	stmt := d.repo.TokenAccountBalance.
		WithContext(ctx)
	if len(tokenAccountPubkeys) > 0 {
		stmt = stmt.Where(d.repo.TokenAccountBalance.Pubkey.In(tokenAccountPubkeys...))
	}
	return stmt.Find()
}

func (d repositoryImpl) GetTokenPairsByIDS(ctx context.Context, tokenPairIds []string) ([]*model.TokenPair, error) {
	stmt := d.repo.TokenPair.
		WithContext(ctx)
	if len(tokenPairIds) > 0 {
		stmt = stmt.Where(d.repo.TokenPair.ID.In(tokenPairIds...))
	}
	return stmt.Find()
}

func (d repositoryImpl) GetVaultPeriods(
	ctx context.Context,
	vault string, vaultPeriod *string,
	paginationParams PaginationParams,
) ([]*model.VaultPeriod, error) {
	stmt := d.repo.
		VaultPeriod.WithContext(ctx).
		Join(d.repo.Vault, d.repo.VaultPeriod.Vault.EqCol(d.repo.Vault.Pubkey)).
		Where(d.repo.VaultPeriod.Vault.Eq(vault)).
		Where(d.repo.Vault.Enabled.Is(true))
	if vaultPeriod != nil {
		stmt = stmt.Where(d.repo.VaultPeriod.Pubkey.Eq(*vaultPeriod))
	}
	if paginationParams.Limit != nil {
		stmt = stmt.Limit(*paginationParams.Limit)
	}
	if paginationParams.Offset != nil {
		stmt = stmt.Offset(*paginationParams.Offset)
	}
	return stmt.Find()
}
