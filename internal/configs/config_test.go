package configs

import (
	"testing"

	"github.com/test-go/testify/assert"
)

func TestNewConfig(t *testing.T) {
	config, err := NewAppConfig()
	assert.NoError(t, err)
	assert.Contains(t, []Environment{NilEnv, LocalnetEnv, DevnetEnv}, config.Environment)
}

func TestIsDev(t *testing.T) {
	assert.Equal(t, IsDev(DevnetEnv), true)
	assert.Equal(t, IsDev(LocalnetEnv), false)
	assert.Equal(t, IsDev("random"), false)
}

func TestIsLocal(t *testing.T) {
	assert.Equal(t, IsLocal(NilEnv), true)
	assert.Equal(t, IsLocal(LocalnetEnv), true)
	assert.Equal(t, IsLocal(DevnetEnv), false)
	assert.Equal(t, IsLocal("random"), false)
}

func TestIsProd(t *testing.T) {
	assert.Equal(t, IsMainnet(MainnetEnv), true)
	assert.Equal(t, IsMainnet(LocalnetEnv), false)
	assert.Equal(t, IsMainnet("random"), false)
}
