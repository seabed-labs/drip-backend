package configs

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Environment Environment `yaml:"environment" env:"ENV"`
	Wallet      string      `yaml:"wallet"      env:"DRIP_BACKEND_WALLET"`
	Port        int         `yaml:"port"        env:"DRIP_BACKEND_PORT"`
}

type Environment string

const (
	NilEnv      = Environment("")
	LocalnetEnv = Environment("LOCALNET")
	DevnetEnv   = Environment("DEVNET")
	MainnetEnv  = Environment("MAINNET")
)

type EnvVar string

const (
	ENV                   EnvVar = "ENV"
	PROJECT_ROOT_OVERRIDE EnvVar = "PROJECT_ROOT_OVERRIDE"
)

func NewConfig() (*Config, error) {
	LoadEnv()

	environment := GetEnv(Environment(os.Getenv(string(ENV))))

	configFileName := fmt.Sprintf("./internal/configs/%s.yaml", environment)
	configFileName = fmt.Sprintf("%s/%s", GetProjectRoot(), configFileName)

	log.WithField("configFileName", configFileName).Infof("loading config file")
	configFile, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config Config
	if err := cleanenv.ReadConfig(configFileName, &config); err != nil {
		return nil, err
	}

	log.Info("loaded drip-backend configs")

	return &config, nil
}

func IsLocal(env Environment) bool {
	return env == LocalnetEnv || env == NilEnv
}

func IsDev(env Environment) bool {
	return env == DevnetEnv
}

func IsProd(env Environment) bool {
	return env == MainnetEnv
}

func GetEnv(env Environment) Environment {
	switch env {
	case MainnetEnv:
		return MainnetEnv
	case DevnetEnv:
		return DevnetEnv
	case NilEnv:
		fallthrough
	case LocalnetEnv:
		fallthrough
	default:
		return LocalnetEnv
	}
}
