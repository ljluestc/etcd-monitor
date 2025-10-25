package algorithms

import (
	"math/rand"
	"testing"
)

func BenchmarkSortMerge_1e3(b *testing.B) {
	benchmarkSortMerge(b, 1000)
}

func BenchmarkSortMerge_1e4(b *testing.B) {
	benchmarkSortMerge(b, 10000)
}

func benchmarkSortMerge(b *testing.B, n int) {
	b.ReportAllocs()
	arr := make([]int, n)
	for i := 0; i < n; i++ { arr[i] = rand.Int() }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SortIntsMerge(arr)
	}
}


