package configs

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Environment  Environment   `yaml:"environment" env:"ENV"`
	Wallet       string        `yaml:"wallet"      env:"DRIP_BACKEND_WALLET"`
	Port         int           `yaml:"port"        env:"PORT"`
	VaultConfigs []VaultConfig `yaml:"vaults"`
}

type VaultConfig struct {
	Vault                       string `yaml:"vault"`
	VaultProtoConfig            string `yaml:"vaultProtoConfig"`
	VaultProtoConfigGranularity uint32 `yaml:"vaultProtoConfigGranularity"`
	VaultTokenAAccount          string `yaml:"vaultTokenAAccount"`
	VaultTokenBAccount          string `yaml:"vaultTokenBAccount"`
	VaultTreasuryTokenBAccount  string `yaml:"vaultTreasuryTokenBAccount"`
	TokenAMint                  string `yaml:"tokenAMint"`
	TokenASymbol                string `yaml:"tokenASymbol"`
	TokenBMint                  string `yaml:"tokenBMint"`
	TokenBSymbol                string `yaml:"tokenBSymbol"`
	SwapTokenMint               string `yaml:"swapTokenMint"`
	SwapTokenAAccount           string `yaml:"swapTokenAAccount"`
	SwapTokenBAccount           string `yaml:"swapTokenBAccount"`
	SwapFeeAccount              string `yaml:"swapFeeAccount"`
	SwapAuthority               string `yaml:"swapAuthority"`
	Swap                        string `yaml:"swap"`
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
	configFileName := fmt.Sprintf("./internal/pkg/configs/%s.yaml", environment)
	configFileName = fmt.Sprintf("%s/%s", GetProjectRoot(), configFileName)

	log.WithField("configFileName", configFileName).Infof("loading configs file")
	configFile, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	defer func(configFile *os.File) {
		if err := configFile.Close(); err != nil {
			log.WithError(err).Errorf("failed to close configs file")
		}
	}(configFile)
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

func IsMainnet(env Environment) bool {
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
