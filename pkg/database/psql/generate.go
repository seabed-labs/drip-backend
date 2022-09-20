package psql

import (
	"fmt"
	"reflect"

	"github.com/dcaf-labs/drip/pkg/configs"

	"github.com/iancoleman/strcase"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"gorm.io/gen"
)

const modelDir = "./pkg/repository/model"
const queryDir = "./pkg/repository/query"

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
		gen.FieldType("token_a_drip_trigger_spread", "uint16"),
		gen.FieldType("token_b_withdrawal_spread", "uint16"),
		gen.FieldType("token_b_referral_spread", "uint16"),
	).AddMethod(ModelUtil{})

	g.GenerateModel("vault",
		gen.FieldType("last_dca_period", "uint64"),
		gen.FieldType("drip_amount", "uint64"),
	).AddMethod(ModelUtil{})

	g.GenerateModel("vault_period",
		gen.FieldType("period_id", "uint64"),
		gen.FieldType("dar", "uint64"),
		gen.FieldType("twap", "decimal.Decimal"),
		gen.FieldType("price_b_over_a", "decimal.Decimal"),
	).AddMethod(ModelUtil{})

	g.GenerateModel("position",
		gen.FieldType("deposited_token_a_amount", "uint64"),
		gen.FieldType("withdrawn_token_b_amount", "uint64"),
		gen.FieldType("dca_period_id_before_deposit", "uint64"),
		gen.FieldType("number_of_swaps", "uint64"),
		gen.FieldType("periodic_drip_amount", "uint64"),
	).AddMethod(ModelUtil{})

	g.GenerateModel("orca_whirlpool",
		gen.FieldType("protocol_fee_owed_a", "decimal.Decimal"),
		gen.FieldType("protocol_fee_owed_b", "decimal.Decimal"),
		gen.FieldType("reward_last_updated_timestamp", "decimal.Decimal"),
		gen.FieldType("liquidity", "decimal.Decimal"),
		gen.FieldType("sqrt_price", "decimal.Decimal"),
		gen.FieldType("fee_growth_global_a", "decimal.Decimal"),
		gen.FieldType("fee_growth_global_b", "decimal.Decimal"),
	).AddMethod(ModelUtil{})

	g.Execute()
}
