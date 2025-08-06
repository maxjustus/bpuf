package unionfind

// https://en.wikipedia.org/wiki/Disjoint-set_data_structure

// UnionFind represents a union-find (disjoint set) data structure
type UnionFind struct {
	Root        []int  // Parent of each element by index
	Initialized []bool // Whether the element has been initialized
	RootCount   int    // Number of roots
	// Upper bound cardinality of each set by root index.
	// this is used to weight the union operation
	// to keep the tree as flat as possible by
	// preferring to make the smaller tree a child
	// of the larger tree in union operations.
	// Not expected to be an exact representation of
	// cardinality as union operations are performed.
	// see https://stackoverflow.com/a/69063833
	Rank []int
}

// NewUnionFind creates a new UnionFind with the specified capacity
func NewUnionFind(capacity int) *UnionFind {
	return &UnionFind{
		Root:        make([]int, capacity),
		Initialized: make([]bool, capacity),
		Rank:        make([]int, capacity),
	}
}

func (uf *UnionFind) addElement(n int) {
	if n >= len(uf.Root) {
		uf.Root = expandSlice(uf.Root, n)
		uf.Initialized = expandSlice(uf.Initialized, n)
		uf.Rank = expandSlice(uf.Rank, n)
	}

	if !uf.Initialized[n] {
		uf.Root[n] = n
		uf.Initialized[n] = true
		uf.Rank[n] = 1
		uf.RootCount++
	}
}

// Find returns the root of the set containing the given index
func (uf *UnionFind) Find(index int) int {
	uf.addElement(index)

	for uf.Root[index] != index {
		uf.Root[index] = uf.Root[uf.Root[index]] // Path compression
		index = uf.Root[index]
	}

	return index
}

// Union merges the sets containing a and b, returning the root of the merged set
func (uf *UnionFind) Union(a, b int) int {
	rootA := uf.Find(a)
	rootB := uf.Find(b)

	if rootA != rootB {
		if uf.Rank[rootA] < uf.Rank[rootB] {
			uf.Root[rootA] = rootB
			uf.Rank[rootB] += uf.Rank[rootA]

			return rootB
		}

		uf.Root[rootB] = rootA
		uf.Rank[rootA] += uf.Rank[rootB]

		return rootA
	}

	return rootA
}
