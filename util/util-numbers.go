package util

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Abs returns the absolute part of an integer
func Abs(a int) int {
	if a < 0 {
		a = -a
	}
	return a
}

// Itos converts integer to string
func Itos(i int) string {
	str := fmt.Sprintf("%v", i)
	return str
}

// Stoi converts string to integer
func Stoi(s string) int {
	s = strings.TrimSpace(s)
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
		log.Printf("could not convert: %v \n", err)
		return 0
	}
	return i
}

// Stof converts string to float64
func Stof(s string) float64 {
	s = strings.TrimSpace(s)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Printf("could not convert: %v \n", err)
		return 0.0
	}
	return f
}

// Round does the rounding - and converts float to integer
func Round(f float64) int {

	var f1 float64
	if f < 0 {
		f1 = math.Ceil(f - 0.5)
	} else {
		f1 = math.Floor(f + 0.5)
	}
	i := int(f1)
	//log.Printf("--%v %v %v",f,f1,i)
	return i
}

// Max returns the smaller of two integers
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Max returns the bigger of two integers
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// ManAbs() is challenging.
// The second param is a "bound" to the absolute of x
// I.e. when we want x in the range of -4...0...4
// then make x = MinAbs(x,4)
func MaxAbs(x, bound int) int {
	bound = Abs(bound)
	if x > 0 {
		if x > bound {
			return bound
		}
		return x
	} else {
		if x < -bound {
			return -bound
		}
		return x
	}
}

// ComputePrimes precomputes a slice, holding prime number quality.
// The returned slice can be globally stored and effectively consulted:
//    if primes[testNumber] {...}
// Warning: Some *large* numbers listed as prime,
// might have undetected factors.
// You may increase testQuality, if you need better quality results.
func PrecomputePrimes(howMany int) []bool {
	const testQuality = 1
	howMany64 := int64(howMany)
	ret := make([]bool, howMany-1)
	for i := int64(0); i < howMany64; i++ {
		bigI := big.NewInt(i)                      // int64 and big
		isPrime := bigI.ProbablyPrime(testQuality) // 1 enough for i < 200
		if isPrime {
			ret[i] = true // takes int64 as index, astounding
		}
	}
	if howMany > 1 { // 0 and 1 are actually undefined - but for our purposes they behave prime
		ret[0] = true
		ret[1] = true
	}
	return ret

}

// Sqrt is an iterative approximation
func Sqrt(n int) int {
	t := uint(n)
	p := uint(1 << 30) // 2^30
	for p > t {
		p >>= 2 // divide by 4, make p 2-powered ceiling of t
	}

	var b, r uint
	for ; p != 0; p >>= 2 {
		b = r | p
		r >>= 1
		if t >= b {
			t -= b
			r |= p
		}
	}
	return int(r)
}
