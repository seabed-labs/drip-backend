package integrationtest

import (
	"context"
	"testing"

	"github.com/dcaf-labs/drip/pkg/integrationtest"
	"github.com/dcaf-labs/drip/pkg/service/clients/coingecko"
	"github.com/dcaf-labs/drip/pkg/service/price"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	"github.com/google/uuid"
	"github.com/test-go/testify/assert"
	"go.uber.org/fx"
)

type dummyInterface struct {
}

func (d dummyInterface) Append(_ fx.Hook) {

}
func TestPriceService(t *testing.T) {
	t.Run("should update market price for 2 tokens", func(t *testing.T) {
		integrationtest.InjectDependencies(
			&integrationtest.APIRecorderOptions{
				Path: "./fixtures/updateTokenMarketPricesForBatch-1",
			},
			func(
				repo repository.Repository,
				coinGeckoClient coingecko.CoinGeckoClient,
			) {
				priceService, err := price.NewPriceService(dummyInterface{}, repo, coinGeckoClient)
				assert.NoError(t, err)
				assert.NoError(t, repo.UpsertTokens(context.Background(), &model.Token{
					Pubkey:      "orcaEKTdK7LKz57vaAYr9QeNsVEPfiu6QeMU1kektZE",
					Symbol:      utils.GetStringPtr("ORCA"),
					Decimals:    6,
					IconURL:     utils.GetStringPtr("https://raw.githubusercontent.com/solana-labs/token-list/main/assets/mainnet/orcaEKTdK7LKz57vaAYr9QeNsVEPfiu6QeMU1kektZE/logo.png"),
					CoinGeckoID: utils.GetStringPtr("orca"),
				}, &model.Token{
					Pubkey:      "mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So",
					Symbol:      utils.GetStringPtr("mSOL"),
					Decimals:    9,
					IconURL:     utils.GetStringPtr("https://raw.githubusercontent.com/solana-labs/token-list/main/assets/mainnet/mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So/logo.png"),
					CoinGeckoID: utils.GetStringPtr("msol"),
				}))
				tokenPairID := uuid.New()
				assert.NoError(t, repo.InsertTokenPairs(context.Background(), &model.TokenPair{
					ID:     tokenPairID.String(),
					TokenA: "orcaEKTdK7LKz57vaAYr9QeNsVEPfiu6QeMU1kektZE",
					TokenB: "mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So",
				}))
				assert.NoError(t, repo.UpsertProtoConfigs(context.Background(), &model.ProtoConfig{
					Pubkey:                  "DG43NUhq6gxUAcNGCJ45zCrBHgPDNNjMovocPHdHUqUD",
					Granularity:             60,
					TokenADripTriggerSpread: 50,
					TokenBWithdrawalSpread:  50,
					Admin:                   "CMRE2hAFTSBs8cTaDdF8qfwjorYeAsrGoeP7tJKHPbBR",
					TokenBReferralSpread:    50,
				}))
				assert.NoError(t, repo.UpsertVaults(context.Background(), &model.Vault{
					Pubkey:                "dcaw8KiS6YbsPFWYwkgA4ztfajE4CJN4rmzDkD6TAMb",
					ProtoConfig:           "DG43NUhq6gxUAcNGCJ45zCrBHgPDNNjMovocPHdHUqUD",
					TokenAAccount:         "J15X2DWTRPVTJsofDrf5se4zkNv1sD1eJPgEHwvuNJer",
					TokenBAccount:         "J15X2DWTRPVTJsofDrf5se4zkNv1sD1eJPgEHwvuNJer",
					TreasuryTokenBAccount: "J15X2DWTRPVTJsofDrf5se4zkNv1sD1eJPgEHwvuNJer",
					TokenPairID:           tokenPairID.String(),
					TokenAMint:            "orcaEKTdK7LKz57vaAYr9QeNsVEPfiu6QeMU1kektZE",
					TokenBMint:            "mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So",
					Enabled:               true,
				}))

				priceService.UpdateTokenMarketPrices()

				tokens, err := repo.GetTokensByAddresses(context.Background(), "orcaEKTdK7LKz57vaAYr9QeNsVEPfiu6QeMU1kektZE", "mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So")
				assert.NoError(t, err)
				for _, token := range tokens {
					assert.NotNil(t, token)
					assert.NotNil(t, token.UIMarketPrice)
					assert.NotZero(t, token.UIMarketPrice)
					assert.NotNil(t, token.MarketCapRank)
				}
			})
	})
}
