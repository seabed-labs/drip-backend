package token

import (
	"context"
	"fmt"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/processor"

	"github.com/dcaf-labs/drip/pkg/service/clients/coingecko"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/utils"
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
	if _, err := s.Every(5).Minutes().Do(impl.UpsertAllSupportedTokensWithMetadata); err != nil {
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

func (i impl) UpsertAllSupportedTokensWithMetadata() {
	ctx := context.Background()
	tokens, err := i.repo.GetAllSupportedTokens(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to GetSupportedTokenAs")
		return
	}
	if err := utils.DoForPaginatedBatch(coingecko.CoinsMarketsPathLimit, len(tokens), func(start, end int) {
		if err := i.updateTokenMetadataForBatch(ctx, tokens[start:end]); err != nil {
			logrus.WithError(err).Error("failed to updateTokenMetadataForBatch")
		}
	}); err != nil {
		logrus.WithError(err).Error("failed to DoForPaginatedBatch")
		return
	}
}

func (i impl) updateTokenMetadataForBatch(ctx context.Context, tokens []*model.Token) error {
	// backfill token metadata
	for _, token := range tokens {
		// this is inefficient on a  fresh db as it will make n network calls
		// subsequent runs should produce better results as data is backfilled
		if err := i.processor.UpsertTokenByAddress(ctx, token.Pubkey); err != nil {
			logrus.WithError(err).Error("failed to UpsertTokenByAddress, continuing...")
		}
	}
	// re-fetch tokens to get tokens with populated coingecko id's
	tokens, err := i.repo.GetTokensByAddresses(ctx, model.GetTokenPubkeys(tokens)...)
	if err != nil {
		return err
	}
	// backfill latest market price
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
				WithError(fmt.Errorf("unexpected missing token")).
				WithField("cgID", marketPrice.ID).
				Warning("missing token, continuing...")
			continue
		}
		token.UIMarketPrice = utils.GetFloat64Ptr(marketPrice.CurrentPrice)
		token.MarketCapRank = marketPrice.MarketCapRank
		tokensToUpsert = append(tokensToUpsert, token)
	}
	return i.repo.UpsertTokens(ctx, tokensToUpsert...)
}
