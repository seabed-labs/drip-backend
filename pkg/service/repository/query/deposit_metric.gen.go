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

	"gorm.io/plugin/dbresolver"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
)

func newDepositMetric(db *gorm.DB) depositMetric {
	_depositMetric := depositMetric{}

	_depositMetric.depositMetricDo.UseDB(db)
	_depositMetric.depositMetricDo.UseModel(&model.DepositMetric{})

	tableName := _depositMetric.depositMetricDo.TableName()
	_depositMetric.ALL = field.NewField(tableName, "*")
	_depositMetric.Signature = field.NewString(tableName, "signature")
	_depositMetric.IxIndex = field.NewInt32(tableName, "ix_index")
	_depositMetric.IxName = field.NewString(tableName, "ix_name")
	_depositMetric.IxVersion = field.NewInt32(tableName, "ix_version")
	_depositMetric.Slot = field.NewInt32(tableName, "slot")
	_depositMetric.Time = field.NewTime(tableName, "time")
	_depositMetric.Vault = field.NewString(tableName, "vault")
	_depositMetric.TokenAMint = field.NewString(tableName, "token_a_mint")
	_depositMetric.Referrer = field.NewString(tableName, "referrer")
	_depositMetric.TokenADepositAmount = field.NewUint64(tableName, "token_a_deposit_amount")
	_depositMetric.TokenAUsdPriceDay = field.NewUint64(tableName, "token_a_usd_price_day")

	_depositMetric.fillFieldMap()

	return _depositMetric
}

type depositMetric struct {
	depositMetricDo depositMetricDo

	ALL                 field.Field
	Signature           field.String
	IxIndex             field.Int32
	IxName              field.String
	IxVersion           field.Int32
	Slot                field.Int32
	Time                field.Time
	Vault               field.String
	TokenAMint          field.String
	Referrer            field.String
	TokenADepositAmount field.Uint64
	TokenAUsdPriceDay   field.Uint64

	fieldMap map[string]field.Expr
}

func (d depositMetric) Table(newTableName string) *depositMetric {
	d.depositMetricDo.UseTable(newTableName)
	return d.updateTableName(newTableName)
}

func (d depositMetric) As(alias string) *depositMetric {
	d.depositMetricDo.DO = *(d.depositMetricDo.As(alias).(*gen.DO))
	return d.updateTableName(alias)
}

func (d *depositMetric) updateTableName(table string) *depositMetric {
	d.ALL = field.NewField(table, "*")
	d.Signature = field.NewString(table, "signature")
	d.IxIndex = field.NewInt32(table, "ix_index")
	d.IxName = field.NewString(table, "ix_name")
	d.IxVersion = field.NewInt32(table, "ix_version")
	d.Slot = field.NewInt32(table, "slot")
	d.Time = field.NewTime(table, "time")
	d.Vault = field.NewString(table, "vault")
	d.TokenAMint = field.NewString(table, "token_a_mint")
	d.Referrer = field.NewString(table, "referrer")
	d.TokenADepositAmount = field.NewUint64(table, "token_a_deposit_amount")
	d.TokenAUsdPriceDay = field.NewUint64(table, "token_a_usd_price_day")

	d.fillFieldMap()

	return d
}

func (d *depositMetric) WithContext(ctx context.Context) *depositMetricDo {
	return d.depositMetricDo.WithContext(ctx)
}

func (d depositMetric) TableName() string { return d.depositMetricDo.TableName() }

func (d depositMetric) Alias() string { return d.depositMetricDo.Alias() }

func (d *depositMetric) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := d.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (d *depositMetric) fillFieldMap() {
	d.fieldMap = make(map[string]field.Expr, 11)
	d.fieldMap["signature"] = d.Signature
	d.fieldMap["ix_index"] = d.IxIndex
	d.fieldMap["ix_name"] = d.IxName
	d.fieldMap["ix_version"] = d.IxVersion
	d.fieldMap["slot"] = d.Slot
	d.fieldMap["time"] = d.Time
	d.fieldMap["vault"] = d.Vault
	d.fieldMap["token_a_mint"] = d.TokenAMint
	d.fieldMap["referrer"] = d.Referrer
	d.fieldMap["token_a_deposit_amount"] = d.TokenADepositAmount
	d.fieldMap["token_a_usd_price_day"] = d.TokenAUsdPriceDay
}

func (d depositMetric) clone(db *gorm.DB) depositMetric {
	d.depositMetricDo.ReplaceDB(db)
	return d
}

type depositMetricDo struct{ gen.DO }

func (d depositMetricDo) Debug() *depositMetricDo {
	return d.withDO(d.DO.Debug())
}

func (d depositMetricDo) WithContext(ctx context.Context) *depositMetricDo {
	return d.withDO(d.DO.WithContext(ctx))
}

func (d depositMetricDo) ReadDB() *depositMetricDo {
	return d.Clauses(dbresolver.Read)
}

func (d depositMetricDo) WriteDB() *depositMetricDo {
	return d.Clauses(dbresolver.Write)
}

func (d depositMetricDo) Clauses(conds ...clause.Expression) *depositMetricDo {
	return d.withDO(d.DO.Clauses(conds...))
}

func (d depositMetricDo) Returning(value interface{}, columns ...string) *depositMetricDo {
	return d.withDO(d.DO.Returning(value, columns...))
}

