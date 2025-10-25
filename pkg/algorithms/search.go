package algorithms

// IndexOfLinearInts performs a linear search on a slice of ints.
// Returns the index of the first occurrence, or -1 if not found.
func IndexOfLinearInts(values []int, target int) int {
	for i := 0; i < len(values); i++ {
		if values[i] == target {
			return i
		}
	}
	return -1
}

// IndexOfLinear returns the index of target using linear scan for comparable types.
func IndexOfLinear[T comparable](values []T, target T) int {
	for i := 0; i < len(values); i++ {
		if values[i] == target {
			return i
		}
	}
	return -1
}

// BinarySearchInts performs binary search on a sorted slice of ints.
// Returns the index of any matching occurrence, or -1 if not found.
func BinarySearchInts(values []int, target int) int {
	lo, hi := 0, len(values)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		v := values[mid]
		if v == target {
			return mid
		}
		if v < target {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return -1
}

// BinarySearchFirstInts performs binary search on a sorted slice of ints and
// returns the index of the first occurrence of target, or -1 if not found.
func BinarySearchFirstInts(values []int, target int) int {
	lo, hi := 0, len(values)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if values[mid] < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	if lo < len(values) && values[lo] == target {
		return lo
	}
	return -1
}

// BinarySearchBy performs binary search over values using a comparator pair.
// equal(a,b) must agree with less(a,b) ordering.
func BinarySearchBy[T any](values []T, target T, less func(a, b T) bool, equal func(a, b T) bool) int {
	lo, hi := 0, len(values)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		mv := values[mid]
		if equal(mv, target) {
			return mid
		}
		if less(mv, target) {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return -1
}


