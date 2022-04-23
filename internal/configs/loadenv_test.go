package configs

import (
	"os"
	"strings"
	"testing"

	"github.com/test-go/testify/assert"
)

func TestGetProjectRoot(t *testing.T) {
	t.Run("should return project root without override", func(t *testing.T) {
		projectRoot := GetProjectRoot()
		assert.Equal(t, strings.HasSuffix(projectRoot, "drip-backend"), true)
	})

	t.Run("should return project root with override", func(t *testing.T) {
		err := os.Setenv(string(PROJECT_ROOT_OVERRIDE), "docs")
		assert.NoError(t, err)
		defer func() {
			err := os.Unsetenv(string(PROJECT_ROOT_OVERRIDE))
			assert.NoError(t, err)
		}()
		projectRoot := GetProjectRoot()
		assert.Equal(t, strings.HasSuffix(projectRoot, "docs"), true)
	})

}
