package configs

import (
	"fmt"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const PROJECT_DIR = "drip-backend"

// LoadEnv loads env vars from .env at root of repo
func GetProjectRoot() string {
	rootOverride := os.Getenv(string(PROJECT_ROOT_OVERRIDE))
	if rootOverride != "" {
		log.WithField("override", rootOverride).Infof("override project root")
		return rootOverride
	}
	re := regexp.MustCompile(`^(.*` + PROJECT_DIR + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	return string(rootPath)
}

func LoadEnv() {
	log.WithField("env", Environment(os.Getenv(string(ENV)))).Infof("loading env")
	re := regexp.MustCompile(`^(.*` + PROJECT_DIR + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	filePath := fmt.Sprintf("%s/.env", string(rootPath))
	err := godotenv.Load(filePath)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"cwd":      cwd,
			"filePath": filePath,
		}).Warning("problem loading .env file")
	}
}
