package unionfind_test

import (
	"testing"

	"github.com/maxjustus/bpuf/unionfind"

	"github.com/stretchr/testify/assert"
)

func TestUnionFind(t *testing.T) {
	t.Parallel()

	t.Run("Union", func(t *testing.T) {
		uf := unionfind.NewUnionFind(0)
		root := uf.Union(1, 2)
		assert.Equal(t, 1, root)
		root = uf.Union(2, 3)
		assert.Equal(t, 1, root)
		root = uf.Union(3, 4)
		assert.Equal(t, 1, root)
		root = uf.Union(5, 6)
		assert.Equal(t, 5, root)
		root = uf.Union(6, 1)
		assert.Equal(t, 1, root)
		assert.Equal(t, 1, uf.Find(5),
			"5 should be in the same set as 1 after union of 6 and 1")
		root = uf.Union(700, 801)
		assert.Equal(t, 700, root)
		root = uf.Union(801, 1000)
		assert.Equal(t, 700, root)
		assert.Equal(t, 700, uf.Find(1000),
			"1000 should be in the same set as 700 after union of 801 and 1000")
	})
}
