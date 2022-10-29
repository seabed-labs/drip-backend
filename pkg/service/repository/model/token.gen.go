// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameToken = "token"

// Token mapped from table <token>
type Token struct {
	Pubkey      string  `gorm:"column:pubkey;type:varchar;primaryKey" json:"pubkey" db:"pubkey"`
	Symbol      *string `gorm:"column:symbol;type:varchar" json:"symbol" db:"symbol"`
	Decimals    int16   `gorm:"column:decimals;type:int2;not null" json:"decimals" db:"decimals"`
	IconURL     *string `gorm:"column:icon_url;type:varchar" json:"iconUrl" db:"icon_url"`
	CoinGeckoID *string `gorm:"column:coin_gecko_id;type:varchar" json:"coinGeckoId" db:"coin_gecko_id"`
}

// TableName Token's table name
func (*Token) TableName() string {
	return TableNameToken
}
