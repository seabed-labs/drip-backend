// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/repository/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"
)

func newSourceReference(db *gorm.DB) sourceReference {
	_sourceReference := sourceReference{}

	_sourceReference.sourceReferenceDo.UseDB(db)
	_sourceReference.sourceReferenceDo.UseModel(&model.SourceReference{})

	tableName := _sourceReference.sourceReferenceDo.TableName()
	_sourceReference.ALL = field.NewField(tableName, "*")
	_sourceReference.Value = field.NewString(tableName, "value")

	_sourceReference.fillFieldMap()

	return _sourceReference
}

type sourceReference struct {
	sourceReferenceDo sourceReferenceDo

	ALL   field.Field
	Value field.String

	fieldMap map[string]field.Expr
}

func (s sourceReference) Table(newTableName string) *sourceReference {
	s.sourceReferenceDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s sourceReference) As(alias string) *sourceReference {
	s.sourceReferenceDo.DO = *(s.sourceReferenceDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *sourceReference) updateTableName(table string) *sourceReference {
	s.ALL = field.NewField(table, "*")
	s.Value = field.NewString(table, "value")

	s.fillFieldMap()

	return s
}

func (s *sourceReference) WithContext(ctx context.Context) *sourceReferenceDo {
	return s.sourceReferenceDo.WithContext(ctx)
}

func (s sourceReference) TableName() string { return s.sourceReferenceDo.TableName() }

func (s sourceReference) Alias() string { return s.sourceReferenceDo.Alias() }

func (s *sourceReference) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *sourceReference) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 1)
	s.fieldMap["value"] = s.Value
}

func (s sourceReference) clone(db *gorm.DB) sourceReference {
	s.sourceReferenceDo.ReplaceDB(db)
	return s
}

type sourceReferenceDo struct{ gen.DO }

func (s sourceReferenceDo) Debug() *sourceReferenceDo {
	return s.withDO(s.DO.Debug())
}

func (s sourceReferenceDo) WithContext(ctx context.Context) *sourceReferenceDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s sourceReferenceDo) ReadDB() *sourceReferenceDo {
	return s.Clauses(dbresolver.Read)
}

func (s sourceReferenceDo) WriteDB() *sourceReferenceDo {
	return s.Clauses(dbresolver.Write)
}

func (s sourceReferenceDo) Clauses(conds ...clause.Expression) *sourceReferenceDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s sourceReferenceDo) Returning(value interface{}, columns ...string) *sourceReferenceDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s sourceReferenceDo) Not(conds ...gen.Condition) *sourceReferenceDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s sourceReferenceDo) Or(conds ...gen.Condition) *sourceReferenceDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s sourceReferenceDo) Select(conds ...field.Expr) *sourceReferenceDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s sourceReferenceDo) Where(conds ...gen.Condition) *sourceReferenceDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s sourceReferenceDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *sourceReferenceDo {
	return s.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (s sourceReferenceDo) Order(conds ...field.Expr) *sourceReferenceDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s sourceReferenceDo) Distinct(cols ...field.Expr) *sourceReferenceDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s sourceReferenceDo) Omit(cols ...field.Expr) *sourceReferenceDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s sourceReferenceDo) Join(table schema.Tabler, on ...field.Expr) *sourceReferenceDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s sourceReferenceDo) LeftJoin(table schema.Tabler, on ...field.Expr) *sourceReferenceDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s sourceReferenceDo) RightJoin(table schema.Tabler, on ...field.Expr) *sourceReferenceDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s sourceReferenceDo) Group(cols ...field.Expr) *sourceReferenceDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s sourceReferenceDo) Having(conds ...gen.Condition) *sourceReferenceDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s sourceReferenceDo) Limit(limit int) *sourceReferenceDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s sourceReferenceDo) Offset(offset int) *sourceReferenceDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s sourceReferenceDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *sourceReferenceDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s sourceReferenceDo) Unscoped() *sourceReferenceDo {
	return s.withDO(s.DO.Unscoped())
}

func (s sourceReferenceDo) Create(values ...*model.SourceReference) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s sourceReferenceDo) CreateInBatches(values []*model.SourceReference, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s sourceReferenceDo) Save(values ...*model.SourceReference) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s sourceReferenceDo) First() (*model.SourceReference, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.SourceReference), nil
	}
}

func (s sourceReferenceDo) Take() (*model.SourceReference, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.SourceReference), nil
	}
}

func (s sourceReferenceDo) Last() (*model.SourceReference, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.SourceReference), nil
	}
}

func (s sourceReferenceDo) Find() ([]*model.SourceReference, error) {
	result, err := s.DO.Find()
	return result.([]*model.SourceReference), err
}

func (s sourceReferenceDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.SourceReference, err error) {
	buf := make([]*model.SourceReference, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s sourceReferenceDo) FindInBatches(result *[]*model.SourceReference, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s sourceReferenceDo) Attrs(attrs ...field.AssignExpr) *sourceReferenceDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s sourceReferenceDo) Assign(attrs ...field.AssignExpr) *sourceReferenceDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s sourceReferenceDo) Joins(fields ...field.RelationField) *sourceReferenceDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s sourceReferenceDo) Preload(fields ...field.RelationField) *sourceReferenceDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s sourceReferenceDo) FirstOrInit() (*model.SourceReference, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.SourceReference), nil
	}
}

func (s sourceReferenceDo) FirstOrCreate() (*model.SourceReference, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.SourceReference), nil
	}
}

func (s sourceReferenceDo) FindByPage(offset int, limit int) (result []*model.SourceReference, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s sourceReferenceDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s sourceReferenceDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s *sourceReferenceDo) withDO(do gen.Dao) *sourceReferenceDo {
	s.DO = *do.(*gen.DO)
	return s
}
