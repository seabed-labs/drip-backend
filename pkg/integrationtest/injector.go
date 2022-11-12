package integrationtest

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dcaf-labs/drip/pkg/api/middleware"
	controller "github.com/dcaf-labs/drip/pkg/api/routes"
	"github.com/dcaf-labs/drip/pkg/service/alert"
	"github.com/dcaf-labs/drip/pkg/service/base"
	"github.com/dcaf-labs/drip/pkg/service/clients"
	"github.com/dcaf-labs/drip/pkg/service/clients/coingecko"
	"github.com/dcaf-labs/drip/pkg/service/clients/orcawhirlpool"
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/clients/tokenregistry"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/database"
	"github.com/dcaf-labs/drip/pkg/service/repository/query"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

type TestOptions struct {
	FixturePath string
	AppConfig   config.AppConfig
	PSQLConfig  config.PSQLConfig
}

func InjectDependencies(
	testOptions *TestOptions,
	testCase interface{},
) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// only set if it is not explicitly set already
	if os.Getenv("SHOULD_USE_EMBEDDED_DB") == "" {
		if err := os.Setenv("SHOULD_USE_EMBEDDED_DB", "true"); err != nil {
			logrus.WithError(err).Error("could not set SHOULD_USE_EMBEDDED_DB env var")
			os.Exit(1)
		}
	}

	// test http recorder
	httpClientProvider := clients.DefaultClientProvider
	if testOptions != nil {
		r, err := recorder.New(testOptions.FixturePath)
		if err != nil {
			logrus.WithError(err).Error("could not get recorder")
			os.Exit(1)
		}
		defer func(r *recorder.Recorder) {
			if err := r.Stop(); err != nil {
				logrus.WithError(err).Error("could stop recorder")
				os.Exit(1)
			}
		}(r)
		if r.Mode() != recorder.ModeRecordOnce {
			logrus.Error("recorder should be in ModeRecordOnce")
			os.Exit(1)
		}
		r.SetReplayableInteractions(true)
		r.SetMatcher(requestMatcher)
		recorderHTTPClient := r.GetDefaultClient()
		httpClientProvider = func() clients.RetryableHTTPClientProvider {
			return func(options clients.RateLimitHTTPClientOptions) clients.RetryableHTTPClient {
				options.HttpClient = recorderHTTPClient
				options.CallsPerSecond = utils.GetIntPtr(100)
				return clients.DefaultClientProvider()(options)
			}
		}
	}
	providers := []interface{}{
		// Data access
		database.NewDatabase,
		database.NewGORMDatabase,
		query.Use,
		repository.NewRepository,
		repository.NewAccountUpdateQueue,
		// API Clients
		httpClientProvider,
		solana.NewSolanaClient,
		tokenregistry.NewTokenRegistry,
		orcawhirlpool.NewOrcaWhirlpoolClient,
		coingecko.NewCoinGeckoClient,
		// services
		base.NewBase,
		processor.NewProcessor,
		alert.NewAlertService,
		// server
		middleware.NewHandler,
		controller.NewHandler,
	}

	if testOptions != nil && testOptions.AppConfig != nil {
		providers = append(providers, func() config.AppConfig {
			return testOptions.AppConfig
		})
	} else {
		providers = append(providers, config.NewAppConfig)
	}
	if testOptions != nil && testOptions.PSQLConfig != nil {
		providers = append(providers, func() config.PSQLConfig {
			return testOptions.PSQLConfig
		})
	} else {
		providers = append(providers, config.NewPSQLConfig)
	}
	// comment out below for logs
	logrus.SetOutput(ioutil.Discard)
	opts := []fx.Option{
		fx.Provide(providers...),
		fx.Invoke(
			database.RunMigrations,
			testCase,
		),
		fx.NopLogger,
	}
	app := fx.New(opts...)
	defer func() {
		if err := app.Stop(context.Background()); err != nil {
			logrus.WithError(err).Errorf("failed to stop integrationtest app")
		}
	}()
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}

func requestMatcher(r *http.Request, i cassette.Request) bool {
	if r.Body == nil || r.Body == http.NoBody {
		return cassette.DefaultMatcher(r, i)
	}

	var reqBody []byte
	var err error
	reqBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Errorf("failed to read request body")
		os.Exit(1)
	}
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	return r.Method == i.Method && r.URL.String() == i.URL && string(reqBody) == i.Body
}