func (d depositMetricDo) Not(conds ...gen.Condition) *depositMetricDo {
	return d.withDO(d.DO.Not(conds...))
}

func (d depositMetricDo) Or(conds ...gen.Condition) *depositMetricDo {
	return d.withDO(d.DO.Or(conds...))
}

func (d depositMetricDo) Select(conds ...field.Expr) *depositMetricDo {
	return d.withDO(d.DO.Select(conds...))
}

func (d depositMetricDo) Where(conds ...gen.Condition) *depositMetricDo {
	return d.withDO(d.DO.Where(conds...))
}

func (d depositMetricDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *depositMetricDo {
	return d.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (d depositMetricDo) Order(conds ...field.Expr) *depositMetricDo {
	return d.withDO(d.DO.Order(conds...))
}

func (d depositMetricDo) Distinct(cols ...field.Expr) *depositMetricDo {
	return d.withDO(d.DO.Distinct(cols...))
}

func (d depositMetricDo) Omit(cols ...field.Expr) *depositMetricDo {
	return d.withDO(d.DO.Omit(cols...))
}

func (d depositMetricDo) Join(table schema.Tabler, on ...field.Expr) *depositMetricDo {
	return d.withDO(d.DO.Join(table, on...))
}

func (d depositMetricDo) LeftJoin(table schema.Tabler, on ...field.Expr) *depositMetricDo {
	return d.withDO(d.DO.LeftJoin(table, on...))
}

func (d depositMetricDo) RightJoin(table schema.Tabler, on ...field.Expr) *depositMetricDo {
	return d.withDO(d.DO.RightJoin(table, on...))
}

func (d depositMetricDo) Group(cols ...field.Expr) *depositMetricDo {
	return d.withDO(d.DO.Group(cols...))
}

func (d depositMetricDo) Having(conds ...gen.Condition) *depositMetricDo {
	return d.withDO(d.DO.Having(conds...))
}

func (d depositMetricDo) Limit(limit int) *depositMetricDo {
	return d.withDO(d.DO.Limit(limit))
}

func (d depositMetricDo) Offset(offset int) *depositMetricDo {
	return d.withDO(d.DO.Offset(offset))
}

func (d depositMetricDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *depositMetricDo {
	return d.withDO(d.DO.Scopes(funcs...))
}

func (d depositMetricDo) Unscoped() *depositMetricDo {
	return d.withDO(d.DO.Unscoped())
}

func (d depositMetricDo) Create(values ...*model.DepositMetric) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Create(values)
}

func (d depositMetricDo) CreateInBatches(values []*model.DepositMetric, batchSize int) error {
	return d.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (d depositMetricDo) Save(values ...*model.DepositMetric) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Save(values)
}

func (d depositMetricDo) First() (*model.DepositMetric, error) {
	if result, err := d.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.DepositMetric), nil
	}
}

func (d depositMetricDo) Take() (*model.DepositMetric, error) {
	if result, err := d.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.DepositMetric), nil
	}
}

func (d depositMetricDo) Last() (*model.DepositMetric, error) {
	if result, err := d.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.DepositMetric), nil
	}
}

func (d depositMetricDo) Find() ([]*model.DepositMetric, error) {
	result, err := d.DO.Find()
	return result.([]*model.DepositMetric), err
}

func (d depositMetricDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.DepositMetric, err error) {
	buf := make([]*model.DepositMetric, 0, batchSize)
	err = d.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (d depositMetricDo) FindInBatches(result *[]*model.DepositMetric, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return d.DO.FindInBatches(result, batchSize, fc)
}

func (d depositMetricDo) Attrs(attrs ...field.AssignExpr) *depositMetricDo {
	return d.withDO(d.DO.Attrs(attrs...))
}

func (d depositMetricDo) Assign(attrs ...field.AssignExpr) *depositMetricDo {
	return d.withDO(d.DO.Assign(attrs...))
}

func (d depositMetricDo) Joins(fields ...field.RelationField) *depositMetricDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Joins(_f))
	}
	return &d
}

func (d depositMetricDo) Preload(fields ...field.RelationField) *depositMetricDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Preload(_f))
	}
	return &d
}

func (d depositMetricDo) FirstOrInit() (*model.DepositMetric, error) {
	if result, err := d.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.DepositMetric), nil
	}
}

func (d depositMetricDo) FirstOrCreate() (*model.DepositMetric, error) {
	if result, err := d.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.DepositMetric), nil
	}
}

func (d depositMetricDo) FindByPage(offset int, limit int) (result []*model.DepositMetric, count int64, err error) {
	result, err = d.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = d.Offset(-1).Limit(-1).Count()
	return
}

func (d depositMetricDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = d.Count()
	if err != nil {
		return
	}

	err = d.Offset(offset).Limit(limit).Scan(result)
	return
}

func (d depositMetricDo) Scan(result interface{}) (err error) {
	return d.DO.Scan(result)
}

func (d depositMetricDo) Delete(models ...*model.DepositMetric) (result gen.ResultInfo, err error) {
	return d.DO.Delete(models)
}

func (d *depositMetricDo) withDO(do gen.Dao) *depositMetricDo {
	d.DO = *do.(*gen.DO)
	return d
}
