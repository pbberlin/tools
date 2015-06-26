// +build benchm
// go test -tags=benchm  -bench=. -benchtime=6ms

package osmaps

import (
	"fmt"
	"math/rand"
	"testing"
)

var keys2 []string

const cval = "v"

func init() {

	rand.Seed(12345 + 2)
	keys2 = make([]string, 2*1000*1000)
	keys2 = make([]string, 22*1000)

	for i := 0; i < len(keys2); i++ {
		sz := 5 + rand.Intn(10)
		b := make([]byte, sz)
		for j := 0; j < sz; j++ {
			b[j] = byte(97 + rand.Intn(24))
		}
		keys2[i] = string(b)
	}
}

type Map2 map[string]string

func (m Map2) Set(key string, value string) {
	m[key] = value
}
func (m Map2) Get(key string) string {
	return m[key]
}

func DISABLED_BenchmarkMapSet(b *testing.B) {
	m := make(Map2)
	for i := 0; i < b.N; i++ {
		offs := i % len(keys2)
		m.Set(keys2[offs], cval)
	}
}

func BenchmarkOSMSet(b *testing.B) {
	t := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		offs := i % len(keys2)
		t.Set(keys2[offs], cval)
	}

	//
	fmt.Printf("\n\n Histo:\n")
	sum := 0
	popul := 0
	alreadyDone := make(map[*node]bool, 0)
	for i := len(ndStat) - 1; i >= 0; i-- {
		mpNodes := ndStat[i]
		if len(mpNodes) == 0 {
			continue
		}
		fmt.Printf("%2v: %2v\t", i, len(mpNodes))
		for nd1, _ := range mpNodes {
			if alreadyDone[nd1] {
				// fmt.Printf("skipp ")
				continue
			}
			sum++
			fmt.Printf("%2v ", len(nd1.kv))
			popul += len(nd1.kv)
			alreadyDone[nd1] = true
		}
		fmt.Printf("\n")
	}
	fmt.Printf("%2v nodes total - %v\n", sum, popul)

}

func DISABLED_BenchmarkMapGet(b *testing.B) {
	m := make(Map2)
	for i := 0; i < b.N; i++ {
		offs := i % len(keys2)
		m.Set(keys2[offs], keys2[offs])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(keys2[i%len(keys2)])
	}
}

func DISABLED_BenchmarkOSMGet(b *testing.B) {
	t := New()
	for i := 0; i < b.N; i++ {
		offs := i % len(keys2)
		t.Set(keys2[offs], keys2[offs])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Get(keys2[i%len(keys2)])
	}
}
