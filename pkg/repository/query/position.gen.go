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

func newPosition(db *gorm.DB) position {
	_position := position{}

	_position.positionDo.UseDB(db)
	_position.positionDo.UseModel(&model.Position{})

	tableName := _position.positionDo.TableName()
	_position.ALL = field.NewField(tableName, "*")
	_position.Pubkey = field.NewString(tableName, "pubkey")
	_position.Vault = field.NewString(tableName, "vault")
	_position.Authority = field.NewString(tableName, "authority")
	_position.DepositedTokenAAmount = field.NewUint64(tableName, "deposited_token_a_amount")
	_position.WithdrawnTokenBAmount = field.NewUint64(tableName, "withdrawn_token_b_amount")
	_position.DepositTimestamp = field.NewTime(tableName, "deposit_timestamp")
	_position.DcaPeriodIDBeforeDeposit = field.NewUint64(tableName, "dca_period_id_before_deposit")
	_position.NumberOfSwaps = field.NewUint64(tableName, "number_of_swaps")
	_position.PeriodicDripAmount = field.NewUint64(tableName, "periodic_drip_amount")
	_position.IsClosed = field.NewBool(tableName, "is_closed")

	_position.fillFieldMap()

	return _position
}

type position struct {
	positionDo positionDo

	ALL                      field.Field
	Pubkey                   field.String
	Vault                    field.String
	Authority                field.String
	DepositedTokenAAmount    field.Uint64
	WithdrawnTokenBAmount    field.Uint64
	DepositTimestamp         field.Time
	DcaPeriodIDBeforeDeposit field.Uint64
	NumberOfSwaps            field.Uint64
	PeriodicDripAmount       field.Uint64
	IsClosed                 field.Bool

	fieldMap map[string]field.Expr
}

func (p position) Table(newTableName string) *position {
	p.positionDo.UseTable(newTableName)
	return p.updateTableName(newTableName)
}

func (p position) As(alias string) *position {
	p.positionDo.DO = *(p.positionDo.As(alias).(*gen.DO))
	return p.updateTableName(alias)
}

func (p *position) updateTableName(table string) *position {
	p.ALL = field.NewField(table, "*")
	p.Pubkey = field.NewString(table, "pubkey")
	p.Vault = field.NewString(table, "vault")
	p.Authority = field.NewString(table, "authority")
	p.DepositedTokenAAmount = field.NewUint64(table, "deposited_token_a_amount")
	p.WithdrawnTokenBAmount = field.NewUint64(table, "withdrawn_token_b_amount")
	p.DepositTimestamp = field.NewTime(table, "deposit_timestamp")
	p.DcaPeriodIDBeforeDeposit = field.NewUint64(table, "dca_period_id_before_deposit")
	p.NumberOfSwaps = field.NewUint64(table, "number_of_swaps")
	p.PeriodicDripAmount = field.NewUint64(table, "periodic_drip_amount")
	p.IsClosed = field.NewBool(table, "is_closed")

	p.fillFieldMap()

	return p
}

func (p *position) WithContext(ctx context.Context) *positionDo { return p.positionDo.WithContext(ctx) }

func (p position) TableName() string { return p.positionDo.TableName() }

func (p position) Alias() string { return p.positionDo.Alias() }

func (p *position) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := p.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (p *position) fillFieldMap() {
	p.fieldMap = make(map[string]field.Expr, 10)
	p.fieldMap["pubkey"] = p.Pubkey
	p.fieldMap["vault"] = p.Vault
	p.fieldMap["authority"] = p.Authority
	p.fieldMap["deposited_token_a_amount"] = p.DepositedTokenAAmount
	p.fieldMap["withdrawn_token_b_amount"] = p.WithdrawnTokenBAmount
	p.fieldMap["deposit_timestamp"] = p.DepositTimestamp
	p.fieldMap["dca_period_id_before_deposit"] = p.DcaPeriodIDBeforeDeposit
	p.fieldMap["number_of_swaps"] = p.NumberOfSwaps
	p.fieldMap["periodic_drip_amount"] = p.PeriodicDripAmount
	p.fieldMap["is_closed"] = p.IsClosed
}

func (p position) clone(db *gorm.DB) position {
	p.positionDo.ReplaceDB(db)
	return p
}

type positionDo struct{ gen.DO }

func (p positionDo) Debug() *positionDo {
	return p.withDO(p.DO.Debug())
}

func (p positionDo) WithContext(ctx context.Context) *positionDo {
	return p.withDO(p.DO.WithContext(ctx))
}

func (p positionDo) ReadDB() *positionDo {
	return p.Clauses(dbresolver.Read)
}

