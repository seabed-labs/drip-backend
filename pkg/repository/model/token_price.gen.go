// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameTokenPrice = "token_price"

// TokenPrice mapped from table <token_price>
type TokenPrice struct {
	ID     int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" db:"id"`
	Base   *string   `gorm:"column:base" json:"base" db:"base"`
	Quote  *string   `gorm:"column:quote" json:"quote" db:"quote"`
	Close  uint64    `gorm:"column:close;not null" json:"close" db:"close"`
	Date   time.Time `gorm:"column:date;not null" json:"date" db:"date"`
	Source *string   `gorm:"column:source" json:"source" db:"source"`
}

// TableName TokenPrice's table name
func (*TokenPrice) TableName() string {
	return TableNameTokenPrice
}