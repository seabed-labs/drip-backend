package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/dcaf-labs/drip/pkg/service/repository/query"
	"github.com/jmoiron/sqlx"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
)

type AccountUpdateQueue interface {
	AddAccountUpdateQueueItem(ctx context.Context, item *model.AccountUpdateQueueItem) error
	ReQueueAccountUpdateQueueItem(ctx context.Context, item *model.AccountUpdateQueueItem, reason string) error
	AccountUpdateQueueItemDepth(ctx context.Context) (int64, error)
	PopAccountUpdateQueueItem(ctx context.Context) (*model.AccountUpdateQueueItem, error)
}

type TransactionUpdateQueue interface {
	AddTransactionUpdateQueueItem(ctx context.Context, item *model.TransactionUpdateQueueItem) error
	ReQueueTransactionUpdateItem(ctx context.Context, item *model.TransactionUpdateQueueItem, reason string) error
	TransactionUpdateQueueItemDepth(ctx context.Context) (int64, error)
	PopTransactionUpdateQueueItem(ctx context.Context) (*model.TransactionUpdateQueueItem, error)
}

type queueImpl struct {
	repo *query.Query
	db   *sqlx.DB
}

func NewAccountUpdateQueue(
	repo *query.Query,
	db *sqlx.DB,
) AccountUpdateQueue {
	return queueImpl{
		repo: repo,
		db:   db,
	}
}

func NewTransactionUpdateQueue(
	repo *query.Query,
	db *sqlx.DB,
) TransactionUpdateQueue {
	return queueImpl{
		repo: repo,
		db:   db,
	}
}

func (d queueImpl) AddAccountUpdateQueueItem(ctx context.Context, item *model.AccountUpdateQueueItem) error {
	return d.repo.AccountUpdateQueueItem.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "pubkey"}},
			UpdateAll: true,
		}).
		Create(item)
}

func (d queueImpl) ReQueueAccountUpdateQueueItem(ctx context.Context, item *model.AccountUpdateQueueItem, reason string) error {
	item.Try += 1
	retryTimeUnix := time.Now().Unix()
	if item.RetryTime != nil {
		retryTimeUnix = item.RetryTime.Unix()
	}
	retryReason := reason
	if item.Reason != nil {
		retryReason = fmt.Sprintf("%s\n%s", *item.Reason, retryReason)
	}
	item.Reason = pointer.ToString(retryReason)
	retryTimeUnix += 30
	item.RetryTime = utils.GetTimePtr(time.Unix(retryTimeUnix, 0))
	return d.AddAccountUpdateQueueItem(ctx, item)
}

func (d queueImpl) PopAccountUpdateQueueItem(ctx context.Context) (*model.AccountUpdateQueueItem, error) {
	queueItem, err := d.repo.AccountUpdateQueueItem.
		WithContext(ctx).
		Where(
			d.repo.AccountUpdateQueueItem.Try.LteCol(d.repo.AccountUpdateQueueItem.MaxTry),
			field.Or(
				d.repo.AccountUpdateQueueItem.RetryTime.Lt(time.Now()),
				d.repo.AccountUpdateQueueItem.RetryTime.IsNull(),
			),
		).
		Order(d.repo.AccountUpdateQueueItem.Priority, d.repo.AccountUpdateQueueItem.Time).
		Limit(1).
		First()
	if err != nil {
		return nil, err
	}
	return queueItem, d.removeAccountUpdateQueueItem(ctx, queueItem)
}

func (d queueImpl) AccountUpdateQueueItemDepth(ctx context.Context) (int64, error) {
	return d.repo.AccountUpdateQueueItem.WithContext(ctx).Count()
}

func (d queueImpl) removeAccountUpdateQueueItem(ctx context.Context, item *model.AccountUpdateQueueItem) error {
	info, err := d.repo.AccountUpdateQueueItem.
		WithContext(ctx).
		Where(d.repo.AccountUpdateQueueItem.Pubkey.Eq(item.Pubkey)).
		Delete()
	if err != nil {
		return err
	}
	if info.Error != nil {
		return info.Error
	}
	return nil
}

func (d queueImpl) AddTransactionUpdateQueueItem(ctx context.Context, item *model.TransactionUpdateQueueItem) error {
	return d.repo.TransactionUpdateQueueItem.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "signature"}},
			UpdateAll: true,
		}).
		Create(item)
}

func (d queueImpl) ReQueueTransactionUpdateItem(ctx context.Context, item *model.TransactionUpdateQueueItem, reason string) error {
	item.Try += 1
	retryTimeUnix := time.Now().Unix()
	if item.RetryTime != nil {
		retryTimeUnix = item.RetryTime.Unix()
	}
	retryReason := reason
	if item.Reason != nil {
		retryReason = fmt.Sprintf("%s\n%s", *item.Reason, retryReason)
	}
	item.Reason = pointer.ToString(retryReason)
	retryTimeUnix += 30
	item.RetryTime = utils.GetTimePtr(time.Unix(retryTimeUnix, 0))
	return d.AddTransactionUpdateQueueItem(ctx, item)
}

func (d queueImpl) TransactionUpdateQueueItemDepth(ctx context.Context) (int64, error) {
	return d.repo.TransactionUpdateQueueItem.WithContext(ctx).Count()
}

func (d queueImpl) PopTransactionUpdateQueueItem(ctx context.Context) (*model.TransactionUpdateQueueItem, error) {
	queueItem, err := d.repo.TransactionUpdateQueueItem.
		WithContext(ctx).
		Where(
			d.repo.TransactionUpdateQueueItem.Try.LteCol(d.repo.TransactionUpdateQueueItem.MaxTry),
			field.Or(
				d.repo.TransactionUpdateQueueItem.RetryTime.Lt(time.Now()),
				d.repo.TransactionUpdateQueueItem.RetryTime.IsNull(),
			),
		).
		Order(d.repo.TransactionUpdateQueueItem.Priority, d.repo.TransactionUpdateQueueItem.Time).
		Limit(1).
		First()
	if err != nil {
		return nil, err
	}
	return queueItem, d.removeTransactionUpdateQueueItem(ctx, queueItem)
}

func (d queueImpl) removeTransactionUpdateQueueItem(ctx context.Context, item *model.TransactionUpdateQueueItem) error {
	info, err := d.repo.TransactionUpdateQueueItem.
		WithContext(ctx).
		Where(d.repo.TransactionUpdateQueueItem.Signature.Eq(item.Signature)).
		Delete()
	if err != nil {
		return err
	}
	if info.Error != nil {
		return info.Error
	}
	return nil
}
