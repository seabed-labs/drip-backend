package drip

import (
	"context"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/internal/pkg/repository"
	"github.com/dcaf-protocol/drip/internal/pkg/repository/model"
	"gorm.io/gorm/clause"
)

type Drip interface {
	UpsertVaultByAddress(context.Context, string) error
	UpsertVaults(context.Context, ...*model.Vault) error
	GetVaults(context.Context, *string, *string, *string) ([]*model.Vault, error)
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

func (d dripImpl) UpsertVaultByAddress(ctx context.Context, s string) error {
	//TODO implement me
	panic("implement me")
}

func (d dripImpl) UpsertVaults(ctx context.Context, vaults ...*model.Vault) error {
	return d.repo.Vault.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(vaults...)
}
