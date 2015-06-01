package util

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPrimes(t *testing.T) {

	b := new(bytes.Buffer)

	primes := PrecomputePrimes(100)

	for i := 0; i < len(primes); i++ {
		if primes[i] {
			b.WriteString(fmt.Sprintf("%v ", i))
		}
	}
	// fmt.Printf(b.String() + "\n")
	want := "0 1 2 3 5 7 11 13 17 19 23 29 31 37 41 43 47 53 59 61 67 71 73 79 83 89 97 "

	if b.String() != want {
		t.Fatalf("\nwnt %q\ngot %q\n", want, b.String())
	}

}

func TestMaxAbs(t *testing.T) {
	cases := [][]int{
		[]int{4, 5, 4},
		[]int{6, 5, 5},
		[]int{-4, -5, -4},
		[]int{-6, -5, -5},
	}
	for i := 0; i < len(cases); i++ {
		s := cases[i]
		res := MaxAbs(s[0], s[1])
		if res != s[2] {
			t.Fatalf("\nwnt %q\ngot %q\n", s[2], res)
		} else {
			// fmt.Printf("AbsMin % 2v,% 2v => % 2v as desired\n", s[0], s[1], s[2])
		}
	}
}

func TestIntSqrt(t *testing.T) {
	inputs := []int{25, 36, 49, 49284, 110889, 19749136, 1000*1000 + 1}
	wants := []int{5, 6, 7, 222, 333, 4444, 1000}

	for i := 0; i < len(inputs); i++ {
		inp := inputs[i]
		wnt := wants[i]
		got := Sqrt(inp)
		if got != wnt {
			t.Fatalf("\ninp %v\nwnt %v\ngot %v\n", inp, wnt, got)
		} else {
			// fmt.Printf("sqrt inp %v => %v\n", inp, got)
		}
	}
}
