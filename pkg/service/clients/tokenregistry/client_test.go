package tokenregistry

import (
	"testing"

	"github.com/test-go/testify/assert"
)

func TestTokenRegistry(t *testing.T) {

	client := NewTokenRegistry()

	t.Run("GetTokenRegistry should return all tokens", func(t *testing.T) {
		tokenRegistry, err := client.GetTokenRegistry()
		assert.NoError(t, err)
		assert.True(t, len(tokenRegistry.Tokens) >= 13644)
	})

	t.Run("GetTokenRegistry should return all tokens", func(t *testing.T) {
		token, err := client.GetTokenRegistryToken("DUSTcnwRpZjhds1tLY2NpcvVTmKL6JJERD9T274LcqCr")
		assert.NoError(t, err)
		assert.Equal(t, token.Address, "DUSTcnwRpZjhds1tLY2NpcvVTmKL6JJERD9T274LcqCr")
	})
}
