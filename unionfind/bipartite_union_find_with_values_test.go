package unionfind_test

import (
	"math/rand"
	"testing"

	"github.com/maxjustus/bpuf/unionfind"

	"github.com/stretchr/testify/assert"
)

func TestBipartiteUnionFindWithValues(t *testing.T) {
	t.Parallel()

	t.Run("Union()/Find()", func(t *testing.T) {
		uf := unionfind.NewBipartiteUnionFindWithValues[string, int](0)
		root := uf.UnionReturningValue("A", 1)
		assert.Equal(t, 1, root)
		root = uf.UnionReturningValue("B", 1)
		assert.Equal(t, 1, root)
		root = uf.UnionReturningValue("B", 2)
		assert.Equal(t, 1, root)
		root = uf.UnionReturningValue("C", 2)
		assert.Equal(t, 1, root)
		root = uf.UnionReturningValue("C", 3)
		assert.Equal(t, 1, root)
		root = uf.UnionReturningValue("D", 4)
		assert.Equal(t, 4, root)
		root = uf.UnionReturningValue("D", 5)
		assert.Equal(t, 4, root)
		root = uf.UnionReturningValue("D", 1)
		assert.Equal(t, 1, root)
	})
}

func BenchmarkBipartiteUnionFindWithValues(b *testing.B) {
	b.Run("UnionReturningValue() 100k times", func(b *testing.B) {
		uf := unionfind.NewBipartiteUnionFindWithValues[int, int](100000)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for i := 0; i < 100000; i++ {
				UValue := i
				if rand.Intn(2) == 0 { //nolint:gosec // test code using weak random is acceptable
					UValue += rand.Intn(5) //nolint:gosec // test code using weak random is acceptable
				}

				VValue := rand.Intn(30000) //nolint:gosec // test code using weak random is acceptable

				// Produces about 6000 groups or 16 records per group
				uf.Union(UValue, VValue)
			}
		}
	})
}
