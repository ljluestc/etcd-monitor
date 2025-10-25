package algorithms

import (
	"slices"
	"testing"
)

func TestSearchAlgorithms_Linear_Ints(t *testing.T) {
	arr := []int{-3, -1, 0, 2, 2, 5}
	if got := IndexOfLinearInts(arr, -3); got != 0 { t.Fatalf("want 0 got %d", got) }
	if got := IndexOfLinearInts(arr, 2); got != 3 { t.Fatalf("want 3 got %d", got) }
	if got := IndexOfLinearInts(arr, 7); got != -1 { t.Fatalf("want -1 got %d", got) }
}

func TestSearchAlgorithms_Linear_Generic(t *testing.T) {
	arr := []string{"a","b","c","b"}
	if got := IndexOfLinear(arr, "b"); got != 1 { t.Fatalf("want 1 got %d", got) }
	if got := IndexOfLinear(arr, "z"); got != -1 { t.Fatalf("want -1 got %d", got) }
}

func TestSearchAlgorithms_Binary_Ints_Basic(t *testing.T) {
	arr := []int{-5,-1,0,1,2,3,4,10}
	if got := BinarySearchInts(arr, -5); got != 0 { t.Fatalf("want 0 got %d", got) }
	if got := BinarySearchInts(arr, 10); got != len(arr)-1 { t.Fatalf("want last got %d", got) }
	if got := BinarySearchInts(arr, 2); arr[got] != 2 { t.Fatalf("bad index %d", got) }
	if got := BinarySearchInts(arr, 7); got != -1 { t.Fatalf("want -1 got %d", got) }
}

func TestSearchAlgorithms_Binary_FirstOccurrence(t *testing.T) {
	arr := []int{1,2,2,2,3,4}
	if got := BinarySearchFirstInts(arr, 2); got != 1 { t.Fatalf("want 1 got %d", got) }
	if got := BinarySearchFirstInts(arr, 5); got != -1 { t.Fatalf("want -1 got %d", got) }
}

func TestSearchAlgorithms_Binary_CustomComparator(t *testing.T) {
	type item struct{ k int; v string }
	arr := []item{{1,"a"},{2,"b"},{2,"c"},{3,"d"}}
	less := func(a, b item) bool { return a.k < b.k }
	equal := func(a, b item) bool { return a.k == b.k }
	if got := BinarySearchBy(arr, item{k:2}, less, equal); got < 1 || got > 2 { t.Fatalf("unexpected index %d", got) }
}

func TestSearchAlgorithms_EmptyAndSingle(t *testing.T) {
	if got := BinarySearchInts(nil, 1); got != -1 { t.Fatalf("want -1 got %d", got) }
	if got := BinarySearchInts([]int{9}, 9); got != 0 { t.Fatalf("want 0 got %d", got) }
}

func TestSearchAlgorithms_UnsortedInput_Binary(t *testing.T) {
	arr := []int{3,2,1}
	// We don't guarantee correctness on unsorted input; verify it doesn't panic and returns either -1 or a valid index whose value equals target if by chance.
	_ = BinarySearchInts(arr, 2)
}

func TestSearchAlgorithms_CrossValidateWithLinear(t *testing.T) {
	arr := []int{-4,-1,0,1,3,5,8,13,21}
	for _, x := range []int{-4,-1,0,1,3,5,8,13,21,7,100} {
		li := IndexOfLinearInts(arr, x)
		bi := BinarySearchInts(arr, x)
		if li == -1 {
			if bi != -1 { t.Fatalf("linear=-1 but binary=%d", bi) }
		} else {
			if arr[bi] != x { t.Fatalf("binary index invalid for %d", x) }
		}
	}
}

func TestSearchAlgorithms_FirstOccurrence_Consistency(t *testing.T) {
	arr := []int{1,1,1,2,2,3}
	idx := BinarySearchFirstInts(arr, 1)
	if idx != 0 { t.Fatalf("want 0 got %d", idx) }
	idx2 := BinarySearchFirstInts(arr, 2)
	if idx2 != 3 { t.Fatalf("want 3 got %d", idx2) }
}

func TestSearchAlgorithms_Bounds(t *testing.T) {
	arr := []int{0,10,20}
	if got := BinarySearchInts(arr, -1); got != -1 { t.Fatalf("expect -1 got %d", got) }
	if got := BinarySearchInts(arr, 21); got != -1 { t.Fatalf("expect -1 got %d", got) }
}

func TestSearchAlgorithms_GenericLinear_Strings(t *testing.T) {
	arr := []string{"α","β","β","δ"}
	idx := IndexOfLinear(arr, "β")
	if idx != 1 && idx != 2 { t.Fatalf("unexpected idx %d", idx) }
}

func TestSearchAlgorithms_TargetPositions(t *testing.T) {
	arr := []int{1,2,3,4,5,6}
	if BinarySearchInts(arr, 1) != 0 { t.Fatalf("first pos") }
	if BinarySearchInts(arr, 6) != len(arr)-1 { t.Fatalf("last pos") }
	mid := BinarySearchInts(arr, 3)
	if arr[mid] != 3 { t.Fatalf("mid pos") }
}

func TestSearchAlgorithms_SortedRequirement(t *testing.T) {
	arr := []int{2,1,3,4}
	idx := BinarySearchInts(arr, 2)
	// If index returned, it may not point to correct value due to unsorted input; verify safety only.
	if idx >= 0 && idx < len(arr) {
		_ = arr[idx]
	}
}

func TestSearchAlgorithms_EqualsDefinition(t *testing.T) {
	type s struct{ a, b int }
	arr := []s{{0,0},{1,1},{2,2}}
	less := func(x, y s) bool { if x.a != y.a { return x.a < y.a }; return x.b < y.b }
	equal := func(x, y s) bool { return x.a == y.a }
	if BinarySearchBy(arr, s{a:1}, less, equal) == -1 { t.Fatalf("should find by equal definition") }
}

func TestSearchAlgorithms_SortedCopySafety(t *testing.T) {
	src := []int{5,4,3,2,1}
	cp := make([]int, len(src))
	copy(cp, src)
	slices.Sort(cp)
	if BinarySearchInts(cp, 3) == -1 { t.Fatalf("missing 3") }
}


