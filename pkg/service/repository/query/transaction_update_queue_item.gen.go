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

func newTransactionUpdateQueueItem(db *gorm.DB) transactionUpdateQueueItem {
	_transactionUpdateQueueItem := transactionUpdateQueueItem{}

	_transactionUpdateQueueItem.transactionUpdateQueueItemDo.UseDB(db)
	_transactionUpdateQueueItem.transactionUpdateQueueItemDo.UseModel(&model.TransactionUpdateQueueItem{})

	tableName := _transactionUpdateQueueItem.transactionUpdateQueueItemDo.TableName()
	_transactionUpdateQueueItem.ALL = field.NewField(tableName, "*")
	_transactionUpdateQueueItem.Signature = field.NewString(tableName, "signature")
	_transactionUpdateQueueItem.TxJSON = field.NewString(tableName, "tx_json")
	_transactionUpdateQueueItem.Time = field.NewTime(tableName, "time")
	_transactionUpdateQueueItem.Priority = field.NewInt32(tableName, "priority")
	_transactionUpdateQueueItem.Try = field.NewInt32(tableName, "try")
	_transactionUpdateQueueItem.MaxTry = field.NewInt32(tableName, "max_try")
	_transactionUpdateQueueItem.RetryTime = field.NewTime(tableName, "retry_time")
	_transactionUpdateQueueItem.Reason = field.NewString(tableName, "reason")

	_transactionUpdateQueueItem.fillFieldMap()

	return _transactionUpdateQueueItem
}

type transactionUpdateQueueItem struct {
	transactionUpdateQueueItemDo transactionUpdateQueueItemDo

	ALL       field.Field
	Signature field.String
	TxJSON    field.String
	Time      field.Time
	Priority  field.Int32
	Try       field.Int32
	MaxTry    field.Int32
	RetryTime field.Time
	Reason    field.String

	fieldMap map[string]field.Expr
}

func (t transactionUpdateQueueItem) Table(newTableName string) *transactionUpdateQueueItem {
	t.transactionUpdateQueueItemDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t transactionUpdateQueueItem) As(alias string) *transactionUpdateQueueItem {
	t.transactionUpdateQueueItemDo.DO = *(t.transactionUpdateQueueItemDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *transactionUpdateQueueItem) updateTableName(table string) *transactionUpdateQueueItem {
	t.ALL = field.NewField(table, "*")
	t.Signature = field.NewString(table, "signature")
	t.TxJSON = field.NewString(table, "tx_json")
	t.Time = field.NewTime(table, "time")
	t.Priority = field.NewInt32(table, "priority")
	t.Try = field.NewInt32(table, "try")
	t.MaxTry = field.NewInt32(table, "max_try")
	t.RetryTime = field.NewTime(table, "retry_time")
	t.Reason = field.NewString(table, "reason")

	t.fillFieldMap()

	return t
}

func (t *transactionUpdateQueueItem) WithContext(ctx context.Context) *transactionUpdateQueueItemDo {
	return t.transactionUpdateQueueItemDo.WithContext(ctx)
}

func (t transactionUpdateQueueItem) TableName() string {
	return t.transactionUpdateQueueItemDo.TableName()
}

func (t transactionUpdateQueueItem) Alias() string { return t.transactionUpdateQueueItemDo.Alias() }

func (t *transactionUpdateQueueItem) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *transactionUpdateQueueItem) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 8)
	t.fieldMap["signature"] = t.Signature
	t.fieldMap["tx_json"] = t.TxJSON
	t.fieldMap["time"] = t.Time
	t.fieldMap["priority"] = t.Priority
	t.fieldMap["try"] = t.Try
	t.fieldMap["max_try"] = t.MaxTry
	t.fieldMap["retry_time"] = t.RetryTime
	t.fieldMap["reason"] = t.Reason
}

func (t transactionUpdateQueueItem) clone(db *gorm.DB) transactionUpdateQueueItem {
	t.transactionUpdateQueueItemDo.ReplaceDB(db)
	return t
}

type transactionUpdateQueueItemDo struct{ gen.DO }

func (t transactionUpdateQueueItemDo) Debug() *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Debug())
}

func (t transactionUpdateQueueItemDo) WithContext(ctx context.Context) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t transactionUpdateQueueItemDo) ReadDB() *transactionUpdateQueueItemDo {
	return t.Clauses(dbresolver.Read)
}

