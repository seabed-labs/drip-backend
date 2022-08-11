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

func newProtoConfig(db *gorm.DB) protoConfig {
	_protoConfig := protoConfig{}

	_protoConfig.protoConfigDo.UseDB(db)
	_protoConfig.protoConfigDo.UseModel(&model.ProtoConfig{})

	tableName := _protoConfig.protoConfigDo.TableName()
	_protoConfig.ALL = field.NewField(tableName, "*")
	_protoConfig.Pubkey = field.NewString(tableName, "pubkey")
	_protoConfig.Granularity = field.NewUint64(tableName, "granularity")
	_protoConfig.TriggerDcaSpread = field.NewInt16(tableName, "trigger_dca_spread")
	_protoConfig.BaseWithdrawalSpread = field.NewInt16(tableName, "base_withdrawal_spread")
	_protoConfig.Admin = field.NewString(tableName, "admin")

	_protoConfig.fillFieldMap()

	return _protoConfig
}

type protoConfig struct {
	protoConfigDo protoConfigDo

	ALL                  field.Field
	Pubkey               field.String
	Granularity          field.Uint64
	TriggerDcaSpread     field.Int16
	BaseWithdrawalSpread field.Int16
	Admin                field.String

	fieldMap map[string]field.Expr
}

func (p protoConfig) Table(newTableName string) *protoConfig {
	p.protoConfigDo.UseTable(newTableName)
	return p.updateTableName(newTableName)
}

func (p protoConfig) As(alias string) *protoConfig {
	p.protoConfigDo.DO = *(p.protoConfigDo.As(alias).(*gen.DO))
	return p.updateTableName(alias)
}

func (p *protoConfig) updateTableName(table string) *protoConfig {
	p.ALL = field.NewField(table, "*")
	p.Pubkey = field.NewString(table, "pubkey")
	p.Granularity = field.NewUint64(table, "granularity")
	p.TriggerDcaSpread = field.NewInt16(table, "trigger_dca_spread")
	p.BaseWithdrawalSpread = field.NewInt16(table, "base_withdrawal_spread")
	p.Admin = field.NewString(table, "admin")

	p.fillFieldMap()

	return p
}

func (p *protoConfig) WithContext(ctx context.Context) *protoConfigDo {
	return p.protoConfigDo.WithContext(ctx)
}

func (p protoConfig) TableName() string { return p.protoConfigDo.TableName() }

func (p protoConfig) Alias() string { return p.protoConfigDo.Alias() }

func (p *protoConfig) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := p.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (p *protoConfig) fillFieldMap() {
	p.fieldMap = make(map[string]field.Expr, 5)
	p.fieldMap["pubkey"] = p.Pubkey
	p.fieldMap["granularity"] = p.Granularity
	p.fieldMap["trigger_dca_spread"] = p.TriggerDcaSpread
	p.fieldMap["base_withdrawal_spread"] = p.BaseWithdrawalSpread
	p.fieldMap["admin"] = p.Admin
}

func (p protoConfig) clone(db *gorm.DB) protoConfig {
	p.protoConfigDo.ReplaceDB(db)
	return p
}

type protoConfigDo struct{ gen.DO }

func (p protoConfigDo) Debug() *protoConfigDo {
	return p.withDO(p.DO.Debug())
}

func (p protoConfigDo) WithContext(ctx context.Context) *protoConfigDo {
	return p.withDO(p.DO.WithContext(ctx))
}

func (p protoConfigDo) ReadDB() *protoConfigDo {
	return p.Clauses(dbresolver.Read)
}

func (p protoConfigDo) WriteDB() *protoConfigDo {
	return p.Clauses(dbresolver.Write)
}

func (p protoConfigDo) Clauses(conds ...clause.Expression) *protoConfigDo {
	return p.withDO(p.DO.Clauses(conds...))
}

func (p protoConfigDo) Returning(value interface{}, columns ...string) *protoConfigDo {
	return p.withDO(p.DO.Returning(value, columns...))
}

func (p protoConfigDo) Not(conds ...gen.Condition) *protoConfigDo {
	return p.withDO(p.DO.Not(conds...))
}

func (p protoConfigDo) Or(conds ...gen.Condition) *protoConfigDo {
	return p.withDO(p.DO.Or(conds...))
}

func (p protoConfigDo) Select(conds ...field.Expr) *protoConfigDo {
	return p.withDO(p.DO.Select(conds...))
}

func (p protoConfigDo) Where(conds ...gen.Condition) *protoConfigDo {
	return p.withDO(p.DO.Where(conds...))
}

