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

func newWithdrawalMetric(db *gorm.DB) withdrawalMetric {
	_withdrawalMetric := withdrawalMetric{}

	_withdrawalMetric.withdrawalMetricDo.UseDB(db)
	_withdrawalMetric.withdrawalMetricDo.UseModel(&model.WithdrawalMetric{})

	tableName := _withdrawalMetric.withdrawalMetricDo.TableName()
	_withdrawalMetric.ALL = field.NewField(tableName, "*")
	_withdrawalMetric.Signature = field.NewString(tableName, "signature")
	_withdrawalMetric.IxIndex = field.NewInt32(tableName, "ix_index")
	_withdrawalMetric.IxName = field.NewString(tableName, "ix_name")
	_withdrawalMetric.IxVersion = field.NewInt32(tableName, "ix_version")
	_withdrawalMetric.Slot = field.NewInt32(tableName, "slot")
	_withdrawalMetric.Time = field.NewTime(tableName, "time")
	_withdrawalMetric.Vault = field.NewString(tableName, "vault")
	_withdrawalMetric.TokenAMint = field.NewString(tableName, "token_a_mint")
	_withdrawalMetric.TokenBMint = field.NewString(tableName, "token_b_mint")
	_withdrawalMetric.UserTokenAWithdrawAmount = field.NewUint64(tableName, "user_token_a_withdraw_amount")
	_withdrawalMetric.UserTokenBWithdrawAmount = field.NewUint64(tableName, "user_token_b_withdraw_amount")
	_withdrawalMetric.TreasuryTokenBReceivedAmount = field.NewUint64(tableName, "treasury_token_b_received_amount")
	_withdrawalMetric.ReferralTokenBReceivedAmount = field.NewUint64(tableName, "referral_token_b_received_amount")
	_withdrawalMetric.TokenAUsdPriceDay = field.NewUint64(tableName, "token_a_usd_price_day")
	_withdrawalMetric.TokenBUsdPriceDay = field.NewUint64(tableName, "token_b_usd_price_day")

	_withdrawalMetric.fillFieldMap()

	return _withdrawalMetric
}

type withdrawalMetric struct {
	withdrawalMetricDo withdrawalMetricDo

	ALL                          field.Field
	Signature                    field.String
	IxIndex                      field.Int32
	IxName                       field.String
	IxVersion                    field.Int32
	Slot                         field.Int32
	Time                         field.Time
	Vault                        field.String
	TokenAMint                   field.String
	TokenBMint                   field.String
	UserTokenAWithdrawAmount     field.Uint64
	UserTokenBWithdrawAmount     field.Uint64
	TreasuryTokenBReceivedAmount field.Uint64
	ReferralTokenBReceivedAmount field.Uint64
	TokenAUsdPriceDay            field.Uint64
	TokenBUsdPriceDay            field.Uint64

	fieldMap map[string]field.Expr
}

func (w withdrawalMetric) Table(newTableName string) *withdrawalMetric {
	w.withdrawalMetricDo.UseTable(newTableName)
	return w.updateTableName(newTableName)
}

func (w withdrawalMetric) As(alias string) *withdrawalMetric {
	w.withdrawalMetricDo.DO = *(w.withdrawalMetricDo.As(alias).(*gen.DO))
	return w.updateTableName(alias)
}

func (w *withdrawalMetric) updateTableName(table string) *withdrawalMetric {
	w.ALL = field.NewField(table, "*")
	w.Signature = field.NewString(table, "signature")
	w.IxIndex = field.NewInt32(table, "ix_index")
	w.IxName = field.NewString(table, "ix_name")
	w.IxVersion = field.NewInt32(table, "ix_version")
	w.Slot = field.NewInt32(table, "slot")
	w.Time = field.NewTime(table, "time")
	w.Vault = field.NewString(table, "vault")
	w.TokenAMint = field.NewString(table, "token_a_mint")
	w.TokenBMint = field.NewString(table, "token_b_mint")
	w.UserTokenAWithdrawAmount = field.NewUint64(table, "user_token_a_withdraw_amount")
	w.UserTokenBWithdrawAmount = field.NewUint64(table, "user_token_b_withdraw_amount")
	w.TreasuryTokenBReceivedAmount = field.NewUint64(table, "treasury_token_b_received_amount")
	w.ReferralTokenBReceivedAmount = field.NewUint64(table, "referral_token_b_received_amount")
	w.TokenAUsdPriceDay = field.NewUint64(table, "token_a_usd_price_day")
	w.TokenBUsdPriceDay = field.NewUint64(table, "token_b_usd_price_day")

	w.fillFieldMap()

	return w
}

