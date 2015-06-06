package subsort

import (
	"fmt"
	"sort"
)

type sel1 struct {
	Key string
	Cnt int
}

type SortByCnt []sel1

func (s SortByCnt) Len() int {
	return len(s)
}

func (s SortByCnt) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortByCnt) Less(i, j int) bool {
	if s[i].Cnt > s[j].Cnt {
		return true
	}
	return false
}

// map[string]int is converted into a slice of
// []{Key,Val} and sorted by
func SortMapByCount(m map[string]int) SortByCnt {

	sT := make([]sel1, 0, len(m))

	for key, cnt := range m {
		sT = append(sT, sel1{key, cnt})
	}

	sbc := SortByCnt(sT)
	sort.Sort(sbc)
	return sbc
}

func (sbc SortByCnt) Print() {
	cntr := 0
	for k, val := range sbc {
		_ = k
		// fmt.Printf("%2v: %14v %4v ", k, val.Key, val.Cnt)
		fmt.Printf("%14v %4v ", val.Key, val.Cnt)
		cntr++
		if cntr%6 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}
