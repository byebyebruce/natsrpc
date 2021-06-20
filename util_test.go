package natsrpc

import "testing"

func Benchmark_CombineSubsets(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CombineSubject("a", "b", "c")
	}
}
