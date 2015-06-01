package uruntime

import (
	"fmt"
	"math/rand"
	"testing"
	"unsafe"
)

type testStruct struct {
	// i1, i2 int
}

var preventCompilerOptimization int

func TestMemoryConsumptionOfMap(t *testing.T) {

	fLogger, fDumper := AllocLogger()
	lenPrev := 0

	var m map[int]*testStruct
	m = make(map[int]*testStruct, 10*1000)

	for i := 0; i < 100*1000; i++ {

		key := rand.Intn(i + 1)
		ts := testStruct{}
		m[key] = &ts
		if i%10000 == 0 {
			fmt.Printf("%11v | ", len(m)-lenPrev)
			lenPrev = len(m)
			fLogger()
		}
	}
	fmt.Printf("\n")
	// _ = m
	fDumper()

}

func TestMapBooleanVsEmpty(t *testing.T) {

	var bl bool = true
	var empty struct{}

	const mapSize = 10 * 1000
	const numElements = 10 * 1000

	fmt.Println("boolean consumes ", unsafe.Sizeof(bl), " byte(s)")
	fmt.Println("empty   consumes ", unsafe.Sizeof(empty), " byte(s)")

	fLogger, fDumper := AllocLogger()

	m1 := make(map[int]struct{}, mapSize)
	fLogger()

	m2 := make(map[int]bool, mapSize)
	fLogger()

	// fmt.Println(unsafe.Sizeof(m1))
	// fmt.Println(unsafe.Sizeof(m2))

	cntr := 0
	for i := 0; i < numElements; i++ {
		m1[i] = empty
		if i%(numElements/5) == numElements/5-1 {
			fLogger()
		}
		cntr++
		preventCompilerOptimization = cntr
	}
	fLogger()

	for i := 0; i < numElements; i++ {
		m2[i] = true
		if i%(numElements/5) == numElements/5-1 {
			fLogger()
		}
		cntr++
		preventCompilerOptimization = cntr
	}
	fDumper()

}
