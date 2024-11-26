package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//	"fmt"
)

func TestCpuLinuxTicks(t *testing.T) {
	c, err := GetCpuTicks()
	assert.NoError(t, err)
	t.Run("Total", func(t *testing.T) {
		assert.Greater(t, c[0].total, uint64(0))
		assert.Greater(t, c[0].user, uint64(0))
		assert.Greater(t, c[0].system, uint64(0))
	})
	t.Run("First CPU", func(t *testing.T) {
		assert.Greater(t, c[1].total, uint64(0))
		assert.Greater(t, c[1].user, uint64(0))
		assert.Greater(t, c[1].system, uint64(0))
	})
}
