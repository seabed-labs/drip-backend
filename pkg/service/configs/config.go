package configs

import (
	"github.com/dcaf-labs/solana-go-clients/pkg/drip"
	ag_solanago "github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
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

func NewAppConfig() (AppConfig, error) {
	var config appConfig
	if err := ParseToConfig(&config, ""); err != nil {
		return nil, err
	}
	// Sane defaults
	if config.Environment == NilEnv {
		config.Environment = StagingEnv
	}
	if config.Network == NilNetwork {
		config.Network = DevnetNetwork
	}

	log.Info("loaded drip-backend app configs")
	drip.ProgramID = ag_solanago.MustPublicKeyFromBase58(config.DripProgramID)
	log.
		WithField("programID", drip.ProgramID.String()).
		WithField("ShouldByPassAdminAuth", config.ShouldByPassAdminAuth).
		Info("set programID")
	return config, nil
}

func NewPSQLConfig() (PSQLConfig, error) {
	var config psqlConfig
	if err := ParseToConfig(&config, ""); err != nil {
		return nil, err
	}
	if config.IsTestDB {
		config.DBName = "test_" + uuid.New().String()[0:4]
	}
	log.
		WithField("IsTestDB", config.IsTestDB).
		WithField("DBName", config.DBName).
		Info("loaded drip-backend app configs")
	return &config, nil
}

func IsStagingEnvironment(env Environment) bool {
	return env == StagingEnv || env == NilEnv
}

func IsProductionEnvironment(env Environment) bool {
	return env == ProdEnv
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
