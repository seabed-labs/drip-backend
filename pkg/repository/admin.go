package repository

import (
	context "context"

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

func (d repositoryImpl) AdminGetVaults(ctx context.Context, enabled *bool, limit *int, offset *int) ([]*model.Vault, error) {
	stmt := d.repo.Vault.
		WithContext(ctx)
	if enabled != nil {
		stmt = stmt.Where(d.repo.Vault.Enabled.Is(*enabled))
	}
	if limit != nil && *limit > 0 {
		stmt = stmt.Limit(*limit)
	}
	if offset != nil && *offset > 0 {
		stmt = stmt.Offset(*offset)
	}
	stmt = stmt.Order(d.repo.Vault.Pubkey)
	return stmt.Find()
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

func (d repositoryImpl) GetVaultPeriods(ctx context.Context, vault string, limit int, offset int, vaultPeriod *string) ([]*model.VaultPeriod, error) {
	stmt := d.repo.
		VaultPeriod.WithContext(ctx).
		Join(d.repo.Vault, d.repo.VaultPeriod.Vault.EqCol(d.repo.Vault.Pubkey)).
		Where(d.repo.VaultPeriod.Vault.Eq(vault)).
		Where(d.repo.Vault.Enabled.Is(true))
	if vaultPeriod != nil {
		stmt = stmt.Where(d.repo.VaultPeriod.Pubkey.Eq(*vaultPeriod))
	}
	return stmt.Limit(limit).Offset(offset).Find()
}
