// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"reflect"

	"github.com/shopspring/decimal"
)

const TableNameVaultPeriod = "vault_period"

// VaultPeriod mapped from table <vault_period>
type VaultPeriod struct {
	Pubkey      string          `gorm:"column:pubkey;type:varchar;primaryKey" json:"pubkey" db:"pubkey"`
	Vault       string          `gorm:"column:vault;type:varchar;not null" json:"vault" db:"vault"`
	PeriodID    uint64          `gorm:"column:period_id;type:numeric;not null" json:"periodId" db:"period_id"`
	Twap        decimal.Decimal `gorm:"column:twap;type:numeric;not null" json:"twap" db:"twap"`
	Dar         uint64          `gorm:"column:dar;type:numeric;not null" json:"dar" db:"dar"`
	PriceBOverA decimal.Decimal `gorm:"column:price_b_over_a;type:numeric" json:"priceBOverA" db:"price_b_over_a"`
}

// TableName VaultPeriod's table name
func (*VaultPeriod) TableName() string {
	return TableNameVaultPeriod
}

func (t VaultPeriod) GetAllColumns() []string {
	var res []string
	numFields := reflect.TypeOf(t).NumField()
	i := 0
	for i < numFields {
		field := reflect.TypeOf(t).Field(i)
		res = append(res, field.Tag.Get("db"))
		i++
	}
	return res
}
