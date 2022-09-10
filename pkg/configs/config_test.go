package configs

import (
	"testing"

	"github.com/test-go/testify/assert"
)

func TestNewAppConfig(t *testing.T) {
	config, err := NewAppConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Contains(t, []Network{NilNetwork, LocalNetwork, DevnetNetwork}, config.Network)
}

func TestNewPSQLConfig(t *testing.T) {
	config, err := NewPSQLConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestIsDev(t *testing.T) {
	assert.Equal(t, IsDevnet(DevnetNetwork), true)
	assert.Equal(t, IsDevnet(LocalNetwork), false)
	assert.Equal(t, IsDevnet("random"), false)
}

func TestIsLocal(t *testing.T) {
	assert.Equal(t, IsLocal(NilNetwork), true)
	assert.Equal(t, IsLocal(LocalNetwork), true)
	assert.Equal(t, IsLocal(DevnetNetwork), false)
	assert.Equal(t, IsLocal("random"), false)
}

func TestIsProd(t *testing.T) {
	assert.Equal(t, IsMainnet(MainnetNetwork), true)
	assert.Equal(t, IsMainnet(LocalNetwork), false)
	assert.Equal(t, IsMainnet("random"), false)
}
