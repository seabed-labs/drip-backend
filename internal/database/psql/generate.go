package psql

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/dcaf-protocol/drip/internal/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
)

const modelDir = "./internal/pkg/repository/model"
const queryDir = "./internal/pkg/repository/query"

func GenerateModels(
	config *configs.PSQLConfig,
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
	})
	db, err := gorm.Open(postgres.Open(getConnectionString(config)), &gorm.Config{})
	if err != nil {
		return
	}
	g.UseDB(db)
	tables := g.GenerateAllTable()
	g.ApplyBasic(tables...)
	g.Execute()
}
