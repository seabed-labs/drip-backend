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

func newOrcaWhirlpool(db *gorm.DB) orcaWhirlpool {
	_orcaWhirlpool := orcaWhirlpool{}

	_orcaWhirlpool.orcaWhirlpoolDo.UseDB(db)
	_orcaWhirlpool.orcaWhirlpoolDo.UseModel(&model.OrcaWhirlpool{})

	tableName := _orcaWhirlpool.orcaWhirlpoolDo.TableName()
	_orcaWhirlpool.ALL = field.NewField(tableName, "*")
	_orcaWhirlpool.Pubkey = field.NewString(tableName, "pubkey")
	_orcaWhirlpool.WhirlpoolsConfig = field.NewString(tableName, "whirlpools_config")
	_orcaWhirlpool.TokenMintA = field.NewString(tableName, "token_mint_a")
	_orcaWhirlpool.TokenVaultA = field.NewString(tableName, "token_vault_a")
	_orcaWhirlpool.TokenMintB = field.NewString(tableName, "token_mint_b")
	_orcaWhirlpool.TokenVaultB = field.NewString(tableName, "token_vault_b")
	_orcaWhirlpool.Oracle = field.NewString(tableName, "oracle")
	_orcaWhirlpool.TickSpacing = field.NewInt32(tableName, "tick_spacing")
	_orcaWhirlpool.FeeRate = field.NewInt32(tableName, "fee_rate")
	_orcaWhirlpool.ProtocolFeeRate = field.NewInt32(tableName, "protocol_fee_rate")
	_orcaWhirlpool.TickCurrentIndex = field.NewInt32(tableName, "tick_current_index")
	_orcaWhirlpool.ProtocolFeeOwedA = field.NewUint64(tableName, "protocol_fee_owed_a")
	_orcaWhirlpool.ProtocolFeeOwedB = field.NewUint64(tableName, "protocol_fee_owed_b")
	_orcaWhirlpool.RewardLastUpdatedTimestamp = field.NewUint64(tableName, "reward_last_updated_timestamp")
	_orcaWhirlpool.Liquidity = field.NewUint64(tableName, "liquidity")
	_orcaWhirlpool.SqrtPrice = field.NewUint64(tableName, "sqrt_price")
	_orcaWhirlpool.FeeGrowthGlobalA = field.NewUint64(tableName, "fee_growth_global_a")
	_orcaWhirlpool.FeeGrowthGlobalB = field.NewUint64(tableName, "fee_growth_global_b")
	_orcaWhirlpool.TokenPairID = field.NewString(tableName, "token_pair_id")
	_orcaWhirlpool.ID = field.NewString(tableName, "id")

	_orcaWhirlpool.fillFieldMap()

	return _orcaWhirlpool
}

type orcaWhirlpool struct {
	orcaWhirlpoolDo orcaWhirlpoolDo

	ALL                        field.Field
	Pubkey                     field.String
	WhirlpoolsConfig           field.String
	TokenMintA                 field.String
	TokenVaultA                field.String
	TokenMintB                 field.String
	TokenVaultB                field.String
	Oracle                     field.String
	TickSpacing                field.Int32
	FeeRate                    field.Int32
	ProtocolFeeRate            field.Int32
	TickCurrentIndex           field.Int32
	ProtocolFeeOwedA           field.Uint64
	ProtocolFeeOwedB           field.Uint64
	RewardLastUpdatedTimestamp field.Uint64
	Liquidity                  field.Uint64
	SqrtPrice                  field.Uint64
	FeeGrowthGlobalA           field.Uint64
	FeeGrowthGlobalB           field.Uint64
	TokenPairID                field.String
	ID                         field.String

	fieldMap map[string]field.Expr
}

func (o orcaWhirlpool) Table(newTableName string) *orcaWhirlpool {
	o.orcaWhirlpoolDo.UseTable(newTableName)
	return o.updateTableName(newTableName)
}

func (o orcaWhirlpool) As(alias string) *orcaWhirlpool {
	o.orcaWhirlpoolDo.DO = *(o.orcaWhirlpoolDo.As(alias).(*gen.DO))
	return o.updateTableName(alias)
}

