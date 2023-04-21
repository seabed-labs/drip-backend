package repository

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/repository/query"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

type TransactionProcessingCheckpointRepository interface {
	GetLatestTransactionCheckpoint(ctx context.Context) *model.TransactionProcessingCheckpoint
	UpsertTransactionProcessingCheckpoint(ctx context.Context, slot uint64, signature string) error
}

type transactionProcessingCheckpointImpl struct {
	repo *query.Query
	db   *sqlx.DB
}

func NewTransactionProcessingCheckpointRepository(
	repo *query.Query,
	db *sqlx.DB,
) TransactionProcessingCheckpointRepository {
	return transactionProcessingCheckpointImpl{
		repo: repo,
		db:   db,
	}
}

func (d transactionProcessingCheckpointImpl) GetLatestTransactionCheckpoint(ctx context.Context) *model.TransactionProcessingCheckpoint {
	checkpoint, err := d.repo.
		TransactionProcessingCheckpoint.
		WithContext(ctx).
		Order(d.repo.TransactionProcessingCheckpoint.Slot.Desc()).
		First()
	if err != nil && err.Error() == ErrRecordNotFound {
		return nil
	} else if err != nil {
		logrus.WithError(err).Error("failed to get checkpoint")
	}
	return checkpoint
}

func (d transactionProcessingCheckpointImpl) UpsertTransactionProcessingCheckpoint(ctx context.Context, slot uint64, signature string) error {
	return d.repo.TransactionProcessingCheckpoint.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&model.TransactionProcessingCheckpoint{
			Signature: signature,
			Slot:      slot,
		})
}
