package repository

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (d repositoryImpl) GetProtoConfigs(ctx context.Context, filterParams ProtoConfigParams) ([]*model.ProtoConfig, error) {
	stmt := d.repo.ProtoConfig.WithContext(ctx)
	if filterParams.TokenA != nil || filterParams.TokenB != nil {
		stmt = stmt.
			Join(d.repo.Vault, d.repo.Vault.ProtoConfig.EqCol(d.repo.ProtoConfig.Pubkey)).
			Join(d.repo.TokenPair, d.repo.TokenPair.ID.EqCol(d.repo.Vault.TokenPairID)).
			Where(d.repo.Vault.Enabled.Is(true))
	}
	if filterParams.TokenA != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenA.Eq(*filterParams.TokenA))
	}
	if filterParams.TokenB != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenB.Eq(*filterParams.TokenB))
	}
	// default ascending
	stmt = stmt.Order(d.repo.ProtoConfig.Granularity)
	return stmt.Find()
}

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

func (d repositoryImpl) GetSPLTokenSwapsByTokenPairIDs(ctx context.Context, tokenPairIDs ...string) ([]*model.TokenSwap, error) {
	stmt := d.repo.TokenSwap.WithContext(ctx)
	if len(tokenPairIDs) > 0 {
		stmt = stmt.Where(d.repo.TokenSwap.TokenPairID.In(tokenPairIDs...))
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

// todo: dedupe with GetAllSupportedTokenAs
func (d repositoryImpl) GetAllSupportTokens(ctx context.Context) ([]*model.Token, error) {
	tokenPairs, err := d.repo.TokenPair.WithContext(ctx).
		Join(d.repo.Vault, d.repo.Vault.TokenPairID.EqCol(d.repo.TokenPair.ID)).
		Where(d.repo.Vault.Enabled.Is(true)).
		Find()
	if err != nil {
		return nil, err
	}
	tokenMintSet := make(map[string]bool)
	for _, pair := range tokenPairs {
		tokenMintSet[pair.TokenA] = true
		tokenMintSet[pair.TokenB] = true
	}
	tokenMints := []string{}
	for mint := range tokenMintSet {
		tokenMints = append(tokenMints, mint)
	}
	if len(tokenMints) == 0 {
		return []*model.Token{}, nil
	}
	return d.GetTokensByMints(ctx, tokenMints)
}

func (d repositoryImpl) GetAllSupportedTokenAs(ctx context.Context) ([]*model.Token, error) {
	tokenPairs, err := d.repo.TokenPair.WithContext(ctx).
		Join(d.repo.Vault, d.repo.Vault.TokenPairID.EqCol(d.repo.TokenPair.ID)).
		Where(d.repo.Vault.Enabled.Is(true)).
		Find()
	if err != nil {
		return nil, err
	}
	tokenMintSet := make(map[string]bool)
	for _, pair := range tokenPairs {
		tokenMintSet[pair.TokenA] = true
	}
	tokenMints := []string{}
	for mint := range tokenMintSet {
		tokenMints = append(tokenMints, mint)
	}
	if len(tokenMints) == 0 {
		return []*model.Token{}, nil
	}
	return d.GetTokensByMints(ctx, tokenMints)
}

func (d repositoryImpl) GetSupportedTokenAs(ctx context.Context, tokenBMint *string) ([]*model.Token, error) {
	if tokenBMint == nil {
		return d.GetAllSupportedTokenAs(ctx)
	}
	return d.repo.Token.WithContext(ctx).
		Distinct(d.repo.Token.ALL).
		Join(d.repo.TokenPair, d.repo.TokenPair.TokenA.EqCol(d.repo.Token.Pubkey)).
		Join(d.repo.Vault, d.repo.Vault.TokenPairID.EqCol(d.repo.TokenPair.ID)).
		Where(d.repo.Vault.Enabled.Is(true)).
		Where(d.repo.TokenPair.TokenB.Eq(*tokenBMint)).
		Where(d.repo.Token.Pubkey.Neq(*tokenBMint)).
		Order(d.repo.Token.Symbol).
		Find()
}

func (d repositoryImpl) GetSupportedTokenBs(ctx context.Context, tokenAMint string) ([]*model.Token, error) {
	return d.repo.Token.WithContext(ctx).
		Distinct(d.repo.Token.ALL).
		Join(d.repo.TokenPair, d.repo.TokenPair.TokenB.EqCol(d.repo.Token.Pubkey)).
		Join(d.repo.Vault, d.repo.Vault.TokenPairID.EqCol(d.repo.TokenPair.ID)).
		Where(d.repo.Vault.Enabled.Is(true)).
		Where(d.repo.TokenPair.TokenA.Eq(tokenAMint)).
		Where(d.repo.Token.Pubkey.Neq(tokenAMint)).
		Order(d.repo.Token.Symbol).
		Find()
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

func (d repositoryImpl) GetTokenAccountBalancesByAddresses(ctx context.Context, tokenAccountPubkeys ...string) ([]*model.TokenAccountBalance, error) {
	stmt := d.repo.TokenAccountBalance.
		WithContext(ctx)
	if len(tokenAccountPubkeys) > 0 {
		stmt = stmt.Where(d.repo.TokenAccountBalance.Pubkey.In(tokenAccountPubkeys...))
	}
	return stmt.Find()
}

func (d repositoryImpl) GetVaultPeriodByAddress(
	ctx context.Context,
	address string,
) (*model.VaultPeriod, error) {
	return d.repo.
		VaultPeriod.WithContext(ctx).
		Where(d.repo.VaultPeriod.Pubkey.Eq(address)).
		First()
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

func (d repositoryImpl) GetActiveWallets(
	ctx context.Context, params GetActiveWalletParams,
) ([]ActiveWallet, error) {
	var res []ActiveWallet
	stmt := d.repo.TokenAccountBalance.WithContext(ctx).
		Select(
			d.repo.TokenAccountBalance.Owner.As("owner"),
			d.repo.TokenAccountBalance.Owner.Count().As("position_count"),
		).
		Join(d.repo.Position, d.repo.Position.Authority.EqCol(d.repo.TokenAccountBalance.Mint)).
		Join(d.repo.Vault, d.repo.Vault.Pubkey.EqCol(d.repo.Position.Vault)).
		Where(d.repo.Vault.Enabled.Is(true))

	if params.Owner != nil {
		stmt = stmt.Where(d.repo.TokenAccountBalance.Owner.Eq(*params.Owner))
	}
	if params.PositionIsClosed != nil {
		stmt = stmt.Where(d.repo.Position.IsClosed.Is(*params.PositionIsClosed))
	}
	if params.Vault != nil {
		stmt = stmt.Where(d.repo.Vault.Pubkey.Eq(*params.Vault))
	}
	err := stmt.
		Group(d.repo.TokenAccountBalance.Owner).
		Scan(&res)
	return res, err
}

func (d repositoryImpl) GetTokensWithSupportedTokenB(ctx context.Context, tokenBMint *string) ([]*model.Token, error) {
	stmt := d.repo.Token.WithContext(ctx).Distinct(d.repo.Token.ALL)
	if tokenBMint != nil {
		stmt = stmt.
			Join(d.repo.TokenPair, d.repo.TokenPair.TokenB.EqCol(d.repo.Token.Pubkey)).
			Join(d.repo.Vault, d.repo.Vault.TokenPairID.EqCol(d.repo.TokenPair.ID)).
			Where(d.repo.Vault.Enabled.Is(true))
	}
	return stmt.Find()
}

func (d repositoryImpl) GetVaultsWithFilter(ctx context.Context, tokenAMint, tokenBMint, protoConfig *string) ([]*model.Vault, error) {
	stmt := d.repo.Vault.WithContext(ctx)
	if tokenAMint != nil || tokenBMint != nil {
		stmt = stmt.Join(d.repo.Vault, d.repo.TokenPair.ID.EqCol(d.repo.Vault.TokenPairID))
	}
	if tokenAMint != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenA.Eq(*tokenAMint))
	}
	if tokenBMint != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenB.Eq(*tokenBMint))
	}
	if protoConfig != nil {
		stmt = stmt.Where(d.repo.Vault.ProtoConfig.Eq(*protoConfig))
	}
	stmt = stmt.Where(d.repo.Vault.Enabled.Is(true)).Order(d.repo.Vault.Pubkey)
	return stmt.Find()
}

