package drip

import (
	"context"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/internal/pkg/repository"
	"github.com/dcaf-protocol/drip/internal/pkg/repository/model"
)

type Drip interface {
	GetVaults(context.Context, *string, *string, *string) ([]*model.Vault, error)
	GetProtoConfigs(context.Context, *string, *string) ([]*model.ProtoConfig, error)
}

type dripImpl struct {
	client solana.Solana
	repo   *repository.Query
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
	if tokenAMint != nil {
		query = query.Where(d.repo.Vault.TokenAMint.Eq(*tokenAMint))
	}
	if tokenBMint != nil {
		query = query.Where(d.repo.Vault.TokenBMint.Eq(*tokenBMint))
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
	if tokenAMint != nil {
		query = query.Where(d.repo.Vault.TokenAMint.Eq(*tokenAMint))
	}
	if tokenBMint != nil {
		query = query.Where(d.repo.Vault.TokenBMint.Eq(*tokenBMint))
	}
	query = query.Where(d.repo.Vault.Enabled.Is(true))
	return query.Find()
}