func (p protoConfigDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *protoConfigDo {
	return p.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (p protoConfigDo) Order(conds ...field.Expr) *protoConfigDo {
	return p.withDO(p.DO.Order(conds...))
}

func (p protoConfigDo) Distinct(cols ...field.Expr) *protoConfigDo {
	return p.withDO(p.DO.Distinct(cols...))
}

func (p protoConfigDo) Omit(cols ...field.Expr) *protoConfigDo {
	return p.withDO(p.DO.Omit(cols...))
}

func (p protoConfigDo) Join(table schema.Tabler, on ...field.Expr) *protoConfigDo {
	return p.withDO(p.DO.Join(table, on...))
}

func (p protoConfigDo) LeftJoin(table schema.Tabler, on ...field.Expr) *protoConfigDo {
	return p.withDO(p.DO.LeftJoin(table, on...))
}

func (p protoConfigDo) RightJoin(table schema.Tabler, on ...field.Expr) *protoConfigDo {
	return p.withDO(p.DO.RightJoin(table, on...))
}

func (p protoConfigDo) Group(cols ...field.Expr) *protoConfigDo {
	return p.withDO(p.DO.Group(cols...))
}

func (p protoConfigDo) Having(conds ...gen.Condition) *protoConfigDo {
	return p.withDO(p.DO.Having(conds...))
}

func (p protoConfigDo) Limit(limit int) *protoConfigDo {
	return p.withDO(p.DO.Limit(limit))
}

func (p protoConfigDo) Offset(offset int) *protoConfigDo {
	return p.withDO(p.DO.Offset(offset))
}

func (p protoConfigDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *protoConfigDo {
	return p.withDO(p.DO.Scopes(funcs...))
}

func (p protoConfigDo) Unscoped() *protoConfigDo {
	return p.withDO(p.DO.Unscoped())
}

func (p protoConfigDo) Create(values ...*model.ProtoConfig) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Create(values)
}

func (p protoConfigDo) CreateInBatches(values []*model.ProtoConfig, batchSize int) error {
	return p.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (p protoConfigDo) Save(values ...*model.ProtoConfig) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Save(values)
}

func (p protoConfigDo) First() (*model.ProtoConfig, error) {
	if result, err := p.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.ProtoConfig), nil
	}
}

func (p protoConfigDo) Take() (*model.ProtoConfig, error) {
	if result, err := p.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.ProtoConfig), nil
	}
}

func (p protoConfigDo) Last() (*model.ProtoConfig, error) {
	if result, err := p.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.ProtoConfig), nil
	}
}

func (p protoConfigDo) Find() ([]*model.ProtoConfig, error) {
	result, err := p.DO.Find()
	return result.([]*model.ProtoConfig), err
}

func (p protoConfigDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ProtoConfig, err error) {
	buf := make([]*model.ProtoConfig, 0, batchSize)
	err = p.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (p protoConfigDo) FindInBatches(result *[]*model.ProtoConfig, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return p.DO.FindInBatches(result, batchSize, fc)
}

func (p protoConfigDo) Attrs(attrs ...field.AssignExpr) *protoConfigDo {
	return p.withDO(p.DO.Attrs(attrs...))
}

func (p protoConfigDo) Assign(attrs ...field.AssignExpr) *protoConfigDo {
	return p.withDO(p.DO.Assign(attrs...))
}

func (p protoConfigDo) Joins(fields ...field.RelationField) *protoConfigDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Joins(_f))
	}
	return &p
}

func (p protoConfigDo) Preload(fields ...field.RelationField) *protoConfigDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Preload(_f))
	}
	return &p
}

func (p protoConfigDo) FirstOrInit() (*model.ProtoConfig, error) {
	if result, err := p.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.ProtoConfig), nil
	}
}

func (p protoConfigDo) FirstOrCreate() (*model.ProtoConfig, error) {
	if result, err := p.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.ProtoConfig), nil
	}
}

func (p protoConfigDo) FindByPage(offset int, limit int) (result []*model.ProtoConfig, count int64, err error) {
	result, err = p.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = p.Offset(-1).Limit(-1).Count()
	return
}

func (p protoConfigDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = p.Count()
	if err != nil {
		return
	}

	err = p.Offset(offset).Limit(limit).Scan(result)
	return
}

func (p protoConfigDo) Scan(result interface{}) (err error) {
	return p.DO.Scan(result)
}

func (p *protoConfigDo) withDO(do gen.Dao) *protoConfigDo {
	p.DO = *do.(*gen.DO)
	return p
}
