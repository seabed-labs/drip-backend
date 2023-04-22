package repository

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
)

func (d repositoryImpl) AdminGetVaultByAddress(ctx context.Context, pubkey string) (*model.Vault, error) {
	return d.repo.
		Vault.WithContext(ctx).
		Where(d.repo.Vault.Pubkey.Eq(pubkey)).
		First()
}

func (d repositoryImpl) AdminGetVaultByTreasuryTokenBAccount(ctx context.Context, pubkey string) (*model.Vault, error) {
	return d.repo.
		Vault.WithContext(ctx).
		Where(d.repo.Vault.TreasuryTokenBAccount.Eq(pubkey)).
		First()
}

func (d repositoryImpl) AdminGetVaultsByAddresses(ctx context.Context, addresses ...string) ([]*model.Vault, error) {
	return d.repo.
		Vault.WithContext(ctx).
		Where(d.repo.Vault.Pubkey.In(addresses...)).
		Find()
}

func (d repositoryImpl) AdminGetVaultsByTokenPairID(ctx context.Context, tokenPairID string) ([]*model.Vault, error) {
	return d.repo.
		Vault.WithContext(ctx).
		Where(d.repo.Vault.TokenPairID.Eq(tokenPairID)).
		Find()
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
			Join(d.repo.TokenAccount, d.repo.TokenAccount.Mint.EqCol(d.repo.Position.Authority))
	}

	// Apply Filters
	if isVaultEnabled != nil {
		stmt = stmt.Where(d.repo.Vault.Enabled.Is(*isVaultEnabled))
	}
	if positionFilterParams.Wallet != nil {
		stmt = stmt.
			Where(
				d.repo.TokenAccount.Owner.Eq(*positionFilterParams.Wallet),
				d.repo.TokenAccount.Amount.Gt(0))
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
	stmt := d.repo.
		Vault.WithContext(ctx).
		Join(d.repo.TokenPair, d.repo.TokenPair.ID.EqCol(d.repo.Vault.TokenPairID))

	if vaultFilterParams.Vault != nil {
		stmt = stmt.Where(d.repo.Vault.Pubkey.Eq(*vaultFilterParams.Vault))
	}
	if vaultFilterParams.IsEnabled != nil {
		stmt = stmt.Where(d.repo.Vault.Enabled.Is(*vaultFilterParams.IsEnabled))
	}
	if vaultFilterParams.VaultProtoConfig != nil {
		stmt = stmt.Where(d.repo.Vault.ProtoConfig.Eq(*vaultFilterParams.VaultProtoConfig))
	}
	if vaultFilterParams.TokenA != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenA.Eq(*vaultFilterParams.TokenA))
	}
	if vaultFilterParams.TokenB != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenB.Eq(*vaultFilterParams.TokenB))
	}

	if paginationParams.Limit != nil {
		stmt = stmt.Limit(*paginationParams.Limit)
	}
	if paginationParams.Offset != nil {
		stmt = stmt.Offset(*paginationParams.Offset)
	}

	return stmt.Order(d.repo.Vault.ProtoConfig, d.repo.Vault.Enabled).Find()
}