func (w *withdrawalMetric) WithContext(ctx context.Context) *withdrawalMetricDo {
	return w.withdrawalMetricDo.WithContext(ctx)
}

func (w withdrawalMetric) TableName() string { return w.withdrawalMetricDo.TableName() }

func (w withdrawalMetric) Alias() string { return w.withdrawalMetricDo.Alias() }

func (w *withdrawalMetric) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := w.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (w *withdrawalMetric) fillFieldMap() {
	w.fieldMap = make(map[string]field.Expr, 15)
	w.fieldMap["signature"] = w.Signature
	w.fieldMap["ix_index"] = w.IxIndex
	w.fieldMap["ix_name"] = w.IxName
	w.fieldMap["ix_version"] = w.IxVersion
	w.fieldMap["slot"] = w.Slot
	w.fieldMap["time"] = w.Time
	w.fieldMap["vault"] = w.Vault
	w.fieldMap["token_a_mint"] = w.TokenAMint
	w.fieldMap["token_b_mint"] = w.TokenBMint
	w.fieldMap["user_token_a_withdraw_amount"] = w.UserTokenAWithdrawAmount
	w.fieldMap["user_token_b_withdraw_amount"] = w.UserTokenBWithdrawAmount
	w.fieldMap["treasury_token_b_received_amount"] = w.TreasuryTokenBReceivedAmount
	w.fieldMap["referral_token_b_received_amount"] = w.ReferralTokenBReceivedAmount
	w.fieldMap["token_a_usd_price_day"] = w.TokenAUsdPriceDay
	w.fieldMap["token_b_usd_price_day"] = w.TokenBUsdPriceDay
}

func (w withdrawalMetric) clone(db *gorm.DB) withdrawalMetric {
	w.withdrawalMetricDo.ReplaceDB(db)
	return w
}

type withdrawalMetricDo struct{ gen.DO }

func (w withdrawalMetricDo) Debug() *withdrawalMetricDo {
	return w.withDO(w.DO.Debug())
}

func (w withdrawalMetricDo) WithContext(ctx context.Context) *withdrawalMetricDo {
	return w.withDO(w.DO.WithContext(ctx))
}

func (w withdrawalMetricDo) ReadDB() *withdrawalMetricDo {
	return w.Clauses(dbresolver.Read)
}

func (w withdrawalMetricDo) WriteDB() *withdrawalMetricDo {
	return w.Clauses(dbresolver.Write)
}

func (w withdrawalMetricDo) Clauses(conds ...clause.Expression) *withdrawalMetricDo {
	return w.withDO(w.DO.Clauses(conds...))
}

func (w withdrawalMetricDo) Returning(value interface{}, columns ...string) *withdrawalMetricDo {
	return w.withDO(w.DO.Returning(value, columns...))
}

func (w withdrawalMetricDo) Not(conds ...gen.Condition) *withdrawalMetricDo {
	return w.withDO(w.DO.Not(conds...))
}

func (w withdrawalMetricDo) Or(conds ...gen.Condition) *withdrawalMetricDo {
	return w.withDO(w.DO.Or(conds...))
}

func (w withdrawalMetricDo) Select(conds ...field.Expr) *withdrawalMetricDo {
	return w.withDO(w.DO.Select(conds...))
}

func (w withdrawalMetricDo) Where(conds ...gen.Condition) *withdrawalMetricDo {
	return w.withDO(w.DO.Where(conds...))
}

