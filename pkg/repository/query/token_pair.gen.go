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

	"github.com/dcaf-labs/drip/pkg/repository/model"
)

func newTokenPair(db *gorm.DB) tokenPair {
	_tokenPair := tokenPair{}

	_tokenPair.tokenPairDo.UseDB(db)
	_tokenPair.tokenPairDo.UseModel(&model.TokenPair{})

	tableName := _tokenPair.tokenPairDo.TableName()
	_tokenPair.ALL = field.NewField(tableName, "*")
	_tokenPair.ID = field.NewString(tableName, "id")
	_tokenPair.TokenA = field.NewString(tableName, "token_a")
	_tokenPair.TokenB = field.NewString(tableName, "token_b")

	_tokenPair.fillFieldMap()

	return _tokenPair
}

type tokenPair struct {
	tokenPairDo tokenPairDo

	ALL    field.Field
	ID     field.String
	TokenA field.String
	TokenB field.String

	fieldMap map[string]field.Expr
}

func (t tokenPair) Table(newTableName string) *tokenPair {
	t.tokenPairDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t tokenPair) As(alias string) *tokenPair {
	t.tokenPairDo.DO = *(t.tokenPairDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *tokenPair) updateTableName(table string) *tokenPair {
	t.ALL = field.NewField(table, "*")
	t.ID = field.NewString(table, "id")
	t.TokenA = field.NewString(table, "token_a")
	t.TokenB = field.NewString(table, "token_b")

	t.fillFieldMap()

	return t
}

func (t *tokenPair) WithContext(ctx context.Context) *tokenPairDo {
	return t.tokenPairDo.WithContext(ctx)
}

func (t tokenPair) TableName() string { return t.tokenPairDo.TableName() }

func (t tokenPair) Alias() string { return t.tokenPairDo.Alias() }

func (t *tokenPair) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *tokenPair) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 3)
	t.fieldMap["id"] = t.ID
	t.fieldMap["token_a"] = t.TokenA
	t.fieldMap["token_b"] = t.TokenB
}

func (t tokenPair) clone(db *gorm.DB) tokenPair {
	t.tokenPairDo.ReplaceDB(db)
	return t
}

type tokenPairDo struct{ gen.DO }

func (t tokenPairDo) Debug() *tokenPairDo {
	return t.withDO(t.DO.Debug())
}

func (t tokenPairDo) WithContext(ctx context.Context) *tokenPairDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t tokenPairDo) ReadDB() *tokenPairDo {
	return t.Clauses(dbresolver.Read)
}

func (t tokenPairDo) WriteDB() *tokenPairDo {
	return t.Clauses(dbresolver.Write)
}

func (t tokenPairDo) Clauses(conds ...clause.Expression) *tokenPairDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t tokenPairDo) Returning(value interface{}, columns ...string) *tokenPairDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t tokenPairDo) Not(conds ...gen.Condition) *tokenPairDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t tokenPairDo) Or(conds ...gen.Condition) *tokenPairDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t tokenPairDo) Select(conds ...field.Expr) *tokenPairDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t tokenPairDo) Where(conds ...gen.Condition) *tokenPairDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t tokenPairDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *tokenPairDo {
	return t.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (t tokenPairDo) Order(conds ...field.Expr) *tokenPairDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t tokenPairDo) Distinct(cols ...field.Expr) *tokenPairDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t tokenPairDo) Omit(cols ...field.Expr) *tokenPairDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t tokenPairDo) Join(table schema.Tabler, on ...field.Expr) *tokenPairDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t tokenPairDo) LeftJoin(table schema.Tabler, on ...field.Expr) *tokenPairDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t tokenPairDo) RightJoin(table schema.Tabler, on ...field.Expr) *tokenPairDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t tokenPairDo) Group(cols ...field.Expr) *tokenPairDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t tokenPairDo) Having(conds ...gen.Condition) *tokenPairDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t tokenPairDo) Limit(limit int) *tokenPairDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t tokenPairDo) Offset(offset int) *tokenPairDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t tokenPairDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *tokenPairDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t tokenPairDo) Unscoped() *tokenPairDo {
	return t.withDO(t.DO.Unscoped())
}

func (t tokenPairDo) Create(values ...*model.TokenPair) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t tokenPairDo) CreateInBatches(values []*model.TokenPair, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t tokenPairDo) Save(values ...*model.TokenPair) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t tokenPairDo) First() (*model.TokenPair, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.TokenPair), nil
	}
}

func (t tokenPairDo) Take() (*model.TokenPair, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.TokenPair), nil
	}
}

func (t tokenPairDo) Last() (*model.TokenPair, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.TokenPair), nil
	}
}

func (t tokenPairDo) Find() ([]*model.TokenPair, error) {
	result, err := t.DO.Find()
	return result.([]*model.TokenPair), err
}

func (t tokenPairDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TokenPair, err error) {
	buf := make([]*model.TokenPair, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t tokenPairDo) FindInBatches(result *[]*model.TokenPair, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t tokenPairDo) Attrs(attrs ...field.AssignExpr) *tokenPairDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t tokenPairDo) Assign(attrs ...field.AssignExpr) *tokenPairDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t tokenPairDo) Joins(fields ...field.RelationField) *tokenPairDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t tokenPairDo) Preload(fields ...field.RelationField) *tokenPairDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t tokenPairDo) FirstOrInit() (*model.TokenPair, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.TokenPair), nil
	}
}

func (t tokenPairDo) FirstOrCreate() (*model.TokenPair, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.TokenPair), nil
	}
}

func (t tokenPairDo) FindByPage(offset int, limit int) (result []*model.TokenPair, count int64, err error) {
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

func (t tokenPairDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t tokenPairDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t tokenPairDo) Delete(models ...*model.TokenPair) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *tokenPairDo) withDO(do gen.Dao) *tokenPairDo {
	t.DO = *do.(*gen.DO)
	return t
}
