// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameSourceReference = "source_reference"

// SourceReference mapped from table <source_reference>
type SourceReference struct {
	Value string `gorm:"column:value;primaryKey" json:"value"`
}

// TableName SourceReference's table name
func (*SourceReference) TableName() string {
	return TableNameSourceReference
}
