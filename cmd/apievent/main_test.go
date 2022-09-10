package main

import (
	"testing"

	"github.com/test-go/testify/assert"
	"go.uber.org/fx"
)

func TestDependencies(t *testing.T) {
	deps := getDependencies()
	err := fx.ValidateApp(deps...)
	assert.NoError(t, err)
}
