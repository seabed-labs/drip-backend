package coingecko

import (
	"context"
	"testing"

	"github.com/test-go/testify/assert"
)

func TestTokenRegistry(t *testing.T) {

	client := NewCoinGeckoClient()

	t.Run("GetCoinGeckoMeta should return the required meta data", func(t *testing.T) {
		tokenAddress := "So11111111111111111111111111111111111111112"
		cgMeta, err := client.GetCoinGeckoMetadata(context.Background(), tokenAddress)
		assert.NoError(t, err)
		assert.Equal(t, "wrapped-solana", cgMeta.Id)
		assert.Equal(t, "sol", cgMeta.Symbol)
		assert.Equal(t, "Wrapped Solana", cgMeta.Name)
	})
}
