// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameProtoConfig = "proto_config"

// ProtoConfig mapped from table <proto_config>
type ProtoConfig struct {
	Pubkey               string `gorm:"column:pubkey;primaryKey" json:"pubkey" yaml:"pubkey"`
	Granularity          uint64 `gorm:"column:granularity;not null" json:"granularity" yaml:"granularity"`
	TriggerDcaSpread     uint16 `gorm:"column:trigger_dca_spread;not null" json:"triggerDcaSpread" yaml:"triggerDcaSpread"`
	BaseWithdrawalSpread uint16 `gorm:"column:base_withdrawal_spread;not null" json:"baseWithdrawalSpread" yaml:"baseWithdrawalSpread"`
}

// TableName ProtoConfig's table name
func (*ProtoConfig) TableName() string {
	return TableNameProtoConfig
}
