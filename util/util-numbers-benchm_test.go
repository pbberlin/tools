package util

// go test -bench=. -benchtime=6ms

import (
	"math"
	"testing"
)

// var inputsTestSqrtr = []int{25, 36, 49}

var inputsTestSqrtr = []int{49284, 110889, 19749136, 1000*1000 + 1}

func BenchmarkMath(b *testing.B) {

	x := 0
	for i := 0; i < b.N; i++ {
		idx := b.N % len(inputsTestSqrtr)
		x = int(math.Sqrt(float64(inputsTestSqrtr[idx])))
	}
	_ = x
}

func BenchmarkApprox(b *testing.B) {
	x := 0
	for i := 0; i < b.N; i++ {
		idx := b.N % len(inputsTestSqrtr)
		x = Sqrt(inputsTestSqrtr[idx])
	}
	_ = x
}
