package configs

import (
	"github.com/dcaf-labs/solana-go-clients/pkg/drip"
	ag_solanago "github.com/gagliardetto/solana-go"
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type AppConfig struct {
	Environment    Environment `yaml:"environment" env:"ENV" env-default:"STAGING"`
	Network        Network     `yaml:"network" env:"NETWORK" env-default:"DEVNET"`
	DripProgramID  string      `yaml:"dripProgramID" env:"DRIP_PROGRAM_ID"  env-default:"dripTrkvSyQKvkyWg7oi4jmeEGMA5scSYowHArJ9Vwk"`
	GoogleClientID string      `yaml:"googleClientID" env:"GOOGLE_CLIENT_ID"  env-default:"540992596258-sa2h4lmtelo44tonpu9htsauk5uabdon.apps.googleusercontent.com"`
	Wallet         string      `yaml:"wallet"      env:"DRIP_BACKEND_WALLET"`
	Port           int         `yaml:"port"        env:"PORT"`
}

type PSQLConfig struct {
	URL      string `env:"DATABASE_URL"`
	User     string `yaml:"psql_username" env:"PSQL_USER"`
	Password string `yaml:"psql_password" env:"PSQL_PASS"`
	DBName   string `yaml:"psql_database" env:"PSQL_DBNAME"`
	Port     int    `yaml:"psql_port" env:"PSQL_PORT"`
	Host     string `yaml:"psql_host" env:"PSQL_HOST"`
	IsTestDB bool   `yaml:"is_test_db" env:"IS_TEST_DB"`
}

type Network string

const (
	NilNetwork     = Network("")
	LocalNetwork   = Network("LOCALNET")
	DevnetNetwork  = Network("DEVNET")
	MainnetNetwork = Network("MAINNET")
)

type Environment string

const (
	NilEnv     = Environment("")
	StagingEnv = Environment("STAGING")
	ProdEnv    = Environment("PROD")
)

type EnvVar string

const (
	ENV                   EnvVar = "ENV"
	PROJECT_ROOT_OVERRIDE EnvVar = "PROJECT_ROOT_OVERRIDE"
)

// Note: config has to be a pointer
func ParseToConfig(config interface{}, configFileName string) error {
	LoadEnv()
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
	drip.ProgramID = ag_solanago.MustPublicKeyFromBase58(config.DripProgramID)
	log.
		WithField("programID", drip.ProgramID.String()).
		Info("set programID")
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

func IsStaging(env Environment) bool {
	return env == StagingEnv || env == NilEnv
}

func IsProd(env Environment) bool {
	return env == ProdEnv
}

func IsLocal(network Network) bool {
	return network == LocalNetwork || network == NilNetwork
}

func IsDevnet(network Network) bool {
	return network == DevnetNetwork
}

func IsMainnet(network Network) bool {
	return network == MainnetNetwork
}
