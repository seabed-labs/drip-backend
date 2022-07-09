// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameTokenAccountBalance = "token_account_balance"

// TokenAccountBalance mapped from table <token_account_balance>
type TokenAccountBalance struct {
	Pubkey string `gorm:"column:pubkey;primaryKey" json:"pubkey" db:"pubkey"`
	Mint   string `gorm:"column:mint;not null" json:"mint" db:"mint"`
	Owner  string `gorm:"column:owner;not null" json:"owner" db:"owner"`
	Amount uint64 `gorm:"column:amount;not null" json:"amount" db:"amount"`
	State  string `gorm:"column:state;not null" json:"state" db:"state"`
}

// TableName TokenAccountBalance's table name
func (*TokenAccountBalance) TableName() string {
	return TableNameTokenAccountBalance
}