func (o *orcaWhirlpool) updateTableName(table string) *orcaWhirlpool {
	o.ALL = field.NewField(table, "*")
	o.Pubkey = field.NewString(table, "pubkey")
	o.WhirlpoolsConfig = field.NewString(table, "whirlpools_config")
	o.TokenMintA = field.NewString(table, "token_mint_a")
	o.TokenVaultA = field.NewString(table, "token_vault_a")
	o.TokenMintB = field.NewString(table, "token_mint_b")
	o.TokenVaultB = field.NewString(table, "token_vault_b")
	o.Oracle = field.NewString(table, "oracle")
	o.TickSpacing = field.NewInt32(table, "tick_spacing")
	o.FeeRate = field.NewInt32(table, "fee_rate")
	o.ProtocolFeeRate = field.NewInt32(table, "protocol_fee_rate")
	o.TickCurrentIndex = field.NewInt32(table, "tick_current_index")
	o.ProtocolFeeOwedA = field.NewUint64(table, "protocol_fee_owed_a")
	o.ProtocolFeeOwedB = field.NewUint64(table, "protocol_fee_owed_b")
	o.RewardLastUpdatedTimestamp = field.NewUint64(table, "reward_last_updated_timestamp")
	o.Liquidity = field.NewUint64(table, "liquidity")
	o.SqrtPrice = field.NewUint64(table, "sqrt_price")
	o.FeeGrowthGlobalA = field.NewUint64(table, "fee_growth_global_a")
	o.FeeGrowthGlobalB = field.NewUint64(table, "fee_growth_global_b")
	o.TokenPairID = field.NewString(table, "token_pair_id")
	o.ID = field.NewString(table, "id")

	o.fillFieldMap()

	return o
}

func (o *orcaWhirlpool) WithContext(ctx context.Context) *orcaWhirlpoolDo {
	return o.orcaWhirlpoolDo.WithContext(ctx)
}

func (o orcaWhirlpool) TableName() string { return o.orcaWhirlpoolDo.TableName() }

func (o orcaWhirlpool) Alias() string { return o.orcaWhirlpoolDo.Alias() }

func (o *orcaWhirlpool) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := o.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (o *orcaWhirlpool) fillFieldMap() {
	o.fieldMap = make(map[string]field.Expr, 20)
	o.fieldMap["pubkey"] = o.Pubkey
	o.fieldMap["whirlpools_config"] = o.WhirlpoolsConfig
	o.fieldMap["token_mint_a"] = o.TokenMintA
	o.fieldMap["token_vault_a"] = o.TokenVaultA
	o.fieldMap["token_mint_b"] = o.TokenMintB
	o.fieldMap["token_vault_b"] = o.TokenVaultB
	o.fieldMap["oracle"] = o.Oracle
	o.fieldMap["tick_spacing"] = o.TickSpacing
	o.fieldMap["fee_rate"] = o.FeeRate
	o.fieldMap["protocol_fee_rate"] = o.ProtocolFeeRate
	o.fieldMap["tick_current_index"] = o.TickCurrentIndex
	o.fieldMap["protocol_fee_owed_a"] = o.ProtocolFeeOwedA
	o.fieldMap["protocol_fee_owed_b"] = o.ProtocolFeeOwedB
	o.fieldMap["reward_last_updated_timestamp"] = o.RewardLastUpdatedTimestamp
	o.fieldMap["liquidity"] = o.Liquidity
	o.fieldMap["sqrt_price"] = o.SqrtPrice
	o.fieldMap["fee_growth_global_a"] = o.FeeGrowthGlobalA
	o.fieldMap["fee_growth_global_b"] = o.FeeGrowthGlobalB
	o.fieldMap["token_pair_id"] = o.TokenPairID
	o.fieldMap["id"] = o.ID
}

func (o orcaWhirlpool) clone(db *gorm.DB) orcaWhirlpool {
	o.orcaWhirlpoolDo.ReplaceDB(db)
	return o
}

type orcaWhirlpoolDo struct{ gen.DO }

func (o orcaWhirlpoolDo) Debug() *orcaWhirlpoolDo {
	return o.withDO(o.DO.Debug())
}

func (o orcaWhirlpoolDo) WithContext(ctx context.Context) *orcaWhirlpoolDo {
	return o.withDO(o.DO.WithContext(ctx))
}

func (o orcaWhirlpoolDo) ReadDB() *orcaWhirlpoolDo {
	return o.Clauses(dbresolver.Read)
}

func (o orcaWhirlpoolDo) WriteDB() *orcaWhirlpoolDo {
	return o.Clauses(dbresolver.Write)
}

func (o orcaWhirlpoolDo) Clauses(conds ...clause.Expression) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Clauses(conds...))
}

func (o orcaWhirlpoolDo) Returning(value interface{}, columns ...string) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Returning(value, columns...))
}

func (o orcaWhirlpoolDo) Not(conds ...gen.Condition) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Not(conds...))
}

func (o orcaWhirlpoolDo) Or(conds ...gen.Condition) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Or(conds...))
}

func (o orcaWhirlpoolDo) Select(conds ...field.Expr) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Select(conds...))
}

func (o orcaWhirlpoolDo) Where(conds ...gen.Condition) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Where(conds...))
}

