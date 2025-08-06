package unionfind

// BipartiteUnionFindWithValues provides a bipartite union-find structure with generic values
type BipartiteUnionFindWithValues[U, V comparable] struct {
	*BipartiteUnionFind
	UValues *EnumeratedValues[U]
	VValues *EnumeratedValues[V]
}

// NewBipartiteUnionFindWithValues creates a new BipartiteUnionFindWithValues.
// A bipartite graph is a graph whose vertices can be divided into two disjoint sets U and V
// such that every edge connects a vertex in U to one in V.
// https://en.wikipedia.org/wiki/Bipartite_graph
func NewBipartiteUnionFindWithValues[U, V comparable](capacity int) *BipartiteUnionFindWithValues[U, V] {
	return &BipartiteUnionFindWithValues[U, V]{
		BipartiteUnionFind: NewBipartiteUnionFind(capacity),
		UValues:            NewEnumeratedValues[U](capacity),
		VValues:            NewEnumeratedValues[V](capacity),
	}
}

// FindVRootForU finds the V root for the given U element
func (buf *BipartiteUnionFindWithValues[U, V]) FindVRootForU(u U) (V, bool) {
	uIndex := buf.UValues.FetchIndex(u)
	return buf.FindVRootForUIndex(uIndex)
}

// FindVRootForUIndex finds the V root for the given U element index
func (buf *BipartiteUnionFindWithValues[U, V]) FindVRootForUIndex(uIndex int) (V, bool) {
	vIndex, ok := buf.FindAssociatedRoot(uIndex)
	if !ok {
		var zero V
		return zero, false
	}

	return buf.VValues.At(vIndex), true
}

// Union connects U and V elements, returning the root index
func (buf *BipartiteUnionFindWithValues[U, V]) Union(u U, v V) int {
	uIndex := buf.UValues.FetchIndex(u)
	vIndex := buf.VValues.FetchIndex(v)

	return buf.BipartiteUnionFind.Union(uIndex, vIndex)
}

// UnionReturningValue connects U and V elements and returns the root value.
// This is handy for directly getting new root value, but is slower
// because it looks up value by index.
func (buf *BipartiteUnionFindWithValues[U, V]) UnionReturningValue(u U, v V) V {
	idx := buf.Union(u, v)
	return buf.VValues.At(idx)
}

// FindReturningValue finds the root and returns it as a value
func (buf *BipartiteUnionFindWithValues[U, V]) FindReturningValue(v V) V {
	vIndex := buf.VValues.FetchIndex(v)
	return buf.VValues.At(buf.Find(vIndex))
}
