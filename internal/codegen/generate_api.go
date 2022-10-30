package codegen

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/deepmap/oapi-codegen/pkg/codegen"
	"github.com/deepmap/oapi-codegen/pkg/util"
	"github.com/sirupsen/logrus"
)

const apiDocs = "./docs/swagger.yaml"
const output = "./pkg/api/apispec/generated.go"

func GenerateAPIServer() {
	apiDocsPath := fmt.Sprintf("%s/%s", config.GetProjectRoot(), apiDocs)
	outputPath := fmt.Sprintf("%s/%s", config.GetProjectRoot(), output)
	logrus.
		WithField("apiDocsPath", apiDocsPath).
		WithField("outputPath", outputPath).
		Info("starting repo models codegen")
	swagger, err := util.LoadSwagger(apiDocs)
	if err != nil {
		logrus.WithError(err).Error("failed to LoadSwagger")
		os.Exit(1)
	}

	code, err := codegen.Generate(swagger, "apispec", codegen.Options{
		GenerateEchoServer: true,
		GenerateClient:     true,
		GenerateTypes:      true,
		EmbedSpec:          true,
		SkipFmt:            false,
		SkipPrune:          false,
		AliasTypes:         false,
		IncludeTags:        nil,
		ExcludeTags:        nil,
		UserTemplates:      nil,
		ImportMapping:      nil,
		ExcludeSchemas:     nil,
	})

	err = ioutil.WriteFile(outputPath, []byte(code), 0644)
	if err != nil {
		logrus.WithError(err).Error("error writing generated code to file")
		os.Exit(1)
	}
}