func (o orcaWhirlpoolDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *orcaWhirlpoolDo {
	return o.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (o orcaWhirlpoolDo) Order(conds ...field.Expr) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Order(conds...))
}

func (o orcaWhirlpoolDo) Distinct(cols ...field.Expr) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Distinct(cols...))
}

func (o orcaWhirlpoolDo) Omit(cols ...field.Expr) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Omit(cols...))
}

func (o orcaWhirlpoolDo) Join(table schema.Tabler, on ...field.Expr) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Join(table, on...))
}

func (o orcaWhirlpoolDo) LeftJoin(table schema.Tabler, on ...field.Expr) *orcaWhirlpoolDo {
	return o.withDO(o.DO.LeftJoin(table, on...))
}

func (o orcaWhirlpoolDo) RightJoin(table schema.Tabler, on ...field.Expr) *orcaWhirlpoolDo {
	return o.withDO(o.DO.RightJoin(table, on...))
}

func (o orcaWhirlpoolDo) Group(cols ...field.Expr) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Group(cols...))
}

func (o orcaWhirlpoolDo) Having(conds ...gen.Condition) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Having(conds...))
}

func (o orcaWhirlpoolDo) Limit(limit int) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Limit(limit))
}

func (o orcaWhirlpoolDo) Offset(offset int) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Offset(offset))
}

func (o orcaWhirlpoolDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Scopes(funcs...))
}

func (o orcaWhirlpoolDo) Unscoped() *orcaWhirlpoolDo {
	return o.withDO(o.DO.Unscoped())
}

func (o orcaWhirlpoolDo) Create(values ...*model.OrcaWhirlpool) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Create(values)
}

func (o orcaWhirlpoolDo) CreateInBatches(values []*model.OrcaWhirlpool, batchSize int) error {
	return o.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (o orcaWhirlpoolDo) Save(values ...*model.OrcaWhirlpool) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Save(values)
}

func (o orcaWhirlpoolDo) First() (*model.OrcaWhirlpool, error) {
	if result, err := o.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrcaWhirlpool), nil
	}
}

func (o orcaWhirlpoolDo) Take() (*model.OrcaWhirlpool, error) {
	if result, err := o.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrcaWhirlpool), nil
	}
}

func (o orcaWhirlpoolDo) Last() (*model.OrcaWhirlpool, error) {
	if result, err := o.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrcaWhirlpool), nil
	}
}

func (o orcaWhirlpoolDo) Find() ([]*model.OrcaWhirlpool, error) {
	result, err := o.DO.Find()
	return result.([]*model.OrcaWhirlpool), err
}

func (o orcaWhirlpoolDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.OrcaWhirlpool, err error) {
	buf := make([]*model.OrcaWhirlpool, 0, batchSize)
	err = o.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (o orcaWhirlpoolDo) FindInBatches(result *[]*model.OrcaWhirlpool, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return o.DO.FindInBatches(result, batchSize, fc)
}

func (o orcaWhirlpoolDo) Attrs(attrs ...field.AssignExpr) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Attrs(attrs...))
}

func (o orcaWhirlpoolDo) Assign(attrs ...field.AssignExpr) *orcaWhirlpoolDo {
	return o.withDO(o.DO.Assign(attrs...))
}

func (o orcaWhirlpoolDo) Joins(fields ...field.RelationField) *orcaWhirlpoolDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Joins(_f))
	}
	return &o
}

func (o orcaWhirlpoolDo) Preload(fields ...field.RelationField) *orcaWhirlpoolDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Preload(_f))
	}
	return &o
}

func (o orcaWhirlpoolDo) FirstOrInit() (*model.OrcaWhirlpool, error) {
	if result, err := o.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrcaWhirlpool), nil
	}
}

func (o orcaWhirlpoolDo) FirstOrCreate() (*model.OrcaWhirlpool, error) {
	if result, err := o.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.OrcaWhirlpool), nil
	}
}

func (o orcaWhirlpoolDo) FindByPage(offset int, limit int) (result []*model.OrcaWhirlpool, count int64, err error) {
	result, err = o.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = o.Offset(-1).Limit(-1).Count()
	return
}

func (o orcaWhirlpoolDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = o.Count()
	if err != nil {
		return
	}

	err = o.Offset(offset).Limit(limit).Scan(result)
	return
}

func (o orcaWhirlpoolDo) Scan(result interface{}) (err error) {
	return o.DO.Scan(result)
}

func (o orcaWhirlpoolDo) Delete(models ...*model.OrcaWhirlpool) (result gen.ResultInfo, err error) {
	return o.DO.Delete(models)
}

func (o *orcaWhirlpoolDo) withDO(do gen.Dao) *orcaWhirlpoolDo {
	o.DO = *do.(*gen.DO)
	return o
}
