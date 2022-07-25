// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"reflect"

	"github.com/shopspring/decimal"
)

const TableNameOrcaWhirlpool = "orca_whirlpool"

// OrcaWhirlpool mapped from table <orca_whirlpool>
type OrcaWhirlpool struct {
	Pubkey                     string          `gorm:"column:pubkey;primaryKey" json:"pubkey" db:"pubkey"`
	WhirlpoolsConfig           string          `gorm:"column:whirlpools_config;not null" json:"whirlpoolsConfig" db:"whirlpools_config"`
	TokenMintA                 string          `gorm:"column:token_mint_a;not null" json:"tokenMintA" db:"token_mint_a"`
	TokenVaultA                string          `gorm:"column:token_vault_a;not null" json:"tokenVaultA" db:"token_vault_a"`
	TokenMintB                 string          `gorm:"column:token_mint_b;not null" json:"tokenMintB" db:"token_mint_b"`
	TokenVaultB                string          `gorm:"column:token_vault_b;not null" json:"tokenVaultB" db:"token_vault_b"`
	TickSpacing                int32           `gorm:"column:tick_spacing;not null" json:"tickSpacing" db:"tick_spacing"`
	FeeRate                    int32           `gorm:"column:fee_rate;not null" json:"feeRate" db:"fee_rate"`
	ProtocolFeeRate            int32           `gorm:"column:protocol_fee_rate;not null" json:"protocolFeeRate" db:"protocol_fee_rate"`
	TickCurrentIndex           int32           `gorm:"column:tick_current_index;not null" json:"tickCurrentIndex" db:"tick_current_index"`
	ProtocolFeeOwedA           decimal.Decimal `gorm:"column:protocol_fee_owed_a;not null" json:"protocolFeeOwedA" db:"protocol_fee_owed_a"`
	ProtocolFeeOwedB           decimal.Decimal `gorm:"column:protocol_fee_owed_b;not null" json:"protocolFeeOwedB" db:"protocol_fee_owed_b"`
	RewardLastUpdatedTimestamp decimal.Decimal `gorm:"column:reward_last_updated_timestamp;not null" json:"rewardLastUpdatedTimestamp" db:"reward_last_updated_timestamp"`
	Liquidity                  decimal.Decimal `gorm:"column:liquidity;not null" json:"liquidity" db:"liquidity"`
	SqrtPrice                  decimal.Decimal `gorm:"column:sqrt_price;not null" json:"sqrtPrice" db:"sqrt_price"`
	FeeGrowthGlobalA           decimal.Decimal `gorm:"column:fee_growth_global_a;not null" json:"feeGrowthGlobalA" db:"fee_growth_global_a"`
	FeeGrowthGlobalB           decimal.Decimal `gorm:"column:fee_growth_global_b;not null" json:"feeGrowthGlobalB" db:"fee_growth_global_b"`
	TokenPairID                string          `gorm:"column:token_pair_id;not null" json:"tokenPairId" db:"token_pair_id"`
}

// TableName OrcaWhirlpool's table name
func (*OrcaWhirlpool) TableName() string {
	return TableNameOrcaWhirlpool
}

func (t OrcaWhirlpool) GetAllColumns() []string {
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
