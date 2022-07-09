package psql

import (
	"fmt"
	"reflect"

	"github.com/dcaf-protocol/drip/pkg/configs"

	"github.com/iancoleman/strcase"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"gorm.io/gen"
)

const modelDir = "./internal/pkg/repository/model"
const queryDir = "./internal/pkg/repository/query"

type ModelUtil struct{}

func (t ModelUtil) GetAllColumns() []string {
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

func GenerateModels(
	db *gorm.DB,
) {
	modelPkgPath := fmt.Sprintf("%s/%s", configs.GetProjectRoot(), modelDir)
	queryPkgPath := fmt.Sprintf("%s/%s", configs.GetProjectRoot(), queryDir)
	logrus.
		WithField("modelPkgPath", modelPkgPath).
		WithField("outPath", queryPkgPath).
		Info("starting repo models codegen")
	g := gen.NewGenerator(gen.Config{
		OutPath:        queryPkgPath,
		ModelPkgPath:   modelPkgPath,
		FieldNullable:  true,
		FieldCoverable: true,
		FieldSignable:  true,
	})
	// TODO(Mocha): map UUID
	dataMap := map[string]func(detailType string) (dataType string){
		"numeric": func(detailType string) (dataType string) { return "uint64" },
	}
	g.WithDataTypeMap(dataMap)

	g.WithNewTagNameStrategy(func(columnName string) string {
		return fmt.Sprintf("yaml:\"%s\"", strcase.ToLowerCamel(columnName))
	})
	g.WithNewTagNameStrategy(func(columnName string) string {
		return fmt.Sprintf("db:\"%s\"", strcase.ToSnake(columnName))
	})
	g.WithJSONTagNameStrategy(strcase.ToLowerCamel)
	g.UseDB(db)
	tables := g.GenerateAllTable()
	g.ApplyBasic(tables...)

	g.GenerateModel("proto_config",
		gen.FieldType("granularity", "uint64"),
		gen.FieldType("trigger_dca_spread", "uint16"),
		gen.FieldType("base_withdrawal_spread", "uint16"),
	).AddMethod(ModelUtil{})

	g.GenerateModel("vault",
		gen.FieldType("last_dca_period", "uint64"),
		gen.FieldType("drip_amount", "uint64"),
	).AddMethod(ModelUtil{})

	g.GenerateModel("vault_period",
		gen.FieldType("period_id", "uint64"),
		gen.FieldType("dar", "uint64"),
		gen.FieldType("twap", "decimal.Decimal"),
	).AddMethod(ModelUtil{})

	g.GenerateModel("position",
		gen.FieldType("deposited_token_a_amount", "uint64"),
		gen.FieldType("withdrawn_token_b_amount", "uint64"),
		gen.FieldType("dca_period_id_before_deposit", "uint64"),
		gen.FieldType("number_of_swaps", "uint64"),
		gen.FieldType("periodic_drip_amount", "uint64"),
	).AddMethod(ModelUtil{})

	g.Execute()
}