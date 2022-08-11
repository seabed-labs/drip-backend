// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import "reflect"

const TableNameProtoConfig = "proto_config"

// ProtoConfig mapped from table <proto_config>
type ProtoConfig struct {
	Pubkey               string `gorm:"column:pubkey;type:varchar;primaryKey" json:"pubkey" db:"pubkey"`
	Granularity          uint64 `gorm:"column:granularity;type:numeric;not null" json:"granularity" db:"granularity"`
	TriggerDcaSpread     uint16 `gorm:"column:trigger_dca_spread;type:int2;not null" json:"triggerDcaSpread" db:"trigger_dca_spread"`
	BaseWithdrawalSpread uint16 `gorm:"column:base_withdrawal_spread;type:int2;not null" json:"baseWithdrawalSpread" db:"base_withdrawal_spread"`
	Admin                string `gorm:"column:admin;type:varchar;not null" json:"admin" db:"admin"`
}

// TableName ProtoConfig's table name
func (*ProtoConfig) TableName() string {
	return TableNameProtoConfig
}

func (t ProtoConfig) GetAllColumns() []string {
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