func (d repositoryImpl) GetVaultByAddress(ctx context.Context, address string) (*model.Vault, error) {
	return d.repo.
		Vault.WithContext(ctx).
		Where(d.repo.Vault.Pubkey.Eq(address)).
		Where(d.repo.Vault.Enabled.Is(true)).
		First()
}

func (d repositoryImpl) GetAdminPositions(
	ctx context.Context, isVaultEnabled *bool,
	positionFilterParams PositionFilterParams,
	params PaginationParams,
) ([]*model.Position, error) {
	stmt := d.repo.Position.WithContext(ctx)

	// Apply Joins
	if isVaultEnabled != nil {
		stmt = stmt.Join(d.repo.Vault, d.repo.Vault.Pubkey.EqCol(d.repo.Position.Vault))
	}
	if positionFilterParams.Wallet != nil {
		stmt = stmt.
			Join(d.repo.TokenAccountBalance, d.repo.TokenAccountBalance.Mint.EqCol(d.repo.Position.Authority))
	}

	// Apply Filters
	if isVaultEnabled != nil {
		stmt = stmt.Where(d.repo.Vault.Enabled.Is(*isVaultEnabled))
	}
	if positionFilterParams.Wallet != nil {
		stmt = stmt.
			Where(
				d.repo.TokenAccountBalance.Owner.Eq(*positionFilterParams.Wallet),
				d.repo.TokenAccountBalance.Amount.Gt(0))
	}
	if positionFilterParams.IsClosed != nil {
		stmt = stmt.Where(d.repo.Position.IsClosed.Is(*positionFilterParams.IsClosed))
	}
	if params.Limit != nil {
		stmt = stmt.Limit(*params.Limit)
	}
	if params.Offset != nil {
		stmt = stmt.Offset(*params.Offset)
	}
	return stmt.Find()
}
