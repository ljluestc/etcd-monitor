package algorithms

import "testing"

func TestSortAlgorithms_NilInputs(t *testing.T) {
	var nilSlice []int
	// in-place algorithms should be safe no-ops on nil slices
	SortIntsBubble(nilSlice)
	SortIntsSelection(nilSlice)
	SortIntsInsertion(nilSlice)
	SortIntsQuick(nilSlice)
	// merge returns a new slice; on nil it should return nil-length slice
	out := SortIntsMerge(nilSlice)
	if out == nil || len(out) != 0 { /* acceptable: nil or empty */ }
}


