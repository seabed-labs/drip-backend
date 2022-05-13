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

func newVaultPeriod(db *gorm.DB) vaultPeriod {
	_vaultPeriod := vaultPeriod{}

	_vaultPeriod.vaultPeriodDo.UseDB(db)
	_vaultPeriod.vaultPeriodDo.UseModel(&model.VaultPeriod{})

	tableName := _vaultPeriod.vaultPeriodDo.TableName()
	_vaultPeriod.ALL = field.NewField(tableName, "*")
	_vaultPeriod.Pubkey = field.NewString(tableName, "pubkey")
	_vaultPeriod.Vault = field.NewString(tableName, "vault")
	_vaultPeriod.PeriodID = field.NewFloat64(tableName, "period_id")
	_vaultPeriod.Twap = field.NewFloat64(tableName, "twap")
	_vaultPeriod.Dar = field.NewFloat64(tableName, "dar")

	_vaultPeriod.fillFieldMap()

	return _vaultPeriod
}

type vaultPeriod struct {
	vaultPeriodDo vaultPeriodDo

	ALL      field.Field
	Pubkey   field.String
	Vault    field.String
	PeriodID field.Float64
	Twap     field.Float64
	Dar      field.Float64

	fieldMap map[string]field.Expr
}

func (v vaultPeriod) Table(newTableName string) *vaultPeriod {
	v.vaultPeriodDo.UseTable(newTableName)
	return v.updateTableName(newTableName)
}

func (v vaultPeriod) As(alias string) *vaultPeriod {
	v.vaultPeriodDo.DO = *(v.vaultPeriodDo.As(alias).(*gen.DO))
	return v.updateTableName(alias)
}

func (v *vaultPeriod) updateTableName(table string) *vaultPeriod {
	v.ALL = field.NewField(table, "*")
	v.Pubkey = field.NewString(table, "pubkey")
	v.Vault = field.NewString(table, "vault")
	v.PeriodID = field.NewFloat64(table, "period_id")
	v.Twap = field.NewFloat64(table, "twap")
	v.Dar = field.NewFloat64(table, "dar")

	v.fillFieldMap()

	return v
}

func (v *vaultPeriod) WithContext(ctx context.Context) *vaultPeriodDo {
	return v.vaultPeriodDo.WithContext(ctx)
}

func (v vaultPeriod) TableName() string { return v.vaultPeriodDo.TableName() }

func (v vaultPeriod) Alias() string { return v.vaultPeriodDo.Alias() }

func (v *vaultPeriod) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := v.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (v *vaultPeriod) fillFieldMap() {
	v.fieldMap = make(map[string]field.Expr, 5)
	v.fieldMap["pubkey"] = v.Pubkey
	v.fieldMap["vault"] = v.Vault
	v.fieldMap["period_id"] = v.PeriodID
	v.fieldMap["twap"] = v.Twap
	v.fieldMap["dar"] = v.Dar
}

func (v vaultPeriod) clone(db *gorm.DB) vaultPeriod {
	v.vaultPeriodDo.ReplaceDB(db)
	return v
}

type vaultPeriodDo struct{ gen.DO }

func (v vaultPeriodDo) Debug() *vaultPeriodDo {
	return v.withDO(v.DO.Debug())
}

func (v vaultPeriodDo) WithContext(ctx context.Context) *vaultPeriodDo {
	return v.withDO(v.DO.WithContext(ctx))
}

func (v vaultPeriodDo) Clauses(conds ...clause.Expression) *vaultPeriodDo {
	return v.withDO(v.DO.Clauses(conds...))
}

func (v vaultPeriodDo) Returning(value interface{}, columns ...string) *vaultPeriodDo {
	return v.withDO(v.DO.Returning(value, columns...))
}

func (v vaultPeriodDo) Not(conds ...gen.Condition) *vaultPeriodDo {
	return v.withDO(v.DO.Not(conds...))
}

func (v vaultPeriodDo) Or(conds ...gen.Condition) *vaultPeriodDo {
	return v.withDO(v.DO.Or(conds...))
}

func (v vaultPeriodDo) Select(conds ...field.Expr) *vaultPeriodDo {
	return v.withDO(v.DO.Select(conds...))
}

func (v vaultPeriodDo) Where(conds ...gen.Condition) *vaultPeriodDo {
	return v.withDO(v.DO.Where(conds...))
}

