package repository

import (
	context "context"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
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

func (d repositoryImpl) RemoveItem(ctx context.Context, item *model.AccountUpdateQueueItem) error {
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

func (d repositoryImpl) Pop(ctx context.Context) (*model.AccountUpdateQueueItem, error) {
	queueitem, err := d.repo.AccountUpdateQueueItem.
		WithContext(ctx).
		Order(d.repo.AccountUpdateQueueItem.Priority, d.repo.AccountUpdateQueueItem.Time).
		First()
	if err != nil {
		return nil, err
	}
	return queueitem, d.RemoveItem(ctx, queueitem)
}
func (d repositoryImpl) Depth(ctx context.Context) (int64, error) {
	return d.repo.AccountUpdateQueueItem.WithContext(ctx).Count()
}
