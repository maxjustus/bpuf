package unionfind

import "slices"

func expandSlice[T any](s []T, newLen int) []T {
	newLen++
	if newLen > len(s) {
		return slices.Grow(s, newLen)[0:newLen]
	}

	return s
}
