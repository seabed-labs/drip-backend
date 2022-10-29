package configs

import (
	"testing"

	"github.com/test-go/testify/assert"
)

func TestNewAppConfig(t *testing.T) {
	config, err := NewAppConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Contains(t, []Network{NilNetwork, LocalNetwork, DevnetNetwork}, config.GetNetwork())
}

func TestNewPSQLConfig(t *testing.T) {
	config, err := NewPSQLConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestIsDev(t *testing.T) {
	assert.Equal(t, IsDevnetNetwork(DevnetNetwork), true)
	assert.Equal(t, IsDevnetNetwork(LocalNetwork), false)
	assert.Equal(t, IsDevnetNetwork("random"), false)
}

func TestIsLocal(t *testing.T) {
	assert.Equal(t, IsLocalNetwork(NilNetwork), true)
	assert.Equal(t, IsLocalNetwork(LocalNetwork), true)
	assert.Equal(t, IsLocalNetwork(DevnetNetwork), false)
	assert.Equal(t, IsLocalNetwork("random"), false)
}

func TestIsProd(t *testing.T) {
	assert.Equal(t, IsMainnetNetwork(MainnetNetwork), true)
	assert.Equal(t, IsMainnetNetwork(LocalNetwork), false)
	assert.Equal(t, IsMainnetNetwork("random"), false)
}
