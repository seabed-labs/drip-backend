package repository

import (
	context "context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/dcaf-labs/drip/pkg/repository/model"
)

func (d repositoryImpl) AdminGetVaultByAddress(ctx context.Context, pubkey string) (*model.Vault, error) {
	return d.repo.
		Vault.WithContext(ctx).
		Where(d.repo.Vault.Pubkey.Eq(pubkey)).
		First()
}

func (d repositoryImpl) AdminGetVaultsByAddresses(ctx context.Context, addresses ...string) ([]*model.Vault, error) {
	return d.repo.
		Vault.WithContext(ctx).
		Where(d.repo.Vault.Pubkey.In(addresses...)).
		Find()
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
		logrus.WithError(err).WithField("query", query).Error("failed to run query")
		return nil, err
	}
	return vaults, nil
}
