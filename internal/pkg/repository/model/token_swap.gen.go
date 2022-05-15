// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameTokenSwap = "token_swap"

// TokenSwap mapped from table <token_swap>
type TokenSwap struct {
	Pubkey        string `gorm:"column:pubkey;primaryKey" json:"pubkey" yaml:"pubkey"`
	Mint          string `gorm:"column:mint;not null" json:"mint" yaml:"mint"`
	Authority     string `gorm:"column:authority;not null" json:"authority" yaml:"authority"`
	FeeAccount    string `gorm:"column:fee_account;not null" json:"fee_account" yaml:"fee_account"`
	TokenAAccount string `gorm:"column:token_a_account;not null" json:"token_a_account" yaml:"token_a_account"`
	TokenBAccount string `gorm:"column:token_b_account;not null" json:"token_b_account" yaml:"token_b_account"`
	Pair          string `gorm:"column:pair;not null" json:"pair" yaml:"pair"`
}

// TableName TokenSwap's table name
func (*TokenSwap) TableName() string {
	return TableNameTokenSwap
}