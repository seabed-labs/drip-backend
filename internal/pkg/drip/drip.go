package drip

import (
	"context"
	"database/sql"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/internal/pkg/repository"
	"github.com/dcaf-protocol/drip/internal/pkg/repository/model"
)

type Drip interface {
	GetVaults(context.Context, *string, *string, *string) ([]*model.Vault, error)
	GetProtoConfigs(context.Context, *string, *string) ([]*model.ProtoConfig, error)
	GetVaultPeriods(context.Context, string, int, int, *string) ([]*model.VaultPeriod, error)
	GetTokens(context.Context, *string, *string) ([]*model.Token, error)
	GetTokenPair(context.Context, string) (*model.TokenPair, error)
}

type dripImpl struct {
	client solana.Solana
	repo   *repository.Query
}

func (d dripImpl) GetTokenPair(ctx context.Context, id string) (*model.TokenPair, error) {
	tokenPairs, err := d.repo.TokenPair.WithContext(ctx).Where(d.repo.TokenPair.ID.Eq(id)).Limit(1).Find()
	if err != nil {
		return nil, err
	}
	if len(tokenPairs) > 0 {
		return tokenPairs[0], nil
	}
	return nil, sql.ErrNoRows
}

func (d dripImpl) GetTokens(ctx context.Context, tokenA *string, tokenB *string) ([]*model.Token, error) {
	//query := d.repo.Token.WithContext(ctx)
	//query = query.Join(d.repo.Vault, d.repo.ProtoConfig.Pubkey.EqCol(d.repo.Vault.ProtoConfig))
	//TODO implement me
	panic("implement me")
}

func NewDripService(
	client solana.Solana,
	repo *repository.Query,
) Drip {
	return dripImpl{
		client: client,
		repo:   repo,
	}
}

func (d dripImpl) GetVaults(ctx context.Context, tokenAMint, tokenBMint, protoConfig *string) ([]*model.Vault, error) {
	query := d.repo.Vault.WithContext(ctx)
	if tokenAMint != nil || tokenBMint != nil {
		query = query.Join(d.repo.Vault, d.repo.TokenPair.ID.EqCol(d.repo.Vault.TokenPairID))
	}
	if tokenAMint != nil {
		query = query.Where(d.repo.TokenPair.TokenA.Eq(*tokenAMint))
	}
	if tokenBMint != nil {
		query = query.Where(d.repo.TokenPair.TokenB.Eq(*tokenBMint))
	}
	if protoConfig != nil {
		query = query.Where(d.repo.Vault.ProtoConfig.Eq(*protoConfig))
	}
	query = query.Where(d.repo.Vault.Enabled.Is(true))
	return query.Find()
}

func (d dripImpl) GetProtoConfigs(ctx context.Context, tokenAMint *string, tokenBMint *string) ([]*model.ProtoConfig, error) {
	query := d.repo.ProtoConfig.WithContext(ctx)
	query = query.Join(d.repo.Vault, d.repo.ProtoConfig.Pubkey.EqCol(d.repo.Vault.ProtoConfig))

	if tokenAMint != nil || tokenBMint != nil {
		query = query.Join(d.repo.Vault, d.repo.TokenPair.ID.EqCol(d.repo.Vault.TokenPairID))
	}
	if tokenAMint != nil {
		query = query.Where(d.repo.TokenPair.TokenA.Eq(*tokenAMint))
	}
	if tokenBMint != nil {
		query = query.Where(d.repo.TokenPair.TokenB.Eq(*tokenBMint))
	}
	query = query.Where(d.repo.Vault.Enabled.Is(true))
	return query.Find()
}

func (d dripImpl) GetVaultPeriods(ctx context.Context, vault string, limit int, offset int, vaultPeriod *string) ([]*model.VaultPeriod, error) {
	query := d.repo.
		VaultPeriod.WithContext(ctx).
		Join(d.repo.Vault, d.repo.VaultPeriod.Vault.EqCol(d.repo.Vault.Pubkey)).
		Where(d.repo.VaultPeriod.Vault.Eq(vault)).
		Where(d.repo.Vault.Enabled.Is(true))
	if vaultPeriod != nil {
		query = query.Where(d.repo.VaultPeriod.Pubkey.Eq(*vaultPeriod))
	}
	return query.Limit(limit).Offset(offset).Find()
}
