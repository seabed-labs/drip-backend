package repository

import (
	"context"
	"database/sql"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/internal/pkg/repository/model"
	"github.com/dcaf-protocol/drip/internal/pkg/repository/query"
)

type Repository interface {
	GetVaults(context.Context, *string, *string, *string) ([]*model.Vault, error)
	GetProtoConfigs(context.Context, *string, *string) ([]*model.ProtoConfig, error)
	GetVaultPeriods(context.Context, string, int, int, *string) ([]*model.VaultPeriod, error)
	GetTokensWithSupportedTokenPair(context.Context, *string, bool) ([]*model.Token, error)
	GetTokenPair(context.Context, string) (*model.TokenPair, error)
	GetTokenPairs(context.Context, *string, *string) ([]*model.TokenPair, error)
}

type repositoryImpl struct {
	client solana.Solana
	repo   *query.Query
}

func NewRepository(
	client solana.Solana,
	repo *query.Query,
) Repository {
	return repositoryImpl{
		client: client,
		repo:   repo,
	}
}

func (d repositoryImpl) GetTokenPair(ctx context.Context, id string) (*model.TokenPair, error) {
	tokenPairs, err := d.repo.TokenPair.WithContext(ctx).Where(d.repo.TokenPair.ID.Eq(id)).Limit(1).Find()
	if err != nil {
		return nil, err
	}
	if len(tokenPairs) > 0 {
		return tokenPairs[0], nil
	}
	return nil, sql.ErrNoRows
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

func (d repositoryImpl) GetVaults(ctx context.Context, tokenAMint, tokenBMint, protoConfig *string) ([]*model.Vault, error) {
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
