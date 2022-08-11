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

func newVault(db *gorm.DB) vault {
	_vault := vault{}

	_vault.vaultDo.UseDB(db)
	_vault.vaultDo.UseModel(&model.Vault{})

	tableName := _vault.vaultDo.TableName()
	_vault.ALL = field.NewField(tableName, "*")
	_vault.Pubkey = field.NewString(tableName, "pubkey")
	_vault.ProtoConfig = field.NewString(tableName, "proto_config")
	_vault.TokenAAccount = field.NewString(tableName, "token_a_account")
	_vault.TokenBAccount = field.NewString(tableName, "token_b_account")
	_vault.TreasuryTokenBAccount = field.NewString(tableName, "treasury_token_b_account")
	_vault.LastDcaPeriod = field.NewUint64(tableName, "last_dca_period")
	_vault.DripAmount = field.NewUint64(tableName, "drip_amount")
	_vault.DcaActivationTimestamp = field.NewTime(tableName, "dca_activation_timestamp")
	_vault.Enabled = field.NewBool(tableName, "enabled")
	_vault.TokenPairID = field.NewString(tableName, "token_pair_id")
	_vault.MaxSlippageBps = field.NewInt32(tableName, "max_slippage_bps")

	_vault.fillFieldMap()

	return _vault
}

type vault struct {
	vaultDo vaultDo

	ALL                    field.Field
	Pubkey                 field.String
	ProtoConfig            field.String
	TokenAAccount          field.String
	TokenBAccount          field.String
	TreasuryTokenBAccount  field.String
	LastDcaPeriod          field.Uint64
	DripAmount             field.Uint64
	DcaActivationTimestamp field.Time
	Enabled                field.Bool
	TokenPairID            field.String
	MaxSlippageBps         field.Int32

	fieldMap map[string]field.Expr
}

func (v vault) Table(newTableName string) *vault {
	v.vaultDo.UseTable(newTableName)
	return v.updateTableName(newTableName)
}

func (v vault) As(alias string) *vault {
	v.vaultDo.DO = *(v.vaultDo.As(alias).(*gen.DO))
	return v.updateTableName(alias)
}

func (v *vault) updateTableName(table string) *vault {
	v.ALL = field.NewField(table, "*")
	v.Pubkey = field.NewString(table, "pubkey")
	v.ProtoConfig = field.NewString(table, "proto_config")
	v.TokenAAccount = field.NewString(table, "token_a_account")
	v.TokenBAccount = field.NewString(table, "token_b_account")
	v.TreasuryTokenBAccount = field.NewString(table, "treasury_token_b_account")
	v.LastDcaPeriod = field.NewUint64(table, "last_dca_period")
	v.DripAmount = field.NewUint64(table, "drip_amount")
	v.DcaActivationTimestamp = field.NewTime(table, "dca_activation_timestamp")
	v.Enabled = field.NewBool(table, "enabled")
	v.TokenPairID = field.NewString(table, "token_pair_id")
	v.MaxSlippageBps = field.NewInt32(table, "max_slippage_bps")

	v.fillFieldMap()

	return v
}

func (v *vault) WithContext(ctx context.Context) *vaultDo { return v.vaultDo.WithContext(ctx) }

func (v vault) TableName() string { return v.vaultDo.TableName() }

func (v vault) Alias() string { return v.vaultDo.Alias() }

func (v *vault) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := v.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (v *vault) fillFieldMap() {
	v.fieldMap = make(map[string]field.Expr, 11)
	v.fieldMap["pubkey"] = v.Pubkey
	v.fieldMap["proto_config"] = v.ProtoConfig
	v.fieldMap["token_a_account"] = v.TokenAAccount
	v.fieldMap["token_b_account"] = v.TokenBAccount
	v.fieldMap["treasury_token_b_account"] = v.TreasuryTokenBAccount
	v.fieldMap["last_dca_period"] = v.LastDcaPeriod
	v.fieldMap["drip_amount"] = v.DripAmount
	v.fieldMap["dca_activation_timestamp"] = v.DcaActivationTimestamp
	v.fieldMap["enabled"] = v.Enabled
	v.fieldMap["token_pair_id"] = v.TokenPairID
	v.fieldMap["max_slippage_bps"] = v.MaxSlippageBps
}

func (v vault) clone(db *gorm.DB) vault {
	v.vaultDo.ReplaceDB(db)
	return v
}

type vaultDo struct{ gen.DO }

func (v vaultDo) Debug() *vaultDo {
	return v.withDO(v.DO.Debug())
}

func (v vaultDo) WithContext(ctx context.Context) *vaultDo {
	return v.withDO(v.DO.WithContext(ctx))
}

func (v vaultDo) ReadDB() *vaultDo {
	return v.Clauses(dbresolver.Read)
}

func (v vaultDo) WriteDB() *vaultDo {
	return v.Clauses(dbresolver.Write)
}

