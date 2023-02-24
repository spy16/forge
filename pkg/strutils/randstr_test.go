package strutils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spy16/forge/pkg/strutils"
)

func TestRandStr(t *testing.T) {
	t.Parallel()

	t.Run("Length", func(t *testing.T) {
		assert.Len(t, strutils.RandStr(10), 10)
		assert.Len(t, strutils.RandStr(10, "a"), 10)
	})

	t.Run("Charset", func(t *testing.T) {
		val := strutils.RandStr(10, "a")
		assert.Equal(t, "aaaaaaaaaa", val)
	})
}
