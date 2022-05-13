package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type AppConfig struct {
	Environment Environment `yaml:"environment" env:"ENV"`
	Wallet      string      `yaml:"wallet"      env:"DRIP_BACKEND_WALLET"`
	Port        int         `yaml:"port"        env:"PORT"`
	//VaultConfigs []VaultConfig `yaml:"vaults"`
}

//type VaultConfig struct {
//	Vault                       string `yaml:"vault"`
//	VaultProtoConfig            string `yaml:"vaultProtoConfig"`
//	VaultProtoConfigGranularity uint32 `yaml:"vaultProtoConfigGranularity"`
//	VaultTokenAAccount          string `yaml:"vaultTokenAAccount"`
//	VaultTokenBAccount          string `yaml:"vaultTokenBAccount"`
//	VaultTreasuryTokenBAccount  string `yaml:"vaultTreasuryTokenBAccount"`
//	TokenAMint                  string `yaml:"tokenAMint"`
//	TokenASymbol                string `yaml:"tokenASymbol"`
//	TokenBMint                  string `yaml:"tokenBMint"`
//	TokenBSymbol                string `yaml:"tokenBSymbol"`
//	SwapTokenMint               string `yaml:"swapTokenMint"`
//	SwapTokenAAccount           string `yaml:"swapTokenAAccount"`
//	SwapTokenBAccount           string `yaml:"swapTokenBAccount"`
//	SwapFeeAccount              string `yaml:"swapFeeAccount"`
//	SwapAuthority               string `yaml:"swapAuthority"`
//	Swap                        string `yaml:"swap"`
//}

type PSQLConfig struct {
	User     string `yaml:"psql_username" env:"PSQL_USER"`
	Password string `yaml:"psql_password" env:"PSQL_PASS"`
	DBName   string `yaml:"psql_database" env:"PSQL_DBNAME"`
	Port     int    `yaml:"psql_port" env:"PSQL_PORT"`
	Host     string `yaml:"psql_host" env:"PSQL_HOST"`
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

// Note: config has to be a pointer
func ParseToConfig(config interface{}, configFileName string) error {
	LoadEnv()
	//configFileName := fmt.Sprintf("./internal/configs/%s.yaml", environment)
	//configFileName = fmt.Sprintf("%s/%s", GetProjectRoot(), configFileName)
	if configFileName != "" {
		log.WithField("configFileName", configFileName).Infof("loading configs file")
		if err := cleanenv.ReadConfig(configFileName, config); err != nil {
			log.WithField("configFileName", configFileName).Warning("config file does not exist")
		}
	}
	return cleanenv.ReadEnv(config)
}

func NewAppConfig() (*AppConfig, error) {
	var config AppConfig
	if err := ParseToConfig(&config, ""); err != nil {
		return nil, err
	}
	log.Info("loaded drip-backend app configs")
	return &config, nil
}

func NewPSQLConfig() (*PSQLConfig, error) {
	var config PSQLConfig
	if err := ParseToConfig(&config, ""); err != nil {
		return nil, err
	}
	log.Info("loaded drip-backend app configs")
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
