package config

import (
	"github.com/dcaf-labs/solana-go-clients/pkg/drip"
	ag_solanago "github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Network string

const (
	NilNetwork     = Network("")
	LocalNetwork   = Network("LOCALNET")
	DevnetNetwork  = Network("DEVNET")
	MainnetNetwork = Network("MAINNET")
)

type Environment string

const (
	NilEnv        = Environment("")
	StagingEnv    = Environment("STAGING")
	ProdEnv       = Environment("PROD")
	ProductionEnv = Environment("PRODUCTION")
)

type EnvVar string

const (
	ENV                   EnvVar = "ENV"
	NETWORK               EnvVar = "NETWORK"
	PROJECT_ROOT_OVERRIDE EnvVar = "PROJECT_ROOT_OVERRIDE"
)

type AppConfig interface {
	GetEnvironment() Environment
	GetNetwork() Network
	GetDripProgramID() string
	GetGoogleClientID() string
	GetWalletPrivateKey() string
	GetServerPort() int
	GetDiscordWebhookID() string
	GetDiscordWebhookAccessToken() string
	GetShouldByPassAdminAuth() bool
}

type appConfig struct {
	Environment               Environment `yaml:"environment" env:"ENV" env-default:"STAGING"`
	Network                   Network     `yaml:"network" env:"NETWORK" env-default:"DEVNET"`
	DripProgramID             string      `yaml:"dripProgramID" env:"DRIP_PROGRAM_ID"  env-default:"dripTrkvSyQKvkyWg7oi4jmeEGMA5scSYowHArJ9Vwk"`
	GoogleClientID            string      `yaml:"googleClientID" env:"GOOGLE_CLIENT_ID"  env-default:"540992596258-sa2h4lmtelo44tonpu9htsauk5uabdon.apps.googleusercontent.com"`
	Wallet                    string      `yaml:"wallet"      env:"DRIP_BACKEND_WALLET"`
	Port                      int         `yaml:"port"        env:"PORT"`
	DiscordWebhookID          string      `yaml:"discordWebhookID" env:"DISCORD_WEBHOOK_ID"`
	DiscordWebhookAccessToken string      `yaml:"discordWebhookAccessToken" env:"DISCORD_ACCESS_TOKEN"`
	ShouldByPassAdminAuth     bool        `yaml:"shouldBypassAdminAuth" env:"SHOULD_BYPASS_ADMIN_AUTH" env-default:"false"`
}

func (a appConfig) GetShouldByPassAdminAuth() bool {
	return a.ShouldByPassAdminAuth
}

func (a appConfig) GetEnvironment() Environment {
	return a.Environment
}

func (a appConfig) GetNetwork() Network {
	return a.Network
}

func (a appConfig) GetDripProgramID() string {
	return a.DripProgramID
}

func (a appConfig) GetGoogleClientID() string {
	return a.GoogleClientID
}

func (a appConfig) GetWalletPrivateKey() string {
	return a.Wallet
}

func (a appConfig) GetServerPort() int {
	return a.Port
}

func (a appConfig) GetDiscordWebhookID() string {
	return a.DiscordWebhookID
}

func (a appConfig) GetDiscordWebhookAccessToken() string {
	return a.DiscordWebhookAccessToken
}

type PSQLConfig interface {
	GetURL() string
	GetUser() string
	GetPassword() string
	GetDBName() string
	GetPort() int
	GetHost() string
	GetIsTestDB() bool

	SetDBName(newDBName string) string
}

type psqlConfig struct {
	URL      string `env:"DATABASE_URL"`
	User     string `yaml:"psql_username" env:"PSQL_USER"`
	Password string `yaml:"psql_password" env:"PSQL_PASS"`
	DBName   string `yaml:"psql_database" env:"PSQL_DBNAME"`
	Port     int    `yaml:"psql_port" env:"PSQL_PORT"`
	Host     string `yaml:"psql_host" env:"PSQL_HOST"`
	IsTestDB bool   `yaml:"is_test_db" env:"IS_TEST_DB"`
}

func (p *psqlConfig) GetURL() string {
	return p.URL
}

func (p *psqlConfig) GetUser() string {
	return p.User
}

func (p *psqlConfig) GetPassword() string {
	return p.Password
}

func (p *psqlConfig) GetDBName() string {
	return p.DBName
}

func (p *psqlConfig) SetDBName(newDBName string) string {
	p.DBName = newDBName
	return p.DBName
}

func (p *psqlConfig) GetPort() int {
	return p.Port
}

func (p *psqlConfig) GetHost() string {
	return p.Host
}

func (p *psqlConfig) GetIsTestDB() bool {
	return p.IsTestDB
}

func NewAppConfig() (AppConfig, error) {
	var config appConfig
	if err := parseToConfig(&config, ""); err != nil {
		return nil, err
	}
	// Sane defaults
	if config.Environment == NilEnv {
		config.Environment = StagingEnv
	} else if config.Environment == ProdEnv {
		// ProdEnv is deprecated
		config.Environment = ProductionEnv
	}
	if config.Network == NilNetwork {
		config.Network = DevnetNetwork
	}

	log.Info("loaded drip-backend app config")
	drip.SetProgramID(ag_solanago.MustPublicKeyFromBase58(config.DripProgramID))
	log.
		WithField("programID", drip.ProgramID.String()).
		WithField("ShouldByPassAdminAuth", config.ShouldByPassAdminAuth).
		Info("set programID")
	return config, nil
}

func NewPSQLConfig() (PSQLConfig, error) {
	var config psqlConfig
	if err := parseToConfig(&config, ""); err != nil {
		return nil, err
	}
	if config.IsTestDB {
		config.DBName = "test_" + uuid.New().String()[0:4]
	}
	log.
		WithField("IsTestDB", config.IsTestDB).
		WithField("DBName", config.DBName).
		Info("loaded drip-backend app config")
	return &config, nil
}
