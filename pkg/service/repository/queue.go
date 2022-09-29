package repository

import (
	context "context"
	"fmt"
	"math"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	"gorm.io/gorm/clause"
)

func (d repositoryImpl) AddItem(ctx context.Context, item *model.AccountUpdateQueueItem) error {
	return d.repo.AccountUpdateQueueItem.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "pubkey"}},
			UpdateAll: true,
		}).
		Create(item)
}

func (d repositoryImpl) ReQueue(ctx context.Context, item *model.AccountUpdateQueueItem) error {
	item.Try += 1
	if item.Try > *item.MaxTry {
		return fmt.Errorf("max retry already attempted")
	}
	retryTimeUnix := time.Now().Unix()
	if item.RetryTime != nil {
		retryTimeUnix = item.RetryTime.Unix()
	}
	retryTimeUnix += int64(math.Pow(2, float64(item.Try)))
	item.RetryTime = utils.GetTimePtr(time.Unix(retryTimeUnix, 0))
	return d.AddItem(ctx, item)
}

func (d repositoryImpl) Pop(ctx context.Context) (*model.AccountUpdateQueueItem, error) {
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
	return queueItem, d.removeItem(ctx, queueItem)
}
func (d repositoryImpl) Depth(ctx context.Context) (int64, error) {
	return d.repo.AccountUpdateQueueItem.WithContext(ctx).Count()
}

func (d repositoryImpl) removeItem(ctx context.Context, item *model.AccountUpdateQueueItem) error {
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
