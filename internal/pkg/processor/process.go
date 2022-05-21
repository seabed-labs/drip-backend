package processor

import (
	"context"

	"github.com/dcaf-protocol/drip/internal/pkg/repository/query"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/internal/pkg/repository/model"
	"gorm.io/gorm/clause"
)

type Processor interface {
	UpsertProtoConfigByAddress(context.Context, string) error
	UpsertProtoConfigs(context.Context, ...*model.ProtoConfig) error

	UpsertVaultByAddress(context.Context, string) error
	UpsertVaults(context.Context, ...*model.Vault) error

	UpsertPositionByAddress(context.Context, string) error
	UpsertPositions(context.Context, ...*model.Position) error

	UpsertVaultPeriodByAddress(context.Context, string) error
	UpsertVaultPeriods(context.Context, ...*model.VaultPeriod) error
}

type ProcessorImpl struct {
	repo   *query.Query
	client solana.Solana
}

func NewProcessor(
	repo *query.Query,
	client solana.Solana,
) Processor {
	return ProcessorImpl{
		repo:   repo,
		client: client,
	}
}

func (p ProcessorImpl) UpsertProtoConfigByAddress(ctx context.Context, s string) error {
	//TODO implement me
	panic("implement me")
}

func (p ProcessorImpl) UpsertProtoConfigs(ctx context.Context, protoConfigs ...*model.ProtoConfig) error {
	return p.repo.ProtoConfig.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(protoConfigs...)
}

func (p ProcessorImpl) UpsertVaultByAddress(ctx context.Context, s string) error {
	//TODO implement me
	panic("implement me")
}

func (p ProcessorImpl) UpsertVaults(ctx context.Context, vaults ...*model.Vault) error {
	return p.repo.Vault.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(vaults...)
}

func (p ProcessorImpl) UpsertPositionByAddress(ctx context.Context, s string) error {
	//TODO implement me
	panic("implement me")
}

func (p ProcessorImpl) UpsertPositions(ctx context.Context, positions ...*model.Position) error {
	return p.repo.Position.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(positions...)
}

func (p ProcessorImpl) UpsertVaultPeriodByAddress(ctx context.Context, s string) error {
	//TODO implement me
	panic("implement me")
}

func (p ProcessorImpl) UpsertVaultPeriods(ctx context.Context, vaultPeriods ...*model.VaultPeriod) error {
	return p.repo.VaultPeriod.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(vaultPeriods...)
}
