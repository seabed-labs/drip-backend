package main

import (
	"os"
	"testing"

	"github.com/test-go/testify/assert"
	"go.uber.org/fx"
)

func TestDependencies(t *testing.T) {
	deps := getDependencies()
	err := fx.ValidateApp(deps...)
	assert.NoError(t, err)

	err = os.Setenv("ENV", "PROD")
	assert.NoError(t, err)
	deps = getDependencies()
	err = fx.ValidateApp(deps...)
	assert.NoError(t, err)
}
