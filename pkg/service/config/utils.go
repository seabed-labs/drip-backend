package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

// Note: config has to be a pointer
func parseToConfig(config interface{}, configFileName string) error {
	LoadEnv()
	if configFileName != "" {
		log.WithField("configFileName", configFileName).Infof("loading config file")
		if err := cleanenv.ReadConfig(configFileName, config); err != nil {
			log.WithField("configFileName", configFileName).Warning("config file does not exist")
		}
	}
	return cleanenv.ReadEnv(config)
}

func IsStagingEnvironment(env Environment) bool {
	return env == StagingEnv || env == NilEnv
}

func IsProductionEnvironment(env Environment) bool {
	return env == ProdEnv || env == ProductionEnv
}

func IsLocalNetwork(network Network) bool {
	return network == LocalNetwork || network == NilNetwork
}

func IsDevnetNetwork(network Network) bool {
	return network == DevnetNetwork
}

func IsMainnetNetwork(network Network) bool {
	return network == MainnetNetwork
}
