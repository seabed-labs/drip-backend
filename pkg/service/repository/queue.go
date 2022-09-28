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
			Columns:   []clause.Column{{Name: "pubkey"}, {Name: "program_id"}},
			UpdateAll: true,
		}).
		Create(item)
}

func (d repositoryImpl) RemoveItem(ctx context.Context, item *model.AccountUpdateQueueItem) error {
	info, err := d.repo.AccountUpdateQueueItem.
		WithContext(ctx).
		Where(
			d.repo.AccountUpdateQueueItem.Pubkey.Eq(item.Pubkey),
			d.repo.AccountUpdateQueueItem.ProgramID.Eq(item.ProgramID),
		).Delete()
	if err != nil {
		return err
	}
	if info.Error != nil {
		return info.Error
	}
	return nil
}

func (d repositoryImpl) GetNextItem(ctx context.Context) (*model.AccountUpdateQueueItem, error) {
	return d.repo.AccountUpdateQueueItem.
		WithContext(ctx).
		Order(d.repo.AccountUpdateQueueItem.Time).
		First()
}
