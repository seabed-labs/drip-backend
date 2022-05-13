// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"github.com/dcaf-protocol/drip/internal/pkg/repository/model"
)

func newUserPosition(db *gorm.DB) userPosition {
	_userPosition := userPosition{}

	_userPosition.userPositionDo.UseDB(db)
	_userPosition.userPositionDo.UseModel(&model.UserPosition{})

	tableName := _userPosition.userPositionDo.TableName()
	_userPosition.ALL = field.NewField(tableName, "*")
	_userPosition.Pubkey = field.NewString(tableName, "pubkey")
	_userPosition.Mint = field.NewString(tableName, "mint")
	_userPosition.Amount = field.NewBool(tableName, "amount")

	_userPosition.fillFieldMap()

	return _userPosition
}

type userPosition struct {
	userPositionDo userPositionDo

	ALL    field.Field
	Pubkey field.String
	Mint   field.String
	Amount field.Bool

	fieldMap map[string]field.Expr
}

func (u userPosition) Table(newTableName string) *userPosition {
	u.userPositionDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userPosition) As(alias string) *userPosition {
	u.userPositionDo.DO = *(u.userPositionDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userPosition) updateTableName(table string) *userPosition {
	u.ALL = field.NewField(table, "*")
	u.Pubkey = field.NewString(table, "pubkey")
	u.Mint = field.NewString(table, "mint")
	u.Amount = field.NewBool(table, "amount")

	u.fillFieldMap()

	return u
}

func (u *userPosition) WithContext(ctx context.Context) *userPositionDo {
	return u.userPositionDo.WithContext(ctx)
}

func (u userPosition) TableName() string { return u.userPositionDo.TableName() }

func (u userPosition) Alias() string { return u.userPositionDo.Alias() }

func (u *userPosition) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userPosition) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 3)
	u.fieldMap["pubkey"] = u.Pubkey
	u.fieldMap["mint"] = u.Mint
	u.fieldMap["amount"] = u.Amount
}

func (u userPosition) clone(db *gorm.DB) userPosition {
	u.userPositionDo.ReplaceDB(db)
	return u
}

type userPositionDo struct{ gen.DO }

func (u userPositionDo) Debug() *userPositionDo {
	return u.withDO(u.DO.Debug())
}

func (u userPositionDo) WithContext(ctx context.Context) *userPositionDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userPositionDo) Clauses(conds ...clause.Expression) *userPositionDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userPositionDo) Returning(value interface{}, columns ...string) *userPositionDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userPositionDo) Not(conds ...gen.Condition) *userPositionDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userPositionDo) Or(conds ...gen.Condition) *userPositionDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userPositionDo) Select(conds ...field.Expr) *userPositionDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userPositionDo) Where(conds ...gen.Condition) *userPositionDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userPositionDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *userPositionDo {
	return u.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (u userPositionDo) Order(conds ...field.Expr) *userPositionDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userPositionDo) Distinct(cols ...field.Expr) *userPositionDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userPositionDo) Omit(cols ...field.Expr) *userPositionDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userPositionDo) Join(table schema.Tabler, on ...field.Expr) *userPositionDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userPositionDo) LeftJoin(table schema.Tabler, on ...field.Expr) *userPositionDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userPositionDo) RightJoin(table schema.Tabler, on ...field.Expr) *userPositionDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userPositionDo) Group(cols ...field.Expr) *userPositionDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userPositionDo) Having(conds ...gen.Condition) *userPositionDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userPositionDo) Limit(limit int) *userPositionDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userPositionDo) Offset(offset int) *userPositionDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userPositionDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *userPositionDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userPositionDo) Unscoped() *userPositionDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userPositionDo) Create(values ...*model.UserPosition) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userPositionDo) CreateInBatches(values []*model.UserPosition, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userPositionDo) Save(values ...*model.UserPosition) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userPositionDo) First() (*model.UserPosition, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserPosition), nil
	}
}

func (u userPositionDo) Take() (*model.UserPosition, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserPosition), nil
	}
}

func (u userPositionDo) Last() (*model.UserPosition, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserPosition), nil
	}
}

func (u userPositionDo) Find() ([]*model.UserPosition, error) {
	result, err := u.DO.Find()
	return result.([]*model.UserPosition), err
}

func (u userPositionDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserPosition, err error) {
	buf := make([]*model.UserPosition, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userPositionDo) FindInBatches(result *[]*model.UserPosition, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userPositionDo) Attrs(attrs ...field.AssignExpr) *userPositionDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userPositionDo) Assign(attrs ...field.AssignExpr) *userPositionDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userPositionDo) Joins(fields ...field.RelationField) *userPositionDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userPositionDo) Preload(fields ...field.RelationField) *userPositionDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userPositionDo) FirstOrInit() (*model.UserPosition, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserPosition), nil
	}
}

func (u userPositionDo) FirstOrCreate() (*model.UserPosition, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserPosition), nil
	}
}

func (u userPositionDo) FindByPage(offset int, limit int) (result []*model.UserPosition, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u userPositionDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u *userPositionDo) withDO(do gen.Dao) *userPositionDo {
	u.DO = *do.(*gen.DO)
	return u
}
