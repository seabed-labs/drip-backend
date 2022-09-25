package tokenregistry

import (
	"context"
	"testing"

	"github.com/test-go/testify/assert"
)

func TestTokenRegistry(t *testing.T) {

	client := NewTokenRegistry()

	t.Run("GetTokenRegistry should return all tokens", func(t *testing.T) {
		tokenRegistry, err := client.GetTokenRegistry(context.Background())
		assert.NoError(t, err)
		assert.True(t, len(tokenRegistry.Tokens) >= 13644)
	})

	t.Run("GetTokenRegistry should return all tokens", func(t *testing.T) {
		token, err := client.GetTokenRegistryToken(context.Background(), "DUSTcnwRpZjhds1tLY2NpcvVTmKL6JJERD9T274LcqCr")
		assert.NoError(t, err)
		assert.Equal(t, token.Address, "DUSTcnwRpZjhds1tLY2NpcvVTmKL6JJERD9T274LcqCr")
	})
}
