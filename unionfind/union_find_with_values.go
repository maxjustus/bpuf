package unionfind

// AlgoUnionFindWithValues represents a union-find structure with generic values
type AlgoUnionFindWithValues[T comparable] struct {
	*UnionFind
	values *EnumeratedValues[T]
}

// NewUnionFindWithValues creates a new AlgoUnionFindWithValues with the specified capacity
func NewUnionFindWithValues[T comparable](capacity int) *AlgoUnionFindWithValues[T] {
	return &AlgoUnionFindWithValues[T]{
		UnionFind: NewUnionFind(capacity),
		values:    NewEnumeratedValues[T](capacity),
	}
}

// Find returns the root index of the set containing the given value
func (uf *AlgoUnionFindWithValues[T]) Find(value T) int {
	return uf.UnionFind.Find(uf.values.FetchIndex(value))
}

// FindReturningValue returns the root value of the set containing the given value
func (uf *AlgoUnionFindWithValues[T]) FindReturningValue(value T) T {
	return uf.values.At(uf.Find(value))
}

// Union merges the sets containing values a and b, returning the root index
func (uf *AlgoUnionFindWithValues[T]) Union(a, b T) int {
	indexA := uf.values.FetchIndex(a)
	indexB := uf.values.FetchIndex(b)

	return uf.UnionFind.Union(indexA, indexB)
}

// UnionReturningValue merges the sets containing values a and b, returning the root value.
// This is handy for directly getting new root value, but is slower
// because it looks up value by index.
func (uf *AlgoUnionFindWithValues[T]) UnionReturningValue(a, b T) T {
	idx := uf.Union(a, b)
	return uf.values.At(idx)
}