func (p positionDo) WriteDB() *positionDo {
	return p.Clauses(dbresolver.Write)
}

func (p positionDo) Clauses(conds ...clause.Expression) *positionDo {
	return p.withDO(p.DO.Clauses(conds...))
}

func (p positionDo) Returning(value interface{}, columns ...string) *positionDo {
	return p.withDO(p.DO.Returning(value, columns...))
}

func (p positionDo) Not(conds ...gen.Condition) *positionDo {
	return p.withDO(p.DO.Not(conds...))
}

func (p positionDo) Or(conds ...gen.Condition) *positionDo {
	return p.withDO(p.DO.Or(conds...))
}

func (p positionDo) Select(conds ...field.Expr) *positionDo {
	return p.withDO(p.DO.Select(conds...))
}

func (p positionDo) Where(conds ...gen.Condition) *positionDo {
	return p.withDO(p.DO.Where(conds...))
}

func (p positionDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *positionDo {
	return p.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (p positionDo) Order(conds ...field.Expr) *positionDo {
	return p.withDO(p.DO.Order(conds...))
}

func (p positionDo) Distinct(cols ...field.Expr) *positionDo {
	return p.withDO(p.DO.Distinct(cols...))
}

func (p positionDo) Omit(cols ...field.Expr) *positionDo {
	return p.withDO(p.DO.Omit(cols...))
}

func (p positionDo) Join(table schema.Tabler, on ...field.Expr) *positionDo {
	return p.withDO(p.DO.Join(table, on...))
}

func (p positionDo) LeftJoin(table schema.Tabler, on ...field.Expr) *positionDo {
	return p.withDO(p.DO.LeftJoin(table, on...))
}

func (p positionDo) RightJoin(table schema.Tabler, on ...field.Expr) *positionDo {
	return p.withDO(p.DO.RightJoin(table, on...))
}

func (p positionDo) Group(cols ...field.Expr) *positionDo {
	return p.withDO(p.DO.Group(cols...))
}

func (p positionDo) Having(conds ...gen.Condition) *positionDo {
	return p.withDO(p.DO.Having(conds...))
}

func (p positionDo) Limit(limit int) *positionDo {
	return p.withDO(p.DO.Limit(limit))
}

func (p positionDo) Offset(offset int) *positionDo {
	return p.withDO(p.DO.Offset(offset))
}

func (p positionDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *positionDo {
	return p.withDO(p.DO.Scopes(funcs...))
}

func (p positionDo) Unscoped() *positionDo {
	return p.withDO(p.DO.Unscoped())
}

func (p positionDo) Create(values ...*model.Position) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Create(values)
}

func (p positionDo) CreateInBatches(values []*model.Position, batchSize int) error {
	return p.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (p positionDo) Save(values ...*model.Position) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Save(values)
}

func (p positionDo) First() (*model.Position, error) {
	if result, err := p.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Position), nil
	}
}

func (p positionDo) Take() (*model.Position, error) {
	if result, err := p.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Position), nil
	}
}

func (p positionDo) Last() (*model.Position, error) {
	if result, err := p.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Position), nil
	}
}

func (p positionDo) Find() ([]*model.Position, error) {
	result, err := p.DO.Find()
	return result.([]*model.Position), err
}

func (p positionDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Position, err error) {
	buf := make([]*model.Position, 0, batchSize)
	err = p.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (p positionDo) FindInBatches(result *[]*model.Position, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return p.DO.FindInBatches(result, batchSize, fc)
}

func (p positionDo) Attrs(attrs ...field.AssignExpr) *positionDo {
	return p.withDO(p.DO.Attrs(attrs...))
}

func (p positionDo) Assign(attrs ...field.AssignExpr) *positionDo {
	return p.withDO(p.DO.Assign(attrs...))
}

func (p positionDo) Joins(fields ...field.RelationField) *positionDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Joins(_f))
	}
	return &p
}

func (p positionDo) Preload(fields ...field.RelationField) *positionDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Preload(_f))
	}
	return &p
}

func (p positionDo) FirstOrInit() (*model.Position, error) {
	if result, err := p.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Position), nil
	}
}

func (p positionDo) FirstOrCreate() (*model.Position, error) {
	if result, err := p.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Position), nil
	}
}

func (p positionDo) FindByPage(offset int, limit int) (result []*model.Position, count int64, err error) {
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

func (p positionDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = p.Count()
	if err != nil {
		return
	}

	err = p.Offset(offset).Limit(limit).Scan(result)
	return
}

func (p positionDo) Scan(result interface{}) (err error) {
	return p.DO.Scan(result)
}

func (p *positionDo) withDO(do gen.Dao) *positionDo {
	p.DO = *do.(*gen.DO)
	return p
}
