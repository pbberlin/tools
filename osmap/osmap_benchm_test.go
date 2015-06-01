package osmap

// go test -bench=. -benchtime=6ms

import (
	"math/rand"
	"testing"
)

var keys2 []string

func init() {

	rand.Seed(12345)
	keys2 = make([]string, 50*1000*1000)

	for i := 0; i < len(keys2); i++ {
		sz := 4 + rand.Intn(5)
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

func BenchmarkMapSet(b *testing.B) {
	m := make(Map2)
	for i := 0; i < b.N; i++ {
		offs := i % len(keys2)
		m.Set(keys2[offs], keys2[offs])
	}
}

func BenchmarkOSMSet(b *testing.B) {
	t := New()
	for i := 0; i < b.N; i++ {
		offs := i % len(keys2)
		t.Set(keys2[offs], keys2[offs])
	}
}

func BenchmarkMapGet(b *testing.B) {
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

func BenchmarkOSMGet(b *testing.B) {
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
