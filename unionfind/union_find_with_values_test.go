package unionfind_test

import (
	"testing"

	"github.com/maxjustus/bpuf/unionfind"

	"github.com/stretchr/testify/assert"
)

func TestUnionFindWithValues(t *testing.T) {
	t.Parallel()

	t.Run("Union/Find", func(t *testing.T) {
		uf := unionfind.NewUnionFindWithValues[string](0)
		root := uf.UnionReturningValue("A", "B")
		assert.Equal(t, "A", root)
		root = uf.UnionReturningValue("B", "C")
		assert.Equal(t, "A", root)
	})
}
