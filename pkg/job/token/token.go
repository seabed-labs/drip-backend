package token

import (
	"context"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/clients/coingecko"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type TokenJob interface {
	UpsertAllSupportedTokensWithMetadata()
}

type impl struct {
	repo            repository.Repository
	processor       processor.Processor
	coinGeckoClient coingecko.CoinGeckoClient
}

func NewTokenJob(
	lifecycle fx.Lifecycle,
	repo repository.Repository,
	processor processor.Processor,
	coinGeckoClient coingecko.CoinGeckoClient,
) (TokenJob, error) {
	impl := impl{
		repo:            repo,
		processor:       processor,
		coinGeckoClient: coinGeckoClient,
	}
	s := gocron.NewScheduler(time.UTC)
	if _, err := s.Every(30).Minutes().Do(impl.UpsertAllSupportedTokensWithMetadata); err != nil {
		return nil, err
	}
	s.StartImmediately().StartAsync()
	lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			s.Stop()
			return nil
		},
	})
	return impl, nil
}

func (i impl) UpsertAllSupportedTokensWithMetadata() {
	logrus.Info("starting UpsertAllSupportedTokensWithMetadata")
	ctx := context.Background()
	tokens, err := i.repo.GetAllSupportedTokens(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to GetSupportedTokenAs")
		return
	}
	if err := i.processor.UpsertTokensByAddresses(ctx, model.GetTokenPubkeys(tokens)...); err != nil {
		logrus.WithError(err).WithField("len(tokens", len(tokens)).Error("failed to UpsertTokensByAddresses")
	}
	logrus.Info("done UpsertAllSupportedTokensWithMetadata")
}