func (v vaultPeriodDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *vaultPeriodDo {
	return v.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (v vaultPeriodDo) Order(conds ...field.Expr) *vaultPeriodDo {
	return v.withDO(v.DO.Order(conds...))
}

func (v vaultPeriodDo) Distinct(cols ...field.Expr) *vaultPeriodDo {
	return v.withDO(v.DO.Distinct(cols...))
}

func (v vaultPeriodDo) Omit(cols ...field.Expr) *vaultPeriodDo {
	return v.withDO(v.DO.Omit(cols...))
}

func (v vaultPeriodDo) Join(table schema.Tabler, on ...field.Expr) *vaultPeriodDo {
	return v.withDO(v.DO.Join(table, on...))
}

func (v vaultPeriodDo) LeftJoin(table schema.Tabler, on ...field.Expr) *vaultPeriodDo {
	return v.withDO(v.DO.LeftJoin(table, on...))
}

func (v vaultPeriodDo) RightJoin(table schema.Tabler, on ...field.Expr) *vaultPeriodDo {
	return v.withDO(v.DO.RightJoin(table, on...))
}

func (v vaultPeriodDo) Group(cols ...field.Expr) *vaultPeriodDo {
	return v.withDO(v.DO.Group(cols...))
}

func (v vaultPeriodDo) Having(conds ...gen.Condition) *vaultPeriodDo {
	return v.withDO(v.DO.Having(conds...))
}

func (v vaultPeriodDo) Limit(limit int) *vaultPeriodDo {
	return v.withDO(v.DO.Limit(limit))
}

func (v vaultPeriodDo) Offset(offset int) *vaultPeriodDo {
	return v.withDO(v.DO.Offset(offset))
}

func (v vaultPeriodDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *vaultPeriodDo {
	return v.withDO(v.DO.Scopes(funcs...))
}

func (v vaultPeriodDo) Unscoped() *vaultPeriodDo {
	return v.withDO(v.DO.Unscoped())
}

func (v vaultPeriodDo) Create(values ...*model.VaultPeriod) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Create(values)
}

func (v vaultPeriodDo) CreateInBatches(values []*model.VaultPeriod, batchSize int) error {
	return v.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (v vaultPeriodDo) Save(values ...*model.VaultPeriod) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Save(values)
}

func (v vaultPeriodDo) First() (*model.VaultPeriod, error) {
	if result, err := v.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.VaultPeriod), nil
	}
}

func (v vaultPeriodDo) Take() (*model.VaultPeriod, error) {
	if result, err := v.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.VaultPeriod), nil
	}
}

func (v vaultPeriodDo) Last() (*model.VaultPeriod, error) {
	if result, err := v.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.VaultPeriod), nil
	}
}

func (v vaultPeriodDo) Find() ([]*model.VaultPeriod, error) {
	result, err := v.DO.Find()
	return result.([]*model.VaultPeriod), err
}

func (v vaultPeriodDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.VaultPeriod, err error) {
	buf := make([]*model.VaultPeriod, 0, batchSize)
	err = v.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (v vaultPeriodDo) FindInBatches(result *[]*model.VaultPeriod, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return v.DO.FindInBatches(result, batchSize, fc)
}

func (v vaultPeriodDo) Attrs(attrs ...field.AssignExpr) *vaultPeriodDo {
	return v.withDO(v.DO.Attrs(attrs...))
}

func (v vaultPeriodDo) Assign(attrs ...field.AssignExpr) *vaultPeriodDo {
	return v.withDO(v.DO.Assign(attrs...))
}

func (v vaultPeriodDo) Joins(fields ...field.RelationField) *vaultPeriodDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Joins(_f))
	}
	return &v
}

func (v vaultPeriodDo) Preload(fields ...field.RelationField) *vaultPeriodDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Preload(_f))
	}
	return &v
}

func (v vaultPeriodDo) FirstOrInit() (*model.VaultPeriod, error) {
	if result, err := v.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.VaultPeriod), nil
	}
}

func (v vaultPeriodDo) FirstOrCreate() (*model.VaultPeriod, error) {
	if result, err := v.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.VaultPeriod), nil
	}
}

func (v vaultPeriodDo) FindByPage(offset int, limit int) (result []*model.VaultPeriod, count int64, err error) {
	result, err = v.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = v.Offset(-1).Limit(-1).Count()
	return
}

func (v vaultPeriodDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = v.Count()
	if err != nil {
		return
	}

	err = v.Offset(offset).Limit(limit).Scan(result)
	return
}

func (v *vaultPeriodDo) withDO(do gen.Dao) *vaultPeriodDo {
	v.DO = *do.(*gen.DO)
	return v
}
