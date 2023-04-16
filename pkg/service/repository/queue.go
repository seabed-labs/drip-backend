package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	"gorm.io/gorm/clause"
)

func (d repositoryImpl) AddAccountUpdateQueueItem(ctx context.Context, item *model.AccountUpdateQueueItem) error {
	return d.repo.AccountUpdateQueueItem.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "pubkey"}},
			UpdateAll: true,
		}).
		Create(item)
}

func (d repositoryImpl) ReQueueAccountUpdateQueueItem(ctx context.Context, item *model.AccountUpdateQueueItem) error {
	item.Try += 1
	if item.Try > *item.MaxTry {
		return fmt.Errorf("max retry already attempted")
	}
	retryTimeUnix := time.Now().Unix()
	if item.RetryTime != nil {
		retryTimeUnix = item.RetryTime.Unix()
	}
	retryTimeUnix += 30
	item.RetryTime = utils.GetTimePtr(time.Unix(retryTimeUnix, 0))
	return d.AddAccountUpdateQueueItem(ctx, item)
}

func (d repositoryImpl) PopAccountUpdateQueueItem(ctx context.Context) (*model.AccountUpdateQueueItem, error) {
	queueItem, err := d.repo.AccountUpdateQueueItem.
		WithContext(ctx).
		Where(d.repo.AccountUpdateQueueItem.RetryTime.Lt(time.Now())).
		Or(d.repo.AccountUpdateQueueItem.RetryTime.IsNull()).
		Order(d.repo.AccountUpdateQueueItem.Priority, d.repo.AccountUpdateQueueItem.Time).
		Limit(1).
		First()
	if err != nil {
		return nil, err
	}
	return queueItem, d.removeAccountUpdateQueueItem(ctx, queueItem)
}

func (d repositoryImpl) AccountUpdateQueueItemDepth(ctx context.Context) (int64, error) {
	return d.repo.AccountUpdateQueueItem.WithContext(ctx).Count()
}

func (d repositoryImpl) removeAccountUpdateQueueItem(ctx context.Context, item *model.AccountUpdateQueueItem) error {
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

func (d repositoryImpl) AddTransactionUpdateQueueItem(ctx context.Context, item *model.TransactionUpdateQueueItem) error {
	return d.repo.TransactionUpdateQueueItem.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "signature"}},
			UpdateAll: true,
		}).
		Create(item)
}

func (d repositoryImpl) ReQueueTransactionUpdateItem(ctx context.Context, item *model.TransactionUpdateQueueItem) error {
	item.Try += 1
	if item.Try > *item.MaxTry {
		return fmt.Errorf("max retry already attempted")
	}
	retryTimeUnix := time.Now().Unix()
	if item.RetryTime != nil {
		retryTimeUnix = item.RetryTime.Unix()
	}
	retryTimeUnix += 30
	item.RetryTime = utils.GetTimePtr(time.Unix(retryTimeUnix, 0))
	return d.AddTransactionUpdateQueueItem(ctx, item)
}

func (d repositoryImpl) TransactionUpdateQueueItemDepth(ctx context.Context) (int64, error) {
	return d.repo.TransactionUpdateQueueItem.WithContext(ctx).Count()
}

func (d repositoryImpl) PopTransactionUpdateQueueItem(ctx context.Context) (*model.TransactionUpdateQueueItem, error) {
	queueItem, err := d.repo.TransactionUpdateQueueItem.
		WithContext(ctx).
		Where(d.repo.TransactionUpdateQueueItem.RetryTime.Lt(time.Now())).
		Or(d.repo.TransactionUpdateQueueItem.RetryTime.IsNull()).
		Order(d.repo.TransactionUpdateQueueItem.Priority, d.repo.TransactionUpdateQueueItem.Time).
		Limit(1).
		First()
	if err != nil {
		return nil, err
	}
	return queueItem, d.removeTransactionUpdateQueueItem(ctx, queueItem)
}

func (d repositoryImpl) removeTransactionUpdateQueueItem(ctx context.Context, item *model.TransactionUpdateQueueItem) error {
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
