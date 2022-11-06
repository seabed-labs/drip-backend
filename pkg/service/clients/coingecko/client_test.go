package coingecko

import (
	"context"
	"testing"

	"github.com/dcaf-labs/drip/pkg/service/clients"

	"github.com/test-go/testify/assert"
)

func TestTokenRegistry(t *testing.T) {

	client := NewCoinGeckoClient(clients.DefaultClientProvider())

	t.Run("GetCoinGeckoMeta should return the required meta data", func(t *testing.T) {
		tokenAddress := "So11111111111111111111111111111111111111112"
		cgMeta, err := client.GetCoinGeckoMetadata(context.Background(), tokenAddress)
		assert.NoError(t, err)
		assert.Equal(t, "wrapped-solana", cgMeta.ID)
		assert.Equal(t, "sol", cgMeta.Symbol)
		assert.Equal(t, "Wrapped Solana", cgMeta.Name)
	})

	t.Run("GetCoinGeckoMeta should return market tokenjob for 3 tokens", func(t *testing.T) {
		res, err := client.GetMarketPriceForTokens(context.Background(), "bonfida", "honey-finance", "solend")
		assert.NoError(t, err)
		assert.Len(t, res, 3)
	})

	t.Run("GetCoinGeckoMeta should return error if coinGeckoIDS are empty", func(t *testing.T) {
		res, err := client.GetMarketPriceForTokens(context.Background())
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "missing coinGeckoIDs", err.Error())
	})

	t.Run("GetSolanaCoinsList should return list of all solana assets", func(t *testing.T) {
		res, err := client.GetSolanaCoinsList(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, res)
		for _, coin := range res {
			assert.NotNil(t, coin.Platforms.Solana)
			assert.NotEqual(t, "", coin.Platforms.Solana)
		}
	})
}
