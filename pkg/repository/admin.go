package repository

import (
	"context"
	"fmt"

	"github.com/dcaf-labs/drip/pkg/repository/model"
	"github.com/sirupsen/logrus"
)

func (d repositoryImpl) AdminGetVaultByAddress(ctx context.Context, pubkey string) (*VaultWithTokenPair, error) {
	stmt := d.repo.
		Vault.WithContext(ctx).
		Join(d.repo.TokenPair, d.repo.TokenPair.ID.EqCol(d.repo.Vault.TokenPairID)).
		Where(d.repo.Vault.Pubkey.Eq(pubkey))
	var vault VaultWithTokenPair
	if err := stmt.Scan(&vault); err != nil {
		return nil, err
	}
	return &vault, nil
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

func (d repositoryImpl) AdminGetVaults(ctx context.Context, vaultFilterParams VaultFilterLikeParams, paginationParams PaginationParams) ([]*model.Vault, error) {
	var vaults []*model.Vault
	query := `SELECT vault.* FROM vault JOIN token_pair ON token_pair.id = vault.token_pair_id`
	var conditions []string
	if vaultFilterParams.IsEnabled != nil {
		conditions = append(conditions, fmt.Sprintf("vault.enabled=%t", *vaultFilterParams.IsEnabled))
	}
	if vaultFilterParams.ProtoConfig != nil {
		conditions = append(conditions, fmt.Sprintf("vault.proto_config=%s", *vaultFilterParams.ProtoConfig))
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
		logrus.WithError(err).WithField("query", query).Error("failed to run query")
		return nil, err
	}
	return vaults, nil
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
