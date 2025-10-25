package algorithms

import "sort"

// SortIntsBubble performs an in-place bubble sort on the provided slice.
func SortIntsBubble(values []int) {
	if len(values) < 2 {
		return
	}
	swapped := true
	for swapped {
		swapped = false
		for i := 1; i < len(values); i++ {
			if values[i-1] > values[i] {
				values[i-1], values[i] = values[i], values[i-1]
				swapped = true
			}
		}
	}
}

// SortIntsSelection performs an in-place selection sort.
func SortIntsSelection(values []int) {
	for i := 0; i < len(values); i++ {
		minIdx := i
		for j := i + 1; j < len(values); j++ {
			if values[j] < values[minIdx] {
				minIdx = j
			}
		}
		values[i], values[minIdx] = values[minIdx], values[i]
	}
}

// SortIntsInsertion performs an in-place insertion sort.
func SortIntsInsertion(values []int) {
	for i := 1; i < len(values); i++ {
		key := values[i]
		j := i - 1
		for j >= 0 && values[j] > key {
			values[j+1] = values[j]
			j--
		}
		values[j+1] = key
	}
}

// SortIntsMerge returns a new sorted slice using merge sort (stable).
func SortIntsMerge(values []int) []int {
	if len(values) < 2 {
		out := make([]int, len(values))
		copy(out, values)
		return out
	}
	mid := len(values) / 2
	left := SortIntsMerge(values[:mid])
	right := SortIntsMerge(values[mid:])
	return merge(left, right)
}

func merge(a, b []int) []int {
	result := make([]int, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}
	result = append(result, a[i:]...)
	result = append(result, b[j:]...)
	return result
}

// SortIntsQuick performs an in-place quicksort (not stable).
func SortIntsQuick(values []int) {
	quickSort(values, 0, len(values)-1)
}

func quickSort(values []int, low, high int) {
	if low >= high {
		return
	}
	p := partition(values, low, high)
	quickSort(values, low, p-1)
	quickSort(values, p+1, high)
}

func partition(values []int, low, high int) int {
	pivot := values[high]
	i := low
	for j := low; j < high; j++ {
		if values[j] <= pivot {
			values[i], values[j] = values[j], values[i]
			i++
		}
	}
	values[i], values[high] = values[high], values[i]
	return i
}

// SortStableWithLess sorts any slice using a less comparator and guarantees stability.
// The slice must be addressable by the standard library sort.SliceStable.
func SortStableWithLess[T any](values []T, less func(i, j int) bool) {
	sort.SliceStable(values, less)
}


