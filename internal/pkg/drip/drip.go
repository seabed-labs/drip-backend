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
	GetVaults(ctx context.Context) ([]*model.Vault, error)
}

type dripImpl struct {
	client solana.Solana
	repo   repository.Query
}

func NewDripService(
	client solana.Solana,
	repo repository.Query,
) Drip {
	return dripImpl{
		client: client,
		repo:   repo,
	}
}

func (d dripImpl) GetVaults(ctx context.Context) ([]*model.Vault, error) {
	return d.repo.Vault.
		WithContext(ctx).
		Find()
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
