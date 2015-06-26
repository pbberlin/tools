package util

import (
	"fmt"
	"sort"
)

// DoubleSliceCap doubles the capacity of a slice
func DoubleSliceCap(ba []byte) (newBa []byte) {
	newBa = make([]byte, len(ba), 2*cap(ba))
	copy(newBa, ba)
	fmt.Printf("len: %v -  newcap: %v <br>\n", len(newBa), cap(newBa))
	return newBa
}

// https://golang.org/x/go/wiki/SliceTricks
func InsertAfter(s []int, idx int, newVal int) []int {
	if idx > len(s)-1 {
		panic("Cannot insert beyond existing length")
	}
	s = append(s, 0)
	copy(s[idx+2:], s[idx+1:])
	s[idx+1] = newVal
	return s
}

func Delete(s []int, idx int) []int {
	if idx > len(s)-1 {
		panic("Cannot delete beyond existing length")
	}
	copy(s[idx:], s[idx+1:])
	s = s[:len(s)-1]
	return s
}

//
// When having multiple URL params of the same name,
func StringSliceToMapKeys(s []string) map[string]bool {
	ret := map[string]bool{}
	for _, v := range s {
		ret[v] = true
	}
	return ret
}

// Intslice2Int condenses several integers into a single int.
// Typically to use it as a map key.
// Intslice2Int unifies permutations.
// []int{0,3,1} and []int{1,3,0} are sorted to []int{0,1,3} and mapped to 30100
// []int{4,3,1} is sorted to []int{1,3,4} and mapped to 40301
func Intslice2Int(arg []int) (rv int) {
	if len(arg) > 10 {
		panic("only for small sets!")
	}
	sorted := make([]int, len(arg))
	copy(sorted, arg)
	sort.Ints(sorted)
	fact := 1
	for i := 0; i < len(sorted); i++ {
		rv += sorted[i] * fact
		fact = fact * 100
	}
	return
}

func IntsEqualsInts(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}

	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
