package repository

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/repository/model"
)

func (d repositoryImpl) GetVaultByAddress(ctx context.Context, adress string) (*VaultWithTokenPair, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetVaultsWithFilter(ctx context.Context, params VaultFilterParams) ([]*VaultWithTokenPair, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) SearchVaultsWithFilter(ctx context.Context, params VaultSearchFilterParams) ([]*VaultWithTokenPair, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetVaultWhitelistsForVaults(ctx context.Context, vaultAddresses ...string) ([]*model.VaultWhitelist, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetProtoConfigs(ctx context.Context) ([]*model.ProtoConfig, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetProtoConfigsByAddresses(ctx context.Context, addresses ...string) ([]*model.ProtoConfig, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetVaultPeriods(ctx context.Context, vault string, vaultPeriodId *string, paginationParams PaginationParams) ([]*model.VaultPeriod, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetTokenPairByID(ctx context.Context, id string) (*model.TokenPair, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetTokenPairsByIDS(ctx context.Context, ids []string) ([]*model.TokenPair, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetTokenPairs(ctx context.Context, params TokenPairFilterParams) ([]*model.TokenPair, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetTokensByMints(ctx context.Context, mints []string) ([]*model.Token, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetTokenSwapByAddress(ctx context.Context, address string) (*model.TokenSwap, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetTokenSwapsByAddresses(ctx context.Context, addresses []string) ([]*model.TokenSwap, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetTokenSwapsByTokenPairIDsWithBalance(ctx context.Context, tokenPairIDs []string) ([]TokenSwapWithBalance, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetOrcaWhirlpoolByAddress(ctx context.Context, address string) (*model.OrcaWhirlpool, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetOrcaWhirlpoolsByTokenPairIDs(ctx context.Context, tokenPairIDs []string) ([]*model.OrcaWhirlpool, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetPositionByNFTMint(ctx context.Context, nftMint string) (*model.Position, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetPositionsWithFilter(ctx context.Context, isVaultEnabled *bool, positionFilterParams PositionFilterParams, paginationParams PaginationParams) ([]*model.Position, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetTokenAccountBalancesByIDS(ctx context.Context, ids []string) ([]*model.TokenAccountBalance, error) {
	//TODO implement me
	panic("implement me")
}

func (d repositoryImpl) GetActiveWallets(ctx context.Context, params GetActiveWalletParams) ([]ActiveWallet, error) {
	//TODO implement me
	panic("implement me")
}

//
//import (
//	"context"
//
//	"github.com/dcaf-labs/drip/pkg/repository/model"
//	"github.com/google/uuid"
//	"github.com/lib/pq"
//)
//
//func (d repositoryImpl) GetPositionByNFTMint(ctx context.Context, nftMint string) (*model.Position, error) {
//	// The position_authority is the nft mint
//	return d.repo.Position.
//		WithContext(ctx).
//		Where(d.repo.Position.Authority.Eq(nftMint)).
//		First()
//}
//
//func (d repositoryImpl) GetTokenPair(ctx context.Context, tokenA string, tokenB string) (*model.TokenPair, error) {
//	return d.repo.TokenPair.WithContext(ctx).
//		Where(d.repo.TokenPair.TokenA.Eq(tokenA)).
//		Where(d.repo.TokenPair.TokenB.Eq(tokenB)).
//		First()
//}
//
//func (d repositoryImpl) GetTokenPairByID(ctx context.Context, id string) (*model.TokenPair, error) {
//	return d.repo.TokenPair.WithContext(ctx).
//		Where(d.repo.TokenPair.ID.Eq(id)).
//		First()
//}
//
//func (d repositoryImpl) GetTokenPairs(ctx context.Context, params TokenPairFilterParams) ([]*model.TokenPair, error) {
//	stmt := d.repo.TokenPair.WithContext(ctx)
//	if params.TokenA != nil {
//		stmt = stmt.Where(d.repo.TokenPair.TokenA.Eq(*params.TokenA))
//	}
//	if params.TokenB != nil {
//		stmt = stmt.Where(d.repo.TokenPair.TokenB.Eq(*params.TokenB))
//	}
//	return stmt.Find()
//}
//
//func (d repositoryImpl) GetTokenSwaps(ctx context.Context, tokenPairID []string) ([]*model.TokenSwap, error) {
//	stmt := d.repo.TokenSwap.WithContext(ctx)
//	if len(tokenPairID) > 0 {
//		stmt = stmt.Where(d.repo.TokenSwap.TokenPairID.In(tokenPairID...))
//	}
//	return stmt.Find()
//}
//
//func (d repositoryImpl) GetTokenSwapsWithBalance(ctx context.Context, tokenPairIDs []string) ([]TokenSwapWithBalance, error) {
//	var tokenSwaps []TokenSwapWithBalance
//	// TODO(Mocha): No clue how to do this in gorm-gen
//	if len(tokenPairIDs) > 0 {
//		var temp []uuid.UUID
//		for _, tokenPairID := range tokenPairIDs {
//			tokenPairUUID, _ := uuid.Parse(tokenPairID)
//			temp = append(temp, tokenPairUUID)
//		}
//		// We should sort by liquidity ratio descending, so that the largest ratio is at the beginning of the list
//		if err := d.db.SelectContext(ctx,
//			&tokenSwaps,
//			`SELECT token_swap.*, token_account_a_balance.amount as token_account_a_balance_amount, token_account_b_balance.amount as token_account_b_balance_amount
//				FROM token_swap
//				JOIN vault
//				ON vault.token_pair_id = token_swap.token_pair_id
//				JOIN token_account_balance token_account_a_balance
//				ON token_account_a_balance.pubkey = token_swap.token_a_account
//				JOIN token_account_balance token_account_b_balance
//				ON token_account_b_balance.pubkey = token_swap.token_b_account
//				WHERE token_account_a_balance.amount != 0
//				AND token_account_b_balance.amount != 0
//				AND vault.enabled = true
//				AND token_swap.token_pair_id=ANY($1)
//				ORDER BY token_swap.token_pair_id desc;`,
//			pq.Array(temp),
//		); err != nil {
//			return nil, err
//		}
//	} else {
//		if err := d.db.SelectContext(ctx,
//			&tokenSwaps,
//			`SELECT token_swap.*, token_account_a_balance.amount as token_account_a_balance_amount, token_account_b_balance.amount as token_account_b_balance_amount
//				FROM token_swap
//				JOIN vault
//				ON vault.token_pair_id = token_swap.token_pair_id
//				JOIN token_account_balance token_account_a_balance
//				ON token_account_a_balance.pubkey = token_swap.token_a_account
//				JOIN token_account_balance token_account_b_balance
//				ON token_account_b_balance.pubkey = token_swap.token_b_account
//				WHERE token_account_a_balance.amount != 0
//				AND token_account_b_balance.amount != 0
//				AND vault.enabled = true;`,
//		); err != nil {
//			return nil, err
//		}
//	}
//	return tokenSwaps, nil
//}
//
//func (d repositoryImpl) GetTokensWithSupportedTokenPair(ctx context.Context, tokenMint *string, supportedTokenA bool) ([]*model.Token, error) {
//	stmt := d.repo.Token.WithContext(ctx).Distinct(d.repo.Token.ALL)
//	if tokenMint != nil {
//		if supportedTokenA {
//			stmt = stmt.
//				Join(d.repo.TokenPair, d.repo.TokenPair.TokenB.EqCol(d.repo.Token.Pubkey)).
//				Join(d.repo.Vault, d.repo.Vault.TokenPairID.EqCol(d.repo.TokenPair.ID)).
//				Where(d.repo.Vault.Enabled.Is(true)).
//				Where(d.repo.TokenPair.TokenA.Eq(*tokenMint))
//		} else {
//			stmt = stmt.
//				Join(d.repo.TokenPair, d.repo.TokenPair.TokenA.EqCol(d.repo.Token.Pubkey)).
//				Join(d.repo.Vault, d.repo.Vault.TokenPairID.EqCol(d.repo.TokenPair.ID)).
//				Where(d.repo.Vault.Enabled.Is(true)).
//				Where(d.repo.TokenPair.TokenB.Eq(*tokenMint))
//		}
//	}
//	return stmt.Find()
//}
//
//func (d repositoryImpl) GetTokenSwapByAddress(ctx context.Context, address string) (*model.TokenSwap, error) {
//	return d.repo.TokenSwap.WithContext(ctx).Where(d.repo.TokenSwap.Pubkey.Eq(address)).First()
//}
//
//func (d repositoryImpl) GetOrcaWhirlpoolByAddress(ctx context.Context, address string) (*model.OrcaWhirlpool, error) {
//	return d.repo.OrcaWhirlpool.WithContext(ctx).Where(d.repo.OrcaWhirlpool.Pubkey.Eq(address)).First()
//}
//
//func (d repositoryImpl) GetOrcaWhirlpoolsByTokenPairIDs(ctx context.Context, tokenPairIDs []string) ([]*model.OrcaWhirlpool, error) {
//	stmt := d.repo.OrcaWhirlpool.
//		WithContext(ctx).
//		Where(d.repo.OrcaWhirlpool.TokenPairID.In(tokenPairIDs...))
//	return stmt.Find()
//}
//
//func (d repositoryImpl) GetVaultWhitelistsByVaultAddress(ctx context.Context, vaultPubkeys []string) ([]*model.VaultWhitelist, error) {
//	if len(vaultPubkeys) == 0 {
//		return nil, nil
//	}
//	return d.repo.VaultWhitelist.
//		WithContext(ctx).
//		Where(d.repo.VaultWhitelist.VaultPubkey.In(vaultPubkeys...)).
//		Find()
//}
//
//func (d repositoryImpl) GetProtoConfigsByAddresses(ctx context.Context, pubkeys []string) ([]*model.ProtoConfig, error) {
//	stmt := d.repo.ProtoConfig.
//		WithContext(ctx)
//	if len(pubkeys) > 0 {
//		stmt = stmt.Where(d.repo.ProtoConfig.Pubkey.In(pubkeys...))
//	}
//	return stmt.Find()
//}
//
//func (d repositoryImpl) GetTokensByMints(ctx context.Context, mints []string) ([]*model.Token, error) {
//	stmt := d.repo.Token.
//		WithContext(ctx)
//	if len(mints) > 0 {
//		stmt = stmt.Where(d.repo.Token.Pubkey.In(mints...))
//	}
//	return stmt.Find()
//}
//
//func (d repositoryImpl) GetTokenAccountBalancesByIDS(ctx context.Context, tokenAccountPubkeys []string) ([]*model.TokenAccountBalance, error) {
//	stmt := d.repo.TokenAccountBalance.
//		WithContext(ctx)
//	if len(tokenAccountPubkeys) > 0 {
//		stmt = stmt.Where(d.repo.TokenAccountBalance.Pubkey.In(tokenAccountPubkeys...))
//	}
//	return stmt.Find()
//}
//
//func (d repositoryImpl) GetTokenPairsByIDS(ctx context.Context, tokenPairIds []string) ([]*model.TokenPair, error) {
//	stmt := d.repo.TokenPair.
//		WithContext(ctx)
//	if len(tokenPairIds) > 0 {
//		stmt = stmt.Where(d.repo.TokenPair.ID.In(tokenPairIds...))
//	}
//	return stmt.Find()
//}
//
//func (d repositoryImpl) GetVaultPeriods(
//	ctx context.Context,
//	vault string, vaultPeriod *string,
//	paginationParams PaginationParams,
//) ([]*model.VaultPeriod, error) {
//	stmt := d.repo.
//		VaultPeriod.WithContext(ctx).
//		Join(d.repo.Vault, d.repo.VaultPeriod.Vault.EqCol(d.repo.Vault.Pubkey)).
//		Where(d.repo.VaultPeriod.Vault.Eq(vault)).
//		Where(d.repo.Vault.Enabled.Is(true))
//	if vaultPeriod != nil {
//		stmt = stmt.Where(d.repo.VaultPeriod.Pubkey.Eq(*vaultPeriod))
//	}
//	if paginationParams.Limit != nil {
//		stmt = stmt.Limit(*paginationParams.Limit)
//	}
//	if paginationParams.Offset != nil {
//		stmt = stmt.Offset(*paginationParams.Offset)
//	}
//	return stmt.Find()
//}
//
//func (d repositoryImpl) GetVaultsWithFilter(ctx context.Context, filter VaultFilterParams) ([]*VaultWithTokenPair, error) {
//	stmt := d.repo.Vault.WithContext(ctx).
//		Join(d.repo.Vault, d.repo.TokenPair.ID.EqCol(d.repo.Vault.TokenPairID))
//	if filter.TokenA != nil {
//		stmt = stmt.Where(d.repo.TokenPair.TokenA.Eq(*filter.TokenA))
//	}
//	if filter.TokenB != nil {
//		stmt = stmt.Where(d.repo.TokenPair.TokenB.Eq(*filter.TokenB))
//	}
//	if filter.ProtoConfig != nil {
//		stmt = stmt.Where(d.repo.Vault.ProtoConfig.Eq(*filter.ProtoConfig))
//	}
//	if filter.Vault != nil {
//		stmt = stmt.Where(d.repo.Vault.Pubkey.Eq(*filter.Vault))
//	}
//	var res []*VaultWithTokenPair
//	if err := stmt.Scan(&res); err != nil {
//		return nil, err
//	}
//	return res, nil
//}
//
//func (d repositoryImpl) GetActiveWallets(
//	ctx context.Context, params GetActiveWalletParams,
//) ([]ActiveWallet, error) {
//	var res []ActiveWallet
//	stmt := d.repo.TokenAccountBalance.WithContext(ctx).
//		Select(
//			d.repo.TokenAccountBalance.Owner.As("owner"),
//			d.repo.TokenAccountBalance.Owner.Count().As("position_count"),
//		).
//		Join(d.repo.Position, d.repo.Position.Authority.EqCol(d.repo.TokenAccountBalance.Mint)).
//		Join(d.repo.Vault, d.repo.Vault.Pubkey.EqCol(d.repo.Position.Vault)).
//		Where(d.repo.Vault.Enabled.Is(true))
//
//	if params.Owner != nil {
//		stmt = stmt.Where(d.repo.TokenAccountBalance.Owner.Eq(*params.Owner))
//	}
//	if params.PositionIsClosed != nil {
//		stmt = stmt.Where(d.repo.Position.IsClosed.Is(*params.PositionIsClosed))
//	}
//	if params.Vault != nil {
//		stmt = stmt.Where(d.repo.Vault.Pubkey.Eq(*params.Vault))
//	}
//	err := stmt.
//		Group(d.repo.TokenAccountBalance.Owner).
//		Scan(&res)
//	return res, err
//}
//
//func (d repositoryImpl) GetTokensWithSupportedTokenB(ctx context.Context, tokenBMint *string) ([]*model.Token, error) {
//	stmt := d.repo.Token.WithContext(ctx).Distinct(d.repo.Token.ALL)
//	if tokenBMint != nil {
//		stmt = stmt.
//			Join(d.repo.TokenPair, d.repo.TokenPair.TokenB.EqCol(d.repo.Token.Pubkey)).
//			Join(d.repo.Vault, d.repo.Vault.TokenPairID.EqCol(d.repo.TokenPair.ID)).
//			Where(d.repo.Vault.Enabled.Is(true))
//	}
//	return stmt.Find()
//}
//
//func (d repositoryImpl) GetVaultByAddress(ctx context.Context, address string) (*VaultWithTokenPair, error) {
//	stmt := d.repo.
//		Vault.WithContext(ctx).
//		Join(d.repo.TokenPair, d.repo.TokenPair.ID.EqCol(d.repo.Vault.TokenPairID)).
//		Where(d.repo.Vault.Pubkey.Eq(address)).
//		Where(d.repo.Vault.Enabled.Is(true))
//	var res VaultWithTokenPair
//	if err := stmt.Scan(&res); err != nil {
//		return nil, err
//	}
//	return &res, nil
//}
//
//func (d repositoryImpl) GetProtoConfigs(ctx context.Context) ([]*model.ProtoConfig, error) {
//	stmt := d.repo.ProtoConfig.WithContext(ctx)
//	stmt = stmt.Join(d.repo.Vault, d.repo.ProtoConfig.Pubkey.EqCol(d.repo.Vault.ProtoConfig))
//	stmt = stmt.Where(d.repo.Vault.Enabled.Is(true))
//	return stmt.Find()
//}
