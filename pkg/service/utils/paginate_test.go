package utils

import (
	"testing"

	"github.com/test-go/testify/assert"
)

func TestDoForPaginatedBatch(t *testing.T) {
	t.Run("should not paginate 0 times", func(t *testing.T) {
		counter := 0
		assert.NoError(t, DoForPaginatedBatch(0, 0, func(start, end int) {
			counter++
		}))
		assert.Equal(t, 0, counter)
	})

	t.Run("should paginate 1 time", func(t *testing.T) {
		counter := 0
		assert.NoError(t, DoForPaginatedBatch(1, 1, func(start, end int) {
			counter++
		}))
		assert.Equal(t, 1, counter)
	})

	t.Run("should paginate 3 times", func(t *testing.T) {
		counter := 0
		assert.NoError(t, DoForPaginatedBatch(1, 3, func(start, end int) {
			counter++
		}))
		assert.Equal(t, 3, counter)
	})
}
