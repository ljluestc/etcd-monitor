package algorithms

import (
	"slices"
	"testing"
)

// FuzzSortAlgorithmsEquivalence ensures multiple sort implementations agree.
func FuzzSortAlgorithmsEquivalence(f *testing.F) {
	seed := [][]int{
		{},
		{1},
		{2,1},
		{3,1,2,3,2,1},
		{-5,0,5,-1,2},
	}
	for _, s := range seed { f.Add(s) }

	f.Fuzz(func(t *testing.T, in []int) {
		// Guard against excessively large input to keep fuzz fast
		if len(in) > 2000 { t.Skip() }

		// Work on copies
		a := make([]int, len(in)); copy(a, in)
		b := make([]int, len(in)); copy(b, in)
		c := make([]int, len(in)); copy(c, in)
		d := make([]int, len(in)); copy(d, in)
		e := make([]int, len(in)); copy(e, in)

		// Sort using various algorithms
		SortIntsBubble(a)
		SortIntsSelection(b)
		SortIntsInsertion(c)
		SortIntsQuick(d)
		m := SortIntsMerge(e)

		// Use standard library as reference
		ref := make([]int, len(in)); copy(ref, in)
		slices.Sort(ref)

		if !slices.Equal(a, ref) { t.Fatalf("bubble mismatch") }
		if !slices.Equal(b, ref) { t.Fatalf("selection mismatch") }
		if !slices.Equal(c, ref) { t.Fatalf("insertion mismatch") }
		if !slices.Equal(d, ref) { t.Fatalf("quick mismatch") }
		if !slices.Equal(m, ref) { t.Fatalf("merge mismatch") }
	})
}

// FuzzBinarySearchVsLinear validates binary search against linear search.
func FuzzBinarySearchVsLinear(f *testing.F) {
	seed := [][]int{{}, {1}, {-1,0,1,2,3}, {1,1,2,2,3}}
	for _, s := range seed { f.Add(s) }

	f.Fuzz(func(t *testing.T, in []int) {
		if len(in) > 2000 { t.Skip() }
		// Work on sorted copy
		s := make([]int, len(in)); copy(s, in)
		slices.Sort(s)
		for _, target := range []int{-1000, -1, 0, 1, 2, 1000} {
			li := IndexOfLinearInts(s, target)
			bi := BinarySearchInts(s, target)
			if li == -1 {
				if bi != -1 { t.Fatalf("linear=-1 binary=%d on %v", bi, s) }
			} else {
				if bi < 0 || bi >= len(s) || s[bi] != target { t.Fatalf("binary wrong index for %d", target) }
			}
		}
	})
}


