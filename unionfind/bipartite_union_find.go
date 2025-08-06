// Package unionfind implements union-find (disjoint set) data structures with path compression
// and union by rank optimization for efficient set operations.
package unionfind

// BipartiteUnionFind facilitates iterative construction of a graph via union find
// given only transitive edges between sets U and V in a bipartite graph.
// Works by caching the last found root for each U in V and then using the cached V for a given U
// as the root for the next union operation on the same U.
// Union(U1, V1), Union(U1, V2) results in Union(V1, V1) and Union(V1, V2) in underlying UnionFind.

// BipartiteUnionFind represents a union-find structure for bipartite graphs
type BipartiteUnionFind struct {
	*UnionFind
	lastRootForUInV            []int
	lastRootForUInVInitialized []bool
}

// NewBipartiteUnionFind creates a new BipartiteUnionFind with the specified capacity
func NewBipartiteUnionFind(capacity int) *BipartiteUnionFind {
	return &BipartiteUnionFind{
		UnionFind:                  NewUnionFind(capacity),
		lastRootForUInV:            make([]int, capacity),
		lastRootForUInVInitialized: make([]bool, capacity),
	}
}

// Union connects elements u and v, returning the root of the merged set
func (buf *BipartiteUnionFind) Union(u, v int) int {
	if len(buf.lastRootForUInV) <= u {
		buf.lastRootForUInV = expandSlice(buf.lastRootForUInV, u)
		buf.lastRootForUInVInitialized = expandSlice(buf.lastRootForUInVInitialized, u)
	}

	var uIdxForV int
	if !buf.lastRootForUInVInitialized[u] {
		uIdxForV = v
	} else {
		uIdxForV = buf.lastRootForUInV[u]
	}

	newRoot := buf.UnionFind.Union(uIdxForV, v)
	buf.lastRootForUInV[u] = newRoot
	buf.lastRootForUInVInitialized[u] = true
	return newRoot
}

// FindAssociatedRoot finds the root associated with element u in the V set
func (buf *BipartiteUnionFind) FindAssociatedRoot(u int) (int, bool) {
	if len(buf.lastRootForUInV) <= u {
		return -1, false
	}

	return buf.Find(buf.lastRootForUInV[u]), true
}