func (w withdrawalMetricDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *withdrawalMetricDo {
	return w.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (w withdrawalMetricDo) Order(conds ...field.Expr) *withdrawalMetricDo {
	return w.withDO(w.DO.Order(conds...))
}

func (w withdrawalMetricDo) Distinct(cols ...field.Expr) *withdrawalMetricDo {
	return w.withDO(w.DO.Distinct(cols...))
}

func (w withdrawalMetricDo) Omit(cols ...field.Expr) *withdrawalMetricDo {
	return w.withDO(w.DO.Omit(cols...))
}

func (w withdrawalMetricDo) Join(table schema.Tabler, on ...field.Expr) *withdrawalMetricDo {
	return w.withDO(w.DO.Join(table, on...))
}

func (w withdrawalMetricDo) LeftJoin(table schema.Tabler, on ...field.Expr) *withdrawalMetricDo {
	return w.withDO(w.DO.LeftJoin(table, on...))
}

func (w withdrawalMetricDo) RightJoin(table schema.Tabler, on ...field.Expr) *withdrawalMetricDo {
	return w.withDO(w.DO.RightJoin(table, on...))
}

func (w withdrawalMetricDo) Group(cols ...field.Expr) *withdrawalMetricDo {
	return w.withDO(w.DO.Group(cols...))
}

func (w withdrawalMetricDo) Having(conds ...gen.Condition) *withdrawalMetricDo {
	return w.withDO(w.DO.Having(conds...))
}

func (w withdrawalMetricDo) Limit(limit int) *withdrawalMetricDo {
	return w.withDO(w.DO.Limit(limit))
}

func (w withdrawalMetricDo) Offset(offset int) *withdrawalMetricDo {
	return w.withDO(w.DO.Offset(offset))
}

func (w withdrawalMetricDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *withdrawalMetricDo {
	return w.withDO(w.DO.Scopes(funcs...))
}

func (w withdrawalMetricDo) Unscoped() *withdrawalMetricDo {
	return w.withDO(w.DO.Unscoped())
}

func (w withdrawalMetricDo) Create(values ...*model.WithdrawalMetric) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Create(values)
}

func (w withdrawalMetricDo) CreateInBatches(values []*model.WithdrawalMetric, batchSize int) error {
	return w.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (w withdrawalMetricDo) Save(values ...*model.WithdrawalMetric) error {
	if len(values) == 0 {
		return nil
	}
	return w.DO.Save(values)
}

func (w withdrawalMetricDo) First() (*model.WithdrawalMetric, error) {
	if result, err := w.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.WithdrawalMetric), nil
	}
}

func (w withdrawalMetricDo) Take() (*model.WithdrawalMetric, error) {
	if result, err := w.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.WithdrawalMetric), nil
	}
}

func (w withdrawalMetricDo) Last() (*model.WithdrawalMetric, error) {
	if result, err := w.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.WithdrawalMetric), nil
	}
}

func (w withdrawalMetricDo) Find() ([]*model.WithdrawalMetric, error) {
	result, err := w.DO.Find()
	return result.([]*model.WithdrawalMetric), err
}

func (w withdrawalMetricDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.WithdrawalMetric, err error) {
	buf := make([]*model.WithdrawalMetric, 0, batchSize)
	err = w.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (w withdrawalMetricDo) FindInBatches(result *[]*model.WithdrawalMetric, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return w.DO.FindInBatches(result, batchSize, fc)
}

func (w withdrawalMetricDo) Attrs(attrs ...field.AssignExpr) *withdrawalMetricDo {
	return w.withDO(w.DO.Attrs(attrs...))
}

func (w withdrawalMetricDo) Assign(attrs ...field.AssignExpr) *withdrawalMetricDo {
	return w.withDO(w.DO.Assign(attrs...))
}

func (w withdrawalMetricDo) Joins(fields ...field.RelationField) *withdrawalMetricDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Joins(_f))
	}
	return &w
}

func (w withdrawalMetricDo) Preload(fields ...field.RelationField) *withdrawalMetricDo {
	for _, _f := range fields {
		w = *w.withDO(w.DO.Preload(_f))
	}
	return &w
}

func (w withdrawalMetricDo) FirstOrInit() (*model.WithdrawalMetric, error) {
	if result, err := w.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.WithdrawalMetric), nil
	}
}

func (w withdrawalMetricDo) FirstOrCreate() (*model.WithdrawalMetric, error) {
	if result, err := w.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.WithdrawalMetric), nil
	}
}

func (w withdrawalMetricDo) FindByPage(offset int, limit int) (result []*model.WithdrawalMetric, count int64, err error) {
	result, err = w.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = w.Offset(-1).Limit(-1).Count()
	return
}

func (w withdrawalMetricDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = w.Count()
	if err != nil {
		return
	}

	err = w.Offset(offset).Limit(limit).Scan(result)
	return
}

func (w withdrawalMetricDo) Scan(result interface{}) (err error) {
	return w.DO.Scan(result)
}

func (w withdrawalMetricDo) Delete(models ...*model.WithdrawalMetric) (result gen.ResultInfo, err error) {
	return w.DO.Delete(models)
}

func (w *withdrawalMetricDo) withDO(do gen.Dao) *withdrawalMetricDo {
	w.DO = *do.(*gen.DO)
	return w
}
