package psql

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/dcaf-protocol/drip/internal/configs"
	"gorm.io/gen"
)

const modelDir = "./internal/pkg/repository/model"
const queryDir = "./internal/pkg/repository"

func GenerateModels(
	db *gorm.DB,
) {
	modelPkgPath := fmt.Sprintf("%s/%s", configs.GetProjectRoot(), modelDir)
	outPath := fmt.Sprintf("%s/%s", configs.GetProjectRoot(), queryDir)
	logrus.
		WithField("modelPkgPath", modelPkgPath).
		WithField("outPath", outPath).
		Info("starting repo models codegen")
	g := gen.NewGenerator(gen.Config{
		OutPath:        outPath,
		ModelPkgPath:   modelPkgPath,
		FieldNullable:  true,
		FieldCoverable: true,
		FieldSignable:  true,
		//fieldNewTagNS:
	})
	// TODO(mocha): convert snake case column name to camel case
	g.WithNewTagNameStrategy(func(columnName string) string {
		return fmt.Sprintf("yaml:\"%s\"", columnName)
	})
	g.UseDB(db)
	tables := g.GenerateAllTable()
	g.ApplyBasic(tables...)
	g.GenerateModel("proto_config",
		gen.FieldType("granularity", "uint64"),
		gen.FieldType("trigger_dca_spread", "uint16"),
		gen.FieldType("base_withdrawal_spread", "uint16"),
	)
	g.GenerateModel("vault",
		gen.FieldType("last_dca_period", "uint64"),
		gen.FieldType("drip_amount", "uint64"),
	)
	g.GenerateModel("vault_period",
		gen.FieldType("period_id", "uint64"),
		gen.FieldType("dar", "uint64"),
		gen.FieldType("twap", "decimal.Decimal"),
	)
	g.Execute()
}
