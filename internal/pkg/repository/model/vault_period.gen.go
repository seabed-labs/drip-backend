// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import "github.com/shopspring/decimal"

const TableNameVaultPeriod = "vault_period"

// VaultPeriod mapped from table <vault_period>
type VaultPeriod struct {
	Pubkey   string          `gorm:"column:pubkey;primaryKey" json:"pubkey" yaml:"pubkey"`
	Vault    string          `gorm:"column:vault;not null" json:"vault" yaml:"vault"`
	PeriodID uint64          `gorm:"column:period_id;not null" json:"periodId" yaml:"periodId"`
	Twap     decimal.Decimal `gorm:"column:twap;not null" json:"twap" yaml:"twap"`
	Dar      uint64          `gorm:"column:dar;not null" json:"dar" yaml:"dar"`
}

// TableName VaultPeriod's table name
func (*VaultPeriod) TableName() string {
	return TableNameVaultPeriod
}
