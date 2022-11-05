package price

import (
	"context"
	"fmt"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/clients/coingecko"

	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type PriceService interface {
	UpdateTokenMarketPrices()
}

type impl struct {
	repo            repository.Repository
	coinGeckoClient coingecko.CoinGeckoClient
}

func NewPriceService(
	lifecycle fx.Lifecycle,
	repo repository.Repository,
	coinGeckoClient coingecko.CoinGeckoClient,
) (PriceService, error) {
	impl := impl{
		repo:            repo,
		coinGeckoClient: coinGeckoClient,
	}
	s := gocron.NewScheduler(time.UTC)
	if _, err := s.Every(5).Minutes().Do(impl.UpdateTokenMarketPrices); err != nil {
		return nil, err
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			s.StartImmediately()
			return nil
		},
		OnStop: func(_ context.Context) error {
			s.Stop()
			return nil
		},
	})
	return impl, nil
}

func (i impl) UpdateTokenMarketPrices() {
	ctx := context.Background()
	tokens, err := i.repo.GetAllSupportedTokens(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to GetSupportedTokenAs")
		return
	}
	if err := utils.DoForPaginatedBatch(coingecko.CoinsMarketsPathLimit, len(tokens), func(start, end int) {
		if err := i.updateTokenMarketPricesForBatch(ctx, tokens[start:end]); err != nil {
			logrus.WithError(err).Error("failed to updateTokenMarketPricesForBatch")
		}
	}); err != nil {
		logrus.WithError(err).Error("failed to DoForPaginatedBatch")
		return
	}
}

func (i impl) updateTokenMarketPricesForBatch(ctx context.Context, tokens []*model.Token) error {
	tokensByCoinGeckoID := model.GetTokensByCoinGeckoID(tokens)
	coinGeckoIDs := model.GetTokenCoinGeckoIDs(tokens)
	marketPrices, err := i.coinGeckoClient.GetMarketPriceForTokens(ctx, coinGeckoIDs...)
	if err != nil {
		return err
	}
	tokensToUpsert := []*model.Token{}
	for _, marketPrice := range marketPrices {
		token, ok := tokensByCoinGeckoID[marketPrice.ID]
		if !ok {
			logrus.
				WithError(fmt.Errorf("unexted missing token")).
				WithField("cgID", marketPrice.ID).Warning("missing token, continuing")
			continue
		}
		token.UIMarketPrice = utils.GetFloat64Ptr(marketPrice.CurrentPrice)
		token.MarketCapRank = marketPrice.MarketCapRank
		tokensToUpsert = append(tokensToUpsert, token)
	}
	return i.repo.UpsertTokens(ctx, tokensToUpsert...)
}