func (t transactionUpdateQueueItemDo) WriteDB() *transactionUpdateQueueItemDo {
	return t.Clauses(dbresolver.Write)
}

func (t transactionUpdateQueueItemDo) Clauses(conds ...clause.Expression) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t transactionUpdateQueueItemDo) Returning(value interface{}, columns ...string) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t transactionUpdateQueueItemDo) Not(conds ...gen.Condition) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t transactionUpdateQueueItemDo) Or(conds ...gen.Condition) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t transactionUpdateQueueItemDo) Select(conds ...field.Expr) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t transactionUpdateQueueItemDo) Where(conds ...gen.Condition) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t transactionUpdateQueueItemDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *transactionUpdateQueueItemDo {
	return t.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (t transactionUpdateQueueItemDo) Order(conds ...field.Expr) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t transactionUpdateQueueItemDo) Distinct(cols ...field.Expr) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t transactionUpdateQueueItemDo) Omit(cols ...field.Expr) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t transactionUpdateQueueItemDo) Join(table schema.Tabler, on ...field.Expr) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t transactionUpdateQueueItemDo) LeftJoin(table schema.Tabler, on ...field.Expr) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t transactionUpdateQueueItemDo) RightJoin(table schema.Tabler, on ...field.Expr) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t transactionUpdateQueueItemDo) Group(cols ...field.Expr) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t transactionUpdateQueueItemDo) Having(conds ...gen.Condition) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t transactionUpdateQueueItemDo) Limit(limit int) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t transactionUpdateQueueItemDo) Offset(offset int) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t transactionUpdateQueueItemDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t transactionUpdateQueueItemDo) Unscoped() *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Unscoped())
}

func (t transactionUpdateQueueItemDo) Create(values ...*model.TransactionUpdateQueueItem) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t transactionUpdateQueueItemDo) CreateInBatches(values []*model.TransactionUpdateQueueItem, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t transactionUpdateQueueItemDo) Save(values ...*model.TransactionUpdateQueueItem) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t transactionUpdateQueueItemDo) First() (*model.TransactionUpdateQueueItem, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.TransactionUpdateQueueItem), nil
	}
}

func (t transactionUpdateQueueItemDo) Take() (*model.TransactionUpdateQueueItem, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.TransactionUpdateQueueItem), nil
	}
}

func (t transactionUpdateQueueItemDo) Last() (*model.TransactionUpdateQueueItem, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.TransactionUpdateQueueItem), nil
	}
}

func (t transactionUpdateQueueItemDo) Find() ([]*model.TransactionUpdateQueueItem, error) {
	result, err := t.DO.Find()
	return result.([]*model.TransactionUpdateQueueItem), err
}

func (t transactionUpdateQueueItemDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TransactionUpdateQueueItem, err error) {
	buf := make([]*model.TransactionUpdateQueueItem, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t transactionUpdateQueueItemDo) FindInBatches(result *[]*model.TransactionUpdateQueueItem, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t transactionUpdateQueueItemDo) Attrs(attrs ...field.AssignExpr) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t transactionUpdateQueueItemDo) Assign(attrs ...field.AssignExpr) *transactionUpdateQueueItemDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t transactionUpdateQueueItemDo) Joins(fields ...field.RelationField) *transactionUpdateQueueItemDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t transactionUpdateQueueItemDo) Preload(fields ...field.RelationField) *transactionUpdateQueueItemDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t transactionUpdateQueueItemDo) FirstOrInit() (*model.TransactionUpdateQueueItem, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.TransactionUpdateQueueItem), nil
	}
}

func (t transactionUpdateQueueItemDo) FirstOrCreate() (*model.TransactionUpdateQueueItem, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.TransactionUpdateQueueItem), nil
	}
}

func (t transactionUpdateQueueItemDo) FindByPage(offset int, limit int) (result []*model.TransactionUpdateQueueItem, count int64, err error) {
	result, err = t.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = t.Offset(-1).Limit(-1).Count()
	return
}

func (t transactionUpdateQueueItemDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t transactionUpdateQueueItemDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t transactionUpdateQueueItemDo) Delete(models ...*model.TransactionUpdateQueueItem) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *transactionUpdateQueueItemDo) withDO(do gen.Dao) *transactionUpdateQueueItemDo {
	t.DO = *do.(*gen.DO)
	return t
}