package unionfind

// EnumeratedValues provides a compact representation of
// the roots in the unionfind structure while still being able to
// map the roots back to the original values (T)
type EnumeratedValues[T comparable] struct {
	ElementIndices  map[T]int
	IndexedElements []T
	lastIndex       int
}

// NewEnumeratedValues creates a new EnumeratedValues with the specified size
func NewEnumeratedValues[T comparable](size int) *EnumeratedValues[T] {
	return &EnumeratedValues[T]{
		ElementIndices:  make(map[T]int, size),
		IndexedElements: make([]T, 0, size),
		lastIndex:       -1,
	}
}

// FetchIndex gets or creates an index for the given element
func (ev *EnumeratedValues[T]) FetchIndex(element T) int {
	if idx, ok := ev.ElementIndices[element]; ok {
		return idx
	}

	ev.lastIndex++
	ev.ElementIndices[element] = ev.lastIndex
	ev.IndexedElements = append(ev.IndexedElements, element)
	return ev.lastIndex
}

// At returns the element at the given index
func (ev *EnumeratedValues[T]) At(index int) T {
	return ev.IndexedElements[index]
}