func (v vaultDo) Clauses(conds ...clause.Expression) *vaultDo {
	return v.withDO(v.DO.Clauses(conds...))
}

func (v vaultDo) Returning(value interface{}, columns ...string) *vaultDo {
	return v.withDO(v.DO.Returning(value, columns...))
}

func (v vaultDo) Not(conds ...gen.Condition) *vaultDo {
	return v.withDO(v.DO.Not(conds...))
}

func (v vaultDo) Or(conds ...gen.Condition) *vaultDo {
	return v.withDO(v.DO.Or(conds...))
}

func (v vaultDo) Select(conds ...field.Expr) *vaultDo {
	return v.withDO(v.DO.Select(conds...))
}

func (v vaultDo) Where(conds ...gen.Condition) *vaultDo {
	return v.withDO(v.DO.Where(conds...))
}

func (v vaultDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *vaultDo {
	return v.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (v vaultDo) Order(conds ...field.Expr) *vaultDo {
	return v.withDO(v.DO.Order(conds...))
}

func (v vaultDo) Distinct(cols ...field.Expr) *vaultDo {
	return v.withDO(v.DO.Distinct(cols...))
}

func (v vaultDo) Omit(cols ...field.Expr) *vaultDo {
	return v.withDO(v.DO.Omit(cols...))
}

func (v vaultDo) Join(table schema.Tabler, on ...field.Expr) *vaultDo {
	return v.withDO(v.DO.Join(table, on...))
}

func (v vaultDo) LeftJoin(table schema.Tabler, on ...field.Expr) *vaultDo {
	return v.withDO(v.DO.LeftJoin(table, on...))
}

func (v vaultDo) RightJoin(table schema.Tabler, on ...field.Expr) *vaultDo {
	return v.withDO(v.DO.RightJoin(table, on...))
}

func (v vaultDo) Group(cols ...field.Expr) *vaultDo {
	return v.withDO(v.DO.Group(cols...))
}

func (v vaultDo) Having(conds ...gen.Condition) *vaultDo {
	return v.withDO(v.DO.Having(conds...))
}

func (v vaultDo) Limit(limit int) *vaultDo {
	return v.withDO(v.DO.Limit(limit))
}

func (v vaultDo) Offset(offset int) *vaultDo {
	return v.withDO(v.DO.Offset(offset))
}

func (v vaultDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *vaultDo {
	return v.withDO(v.DO.Scopes(funcs...))
}

func (v vaultDo) Unscoped() *vaultDo {
	return v.withDO(v.DO.Unscoped())
}

func (v vaultDo) Create(values ...*model.Vault) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Create(values)
}

func (v vaultDo) CreateInBatches(values []*model.Vault, batchSize int) error {
	return v.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (v vaultDo) Save(values ...*model.Vault) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Save(values)
}

func (v vaultDo) First() (*model.Vault, error) {
	if result, err := v.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Vault), nil
	}
}

func (v vaultDo) Take() (*model.Vault, error) {
	if result, err := v.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Vault), nil
	}
}

func (v vaultDo) Last() (*model.Vault, error) {
	if result, err := v.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Vault), nil
	}
}

func (v vaultDo) Find() ([]*model.Vault, error) {
	result, err := v.DO.Find()
	return result.([]*model.Vault), err
}

func (v vaultDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Vault, err error) {
	buf := make([]*model.Vault, 0, batchSize)
	err = v.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (v vaultDo) FindInBatches(result *[]*model.Vault, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return v.DO.FindInBatches(result, batchSize, fc)
}

func (v vaultDo) Attrs(attrs ...field.AssignExpr) *vaultDo {
	return v.withDO(v.DO.Attrs(attrs...))
}

func (v vaultDo) Assign(attrs ...field.AssignExpr) *vaultDo {
	return v.withDO(v.DO.Assign(attrs...))
}

func (v vaultDo) Joins(fields ...field.RelationField) *vaultDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Joins(_f))
	}
	return &v
}

func (v vaultDo) Preload(fields ...field.RelationField) *vaultDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Preload(_f))
	}
	return &v
}

func (v vaultDo) FirstOrInit() (*model.Vault, error) {
	if result, err := v.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Vault), nil
	}
}

func (v vaultDo) FirstOrCreate() (*model.Vault, error) {
	if result, err := v.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Vault), nil
	}
}

func (v vaultDo) FindByPage(offset int, limit int) (result []*model.Vault, count int64, err error) {
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

func (v vaultDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = v.Count()
	if err != nil {
		return
	}

	err = v.Offset(offset).Limit(limit).Scan(result)
	return
}

func (v vaultDo) Scan(result interface{}) (err error) {
	return v.DO.Scan(result)
}

func (v *vaultDo) withDO(do gen.Dao) *vaultDo {
	v.DO = *do.(*gen.DO)
	return v
}
