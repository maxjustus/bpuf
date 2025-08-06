package unionfind_test

import (
	"math/rand"
	"testing"

	"github.com/maxjustus/bpuf/unionfind"

	"github.com/stretchr/testify/assert"
)

func TestBipartiteUnionFind(t *testing.T) {
	t.Parallel()

	t.Run("Union/Find", func(t *testing.T) {
		t.Run("unions and finds associated roots", func(t *testing.T) {
			uf := unionfind.NewBipartiteUnionFind(0)
			root := uf.Union(1, 2)
			assert.Equal(t, 2, root)
			root = uf.Union(1, 3)
			assert.Equal(t, 2, root)
			root = uf.Union(2, 3)
			assert.Equal(t, 2, root)
			root = uf.Union(2, 6)
			assert.Equal(t, 2, root)
			root = uf.Union(3, 6)
			assert.Equal(t, 2, root)
			root = uf.Union(6, 7)
			assert.Equal(t, 7, root)
			root = uf.Union(6, 8)
			assert.Equal(t, 7, root)
			root = uf.Union(7, 8)
			assert.Equal(t, 7, root)
			root = uf.Union(7, 2)
			assert.Equal(t, 2, root)
			root, ok := uf.FindAssociatedRoot(6)
			assert.True(t, ok)
			assert.Equal(t, 2, root, "6 should be in the same set as 2 after union of 7 and 2")
		})
	})
}

func BenchmarkBipartiteUnionFind(b *testing.B) {
	b.Run("Union()", func(b *testing.B) {
		uf := unionfind.NewBipartiteUnionFind(100000)
		for i := 0; i < b.N; i++ {
			uf.Union(i, i+rand.Intn(1000)) //nolint:gosec // test code using weak random is acceptable
		}
	})
}
