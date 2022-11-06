package integrationtest

import (
	"context"
	"testing"

	"github.com/dcaf-labs/drip/pkg/integrationtest"
	"github.com/dcaf-labs/drip/pkg/job/token"
	"github.com/dcaf-labs/drip/pkg/service/clients/coingecko"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/unittest"
	"github.com/gagliardetto/solana-go"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/test-go/testify/assert"
	"go.uber.org/fx"
)

type dummyInterface struct {
}

func (d dummyInterface) Append(_ fx.Hook) {

}
func TestPriceService(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}
	ctrl := gomock.NewController(t)
	t.Run("should update market tokenjob for 2 tokens", func(t *testing.T) {
		mockConfig := config.NewMockAppConfig(ctrl)
		mockConfig.EXPECT().GetWalletPrivateKey().Return(unittest.GetTestPrivateKey()).AnyTimes()
		mockConfig.EXPECT().GetDiscordWebhookID().Return("").AnyTimes()
		mockConfig.EXPECT().GetDiscordWebhookAccessToken().Return("").AnyTimes()
		mockConfig.EXPECT().GetNetwork().Return(config.MainnetNetwork).AnyTimes()
		mockConfig.EXPECT().GetEnvironment().Return(config.ProductionEnv).AnyTimes()
		mockConfig.EXPECT().GetServerPort().Return(8080).AnyTimes()
		integrationtest.InjectDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/updateTokenMarketPricesForBatch-1",
				AppConfig:   mockConfig,
			},
			func(
				repo repository.Repository,
				processor processor.Processor,
				coinGeckoClient coingecko.CoinGeckoClient,
			) {
				orcaTokenAddress := "orcaEKTdK7LKz57vaAYr9QeNsVEPfiu6QeMU1kektZE"
				msolTokenAddress := "mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So"

				setupEnabledTokens(t, repo, orcaTokenAddress, msolTokenAddress)

				priceService, err := token.NewTokenJob(dummyInterface{}, repo, processor, coinGeckoClient)
				assert.NoError(t, err)
				priceService.UpsertAllSupportedTokensWithMetadata()

				orca, err := repo.GetTokenByAddress(context.Background(), orcaTokenAddress)
				assert.NoError(t, err)
				assert.NotNil(t, orca)
				msol, err := repo.GetTokenByAddress(context.Background(), msolTokenAddress)
				assert.NoError(t, err)
				assert.NotNil(t, msol)

				assert.Equal(t, orcaTokenAddress, orca.Pubkey)
				assert.Equal(t, int16(6), orca.Decimals)
				assert.Equal(t, "ORCA", *orca.Symbol)
				assert.Equal(t, "Orca", *orca.Name)
				assert.Equal(t, "https://assets.coingecko.com/coins/images/17547/small/Orca_Logo.png?1628781615", *orca.IconURL)
				assert.Equal(t, "orca", *orca.CoinGeckoID)
				assert.Equal(t, 0.893816, *orca.UIMarketPrice)
				assert.Equal(t, int32(683), *orca.MarketCapRank)

				assert.Equal(t, msolTokenAddress, msol.Pubkey)
				assert.Equal(t, int16(9), msol.Decimals)
				assert.Equal(t, "mSOL", *msol.Symbol)
				assert.Equal(t, "Marinade staked SOL (mSOL)", *msol.Name)
				assert.Equal(t, "https://assets.coingecko.com/coins/images/17752/small/mSOL.png?1644541955", *msol.IconURL)
				assert.Equal(t, "msol", *msol.CoinGeckoID)
				assert.Equal(t, 36.52, *msol.UIMarketPrice)
				assert.Equal(t, int32(155), *msol.MarketCapRank)

			})
	})
}

func setupEnabledTokens(t *testing.T, repo repository.Repository, orcaTokenAddress string, msolTokenAddress string) {
	assert.NoError(t, repo.UpsertTokens(context.Background(), &model.Token{
		Pubkey:   orcaTokenAddress,
		Decimals: 6,
	}, &model.Token{
		Pubkey:   msolTokenAddress,
		Decimals: 9,
	}))
	tokenPairID := uuid.New()
	protoConfig := solana.NewWallet().PublicKey().String()
	assert.NoError(t, repo.InsertTokenPairs(context.Background(), &model.TokenPair{
		ID:     tokenPairID.String(),
		TokenA: orcaTokenAddress,
		TokenB: msolTokenAddress,
	}))
	assert.NoError(t, repo.UpsertProtoConfigs(context.Background(), &model.ProtoConfig{
		Pubkey:                  protoConfig,
		Granularity:             60,
		TokenADripTriggerSpread: 50,
		TokenBWithdrawalSpread:  50,
		Admin:                   solana.NewWallet().PublicKey().String(),
		TokenBReferralSpread:    50,
	}))
	assert.NoError(t, repo.UpsertVaults(context.Background(), &model.Vault{
		Pubkey:                solana.NewWallet().PublicKey().String(),
		ProtoConfig:           protoConfig,
		TokenAAccount:         solana.NewWallet().PublicKey().String(),
		TokenBAccount:         solana.NewWallet().PublicKey().String(),
		TreasuryTokenBAccount: solana.NewWallet().PublicKey().String(),
		TokenPairID:           tokenPairID.String(),
		TokenAMint:            orcaTokenAddress,
		TokenBMint:            msolTokenAddress,
		Enabled:               true,
	}))
}
