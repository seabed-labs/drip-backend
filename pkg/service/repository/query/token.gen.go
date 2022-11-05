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

func newToken(db *gorm.DB) token {
	_token := token{}

	_token.tokenDo.UseDB(db)
	_token.tokenDo.UseModel(&model.Token{})

	tableName := _token.tokenDo.TableName()
	_token.ALL = field.NewField(tableName, "*")
	_token.Pubkey = field.NewString(tableName, "pubkey")
	_token.Symbol = field.NewString(tableName, "symbol")
	_token.Decimals = field.NewInt16(tableName, "decimals")
	_token.IconURL = field.NewString(tableName, "icon_url")
	_token.CoinGeckoID = field.NewString(tableName, "coin_gecko_id")
	_token.UIMarketPrice = field.NewUint64(tableName, "ui_market_price")
	_token.Name = field.NewString(tableName, "name")
	_token.MarketCapRank = field.NewInt32(tableName, "market_cap_rank")

	_token.fillFieldMap()

	return _token
}

type token struct {
	tokenDo tokenDo

	ALL           field.Field
	Pubkey        field.String
	Symbol        field.String
	Decimals      field.Int16
	IconURL       field.String
	CoinGeckoID   field.String
	UIMarketPrice field.Uint64
	Name          field.String
	MarketCapRank field.Int32

	fieldMap map[string]field.Expr
}

func (t token) Table(newTableName string) *token {
	t.tokenDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t token) As(alias string) *token {
	t.tokenDo.DO = *(t.tokenDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *token) updateTableName(table string) *token {
	t.ALL = field.NewField(table, "*")
	t.Pubkey = field.NewString(table, "pubkey")
	t.Symbol = field.NewString(table, "symbol")
	t.Decimals = field.NewInt16(table, "decimals")
	t.IconURL = field.NewString(table, "icon_url")
	t.CoinGeckoID = field.NewString(table, "coin_gecko_id")
	t.UIMarketPrice = field.NewUint64(table, "ui_market_price")
	t.Name = field.NewString(table, "name")
	t.MarketCapRank = field.NewInt32(table, "market_cap_rank")

	t.fillFieldMap()

	return t
}

func (t *token) WithContext(ctx context.Context) *tokenDo { return t.tokenDo.WithContext(ctx) }

func (t token) TableName() string { return t.tokenDo.TableName() }

func (t token) Alias() string { return t.tokenDo.Alias() }

func (t *token) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *token) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 8)
	t.fieldMap["pubkey"] = t.Pubkey
	t.fieldMap["symbol"] = t.Symbol
	t.fieldMap["decimals"] = t.Decimals
	t.fieldMap["icon_url"] = t.IconURL
	t.fieldMap["coin_gecko_id"] = t.CoinGeckoID
	t.fieldMap["ui_market_price"] = t.UIMarketPrice
	t.fieldMap["name"] = t.Name
	t.fieldMap["market_cap_rank"] = t.MarketCapRank
}

func (t token) clone(db *gorm.DB) token {
	t.tokenDo.ReplaceDB(db)
	return t
}

type tokenDo struct{ gen.DO }

func (t tokenDo) Debug() *tokenDo {
	return t.withDO(t.DO.Debug())
}

func (t tokenDo) WithContext(ctx context.Context) *tokenDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t tokenDo) ReadDB() *tokenDo {
	return t.Clauses(dbresolver.Read)
}

func (t tokenDo) WriteDB() *tokenDo {
	return t.Clauses(dbresolver.Write)
}

func (t tokenDo) Clauses(conds ...clause.Expression) *tokenDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t tokenDo) Returning(value interface{}, columns ...string) *tokenDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t tokenDo) Not(conds ...gen.Condition) *tokenDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t tokenDo) Or(conds ...gen.Condition) *tokenDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t tokenDo) Select(conds ...field.Expr) *tokenDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t tokenDo) Where(conds ...gen.Condition) *tokenDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t tokenDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *tokenDo {
	return t.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (t tokenDo) Order(conds ...field.Expr) *tokenDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t tokenDo) Distinct(cols ...field.Expr) *tokenDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t tokenDo) Omit(cols ...field.Expr) *tokenDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t tokenDo) Join(table schema.Tabler, on ...field.Expr) *tokenDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t tokenDo) LeftJoin(table schema.Tabler, on ...field.Expr) *tokenDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t tokenDo) RightJoin(table schema.Tabler, on ...field.Expr) *tokenDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t tokenDo) Group(cols ...field.Expr) *tokenDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t tokenDo) Having(conds ...gen.Condition) *tokenDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t tokenDo) Limit(limit int) *tokenDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t tokenDo) Offset(offset int) *tokenDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t tokenDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *tokenDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t tokenDo) Unscoped() *tokenDo {
	return t.withDO(t.DO.Unscoped())
}

func (t tokenDo) Create(values ...*model.Token) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t tokenDo) CreateInBatches(values []*model.Token, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t tokenDo) Save(values ...*model.Token) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t tokenDo) First() (*model.Token, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Token), nil
	}
}

func (t tokenDo) Take() (*model.Token, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Token), nil
	}
}

func (t tokenDo) Last() (*model.Token, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Token), nil
	}
}

func (t tokenDo) Find() ([]*model.Token, error) {
	result, err := t.DO.Find()
	return result.([]*model.Token), err
}

func (t tokenDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Token, err error) {
	buf := make([]*model.Token, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t tokenDo) FindInBatches(result *[]*model.Token, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t tokenDo) Attrs(attrs ...field.AssignExpr) *tokenDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t tokenDo) Assign(attrs ...field.AssignExpr) *tokenDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t tokenDo) Joins(fields ...field.RelationField) *tokenDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t tokenDo) Preload(fields ...field.RelationField) *tokenDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t tokenDo) FirstOrInit() (*model.Token, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Token), nil
	}
}

func (t tokenDo) FirstOrCreate() (*model.Token, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Token), nil
	}
}

func (t tokenDo) FindByPage(offset int, limit int) (result []*model.Token, count int64, err error) {
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

func (t tokenDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t tokenDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t tokenDo) Delete(models ...*model.Token) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *tokenDo) withDO(do gen.Dao) *tokenDo {
	t.DO = *do.(*gen.DO)
	return t
}
