package sortmapx

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {

	m := NewSortedMapInt2Int()
	m.Set(1, 2)
	m.Set(1, 3)
	m.Set(4, -3)
	m.Set(-4, -13)
	// fmt.Printf("%v \n", m)

	{
		x := m.SortedKV()
		// fmt.Printf("%v \n", x)
		wnt := `[{-4 -13} {1 3} {4 -3}]`
		got := fmt.Sprintf("%v", x)
		if wnt != got {
			t.Errorf("want != got\n%q\n%q\n", wnt, got)
		}
	}

	{
		x, _ := m.SortedPartialSums()
		// fmt.Printf("%v \n", x)
		wnt := `[{-4 -13} {-3 -13} {-2 -13} {-1 -13} {0 -13} {1 -10} {2 -10} {3 -10} {4 -13}]`
		got := fmt.Sprintf("%v", x)
		if wnt != got {
			t.Errorf("want != got\n%q\n%q\n", wnt, got)
		}
	}

}
