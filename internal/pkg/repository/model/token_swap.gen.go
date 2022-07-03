// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameTokenSwap = "token_swap"

// TokenSwap mapped from table <token_swap>
type TokenSwap struct {
	Pubkey        string `gorm:"column:pubkey;type:varchar;not null" json:"pubkey" db:"pubkey"`
	Mint          string `gorm:"column:mint;type:varchar;not null" json:"mint" db:"mint"`
	Authority     string `gorm:"column:authority;type:varchar;not null" json:"authority" db:"authority"`
	FeeAccount    string `gorm:"column:fee_account;type:varchar;not null" json:"feeAccount" db:"fee_account"`
	TokenAAccount string `gorm:"column:token_a_account;type:varchar;not null" json:"tokenAAccount" db:"token_a_account"`
	TokenBAccount string `gorm:"column:token_b_account;type:varchar;not null" json:"tokenBAccount" db:"token_b_account"`
	TokenPairID   string `gorm:"column:token_pair_id;type:uuid;not null" json:"tokenPairId" db:"token_pair_id"`
	TokenAMint    string `gorm:"column:token_a_mint;type:varchar;not null" json:"tokenAMint" db:"token_a_mint"`
	TokenBMint    string `gorm:"column:token_b_mint;type:varchar;not null" json:"tokenBMint" db:"token_b_mint"`
	ID            string `gorm:"column:id;type:uuid;primaryKey" json:"id" db:"id"`
}

// TableName TokenSwap's table name
func (*TokenSwap) TableName() string {
	return TableNameTokenSwap
}
