package algorithms

import (
	"math/rand"
	"slices"
	"testing"
)

func isNonDecreasing(values []int) bool {
	for i := 1; i < len(values); i++ {
		if values[i-1] > values[i] {
			return false
		}
	}
	return true
}

func makeSliceCopy(v []int) []int {
	cp := make([]int, len(v))
	copy(cp, v)
	return cp
}

func TestSortAlgorithms_BasicCases(t *testing.T) {
	inputs := [][]int{
		{},
		{1},
		{1, 2, 3, 4, 5},
		{5, 4, 3, 2, 1},
		{2, 1, 2, 1, 2},
		{-3, -1, -2, 0, 2, 1},
	}

	algos := []struct {
		name string
		fn   func([]int)
	}{
		{"bubble", SortIntsBubble},
		{"selection", SortIntsSelection},
		{"insertion", SortIntsInsertion},
		{"quick", SortIntsQuick},
	}

	for _, in := range inputs {
		for _, a := range algos {
			v := makeSliceCopy(in)
			a.fn(v)
			if !isNonDecreasing(v) {
				t.Fatalf("%s failed to sort %v => %v", a.name, in, v)
			}
		}

		v := SortIntsMerge(in)
		if !isNonDecreasing(v) {
			t.Fatalf("merge failed to sort %v => %v", in, v)
		}
	}
}

func TestSortAlgorithms_Idempotence(t *testing.T) {
	in := []int{5, 1, 4, 2, 8, 5}
	algos := []struct{ name string; fn func([]int) }{
		{"bubble", SortIntsBubble},
		{"selection", SortIntsSelection},
		{"insertion", SortIntsInsertion},
		{"quick", SortIntsQuick},
	}
	for _, a := range algos {
		v := makeSliceCopy(in)
		a.fn(v)
		a.fn(v)
		if !isNonDecreasing(v) {
			t.Fatalf("%s not idempotent: %v", a.name, v)
		}
	}

	v := SortIntsMerge(in)
	v2 := SortIntsMerge(v)
	if !slices.Equal(v, v2) {
		t.Fatalf("merge not idempotent: %v vs %v", v, v2)
	}
}

func TestSortAlgorithms_PermutationPreserved(t *testing.T) {
	in := []int{3, 1, 2, 3, 2, 1}
	orig := make(map[int]int)
	for _, x := range in {
		orig[x]++
	}
	check := func(sorted []int) {
		m := make(map[int]int)
		for _, x := range sorted { m[x]++ }
		if len(m) != len(orig) {
			t.Fatalf("element set changed: %v -> %v", orig, m)
		}
		for k, c := range orig {
			if m[k] != c { t.Fatalf("count mismatch for %d: %d vs %d", k, m[k], c) }
		}
	}

	v1 := makeSliceCopy(in); SortIntsBubble(v1); check(v1)
	v2 := makeSliceCopy(in); SortIntsSelection(v2); check(v2)
	v3 := makeSliceCopy(in); SortIntsInsertion(v3); check(v3)
	v4 := SortIntsMerge(in); check(v4)
	v5 := makeSliceCopy(in); SortIntsQuick(v5); check(v5)
}

func TestSortAlgorithms_Stability_Merge(t *testing.T) {
	type pair struct{ key, idx int }
	in := []pair{{2,0},{1,0},{2,1},{1,1},{2,2}}
	out := make([]pair, len(in))
	copy(out, in)
	SortStableWithLess(out, func(i, j int) bool { return out[i].key < out[j].key })

	// For equal keys, relative order of idx must be preserved
	lastIdxForKey := map[int]int{}
	for _, p := range out {
		if prev, ok := lastIdxForKey[p.key]; ok && p.idx < prev {
			t.Fatalf("stability violated for key=%d: %v", p.key, out)
		}
		lastIdxForKey[p.key] = p.idx
	}
}

func TestSortAlgorithms_Randomized(t *testing.T) {
	rand.Seed(1)
	for n := 1; n <= 50; n += 7 {
		arr := make([]int, n)
		for i := 0; i < n; i++ { arr[i] = rand.Intn(100) - 50 }
		v1 := makeSliceCopy(arr); SortIntsInsertion(v1)
		v2 := SortIntsMerge(arr)
		if !slices.Equal(v1, v2) {
			t.Fatalf("insertion vs merge mismatch\nins=%v\nmer=%v", v1, v2)
		}
	}
}

func TestSortAlgorithms_NearlySorted(t *testing.T) {
	in := []int{1,2,3,4,5,6,7,8,9}
	in[3], in[4] = in[4], in[3]
	v := makeSliceCopy(in)
	SortIntsInsertion(v)
	if !isNonDecreasing(v) { t.Fatalf("not sorted: %v", v) }
}

func TestSortAlgorithms_AllEqual(t *testing.T) {
	in := make([]int, 100)
	v := makeSliceCopy(in)
	SortIntsQuick(v)
	if !isNonDecreasing(v) { t.Fatalf("not sorted: %v", v) }
}

func TestSortAlgorithms_LargeSanity(t *testing.T) {
	rand.Seed(2)
	n := 500
	in := make([]int, n)
	for i := 0; i < n; i++ { in[i] = rand.Intn(10000) - 5000 }
	v := SortIntsMerge(in)
	if !isNonDecreasing(v) { t.Fatalf("merge not sorted") }
}


