package repository

import (
	context "context"
	"fmt"

	"github.com/dcaf-labs/drip/pkg/repository/model"
)

func (d repositoryImpl) AdminGetVaultByAddress(ctx context.Context, pubkey string) (*model.Vault, error) {
	return d.repo.
		Vault.WithContext(ctx).
		Where(d.repo.Vault.Pubkey.Eq(pubkey)).
		First()
}

func (d repositoryImpl) AdminGetVaultsByTokenAccountAddress(ctx context.Context, tokenAccountPubkey string) ([]*model.Vault, error) {
	return d.repo.
		Vault.WithContext(ctx).
		Or(d.repo.Vault.TokenAAccount.Eq(tokenAccountPubkey)).
		Or(d.repo.Vault.TokenBAccount.Eq(tokenAccountPubkey)).
		Or(d.repo.Vault.TreasuryTokenBAccount.Eq(tokenAccountPubkey)).
		Find()
}

func (d repositoryImpl) AdminSetVaultEnabled(ctx context.Context, vaultPubkey string, enabled bool) (*model.Vault, error) {
	var res model.Vault
	_, err := d.repo.Vault.
		WithContext(ctx).
		Returning(&res, res.GetAllColumns()...).
		Where(d.repo.Vault.Pubkey.Eq(vaultPubkey)).
		Update(d.repo.Vault.Enabled, enabled)
	return &res, err
}

func (d repositoryImpl) AdminGetVaults(ctx context.Context, vaultFilterParams VaultFilterParams, paginationParams PaginationParams) ([]*model.Vault, error) {
	var vaults []*model.Vault
	query := `SELECT vault.* FROM vault JOIN token_pair ON token_pair.id = vault.token_pair_id`
	var conditions []string
	if vaultFilterParams.IsEnabled != nil {
		conditions = append(conditions, fmt.Sprintf("vault.enabled=%t", *vaultFilterParams.IsEnabled))
	}
	if vaultFilterParams.LikeVault != nil {
		conditions = append(conditions, fmt.Sprintf("vault.pubkey ~ '(?i)%s'", *vaultFilterParams.LikeVault))
	}
	if vaultFilterParams.LikeTokenA != nil {
		conditions = append(conditions, fmt.Sprintf("token_pair.token_a ~ '(?i)%s'", *vaultFilterParams.LikeTokenA))
	}
	if vaultFilterParams.LikeTokenB != nil {
		conditions = append(conditions, fmt.Sprintf("token_pair.token_b ~ '(?i)%s'", *vaultFilterParams.LikeTokenB))
	}
	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE", query)
	}
	for _, condition := range conditions {
		query = fmt.Sprintf("%s %s", query, condition)
	}
	if paginationParams.Limit != nil {
		query = fmt.Sprintf("%s LIMIT %d", query, *paginationParams.Limit)
	}
	if paginationParams.Offset != nil {
		query = fmt.Sprintf("%s OFFSET %d", query, *paginationParams.Offset)
	}
	query = fmt.Sprintf("%s;", query)
	if err := d.db.SelectContext(ctx,
		&vaults,
		query,
	); err != nil {
		return nil, err
	}
	return vaults, nil
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
	stmt = stmt.Where(d.repo.Vault.Enabled.Is(true))
	return stmt.Find()
}

func (d repositoryImpl) GetVaultByAddress(ctx context.Context, address string) (*model.Vault, error) {
	return d.repo.
		Vault.WithContext(ctx).
		Where(d.repo.Vault.Pubkey.Eq(address)).
		Where(d.repo.Vault.Enabled.Is(true)).
		First()
}

func (d repositoryImpl) GetProtoConfigs(ctx context.Context, tokenAMint *string, tokenBMint *string) ([]*model.ProtoConfig, error) {
	stmt := d.repo.ProtoConfig.WithContext(ctx)
	stmt = stmt.Join(d.repo.Vault, d.repo.ProtoConfig.Pubkey.EqCol(d.repo.Vault.ProtoConfig))

	if tokenAMint != nil || tokenBMint != nil {
		stmt = stmt.Join(d.repo.Vault, d.repo.TokenPair.ID.EqCol(d.repo.Vault.TokenPairID))
	}
	if tokenAMint != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenA.Eq(*tokenAMint))
	}
	if tokenBMint != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenB.Eq(*tokenBMint))
	}
	stmt = stmt.Where(d.repo.Vault.Enabled.Is(true))
	return stmt.Find()
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
