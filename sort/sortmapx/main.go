// package sortmapx sorts maps into
// an ordered *persistent* slice of KV-structs;
// it also constructs slices of dense key ranges
// with partial sums of values; see tests;
// before using, consider the simpler sortmap.StringKeysToSortedArray()
// or sortmap.StringKeysToSortedArray();
// before using, consider the simpler sortmap.SortByCnt(map[string]int).
// It is useful for with unpredictable *sparse* long tails.
// This package only adds value, if we want *repeated* access to
// specific sort offsets, or to min and max values.
// The package furthermore *flattens* the map to a fully ranged array
// of offsets.
package sortmapx

import "sort"

type sortedKV struct {
	K, V int
}

type SortedMapInt2Int struct {
	m          map[int]int
	dirty      bool
	sortedKeys []int      // statified partially: [sortIndex][key]
	sKV        []sortedKV // statified fully:     [sortIndex][key][val]
}

func (m *SortedMapInt2Int) updateSort() {
	keys := make([]int, len(m.m))
	i := 0
	for key, _ := range m.m {
		keys[i] = key
		i++
	}
	sort.Ints(keys)
	m.sortedKeys = keys
}

func NewSortedMapInt2Int() (sm SortedMapInt2Int) {
	sm = SortedMapInt2Int{}
	sm.m = make(map[int]int)
	sm.sortedKeys = make([]int, 0)
	return
}

func (m *SortedMapInt2Int) Set(k, v int) {
	m.m[k] = v
	m.dirty = true // that's why m is pointer *SortedMapInt2Int
}
func (m *SortedMapInt2Int) Inc(k int) {
	m.m[k]++
	m.dirty = true // that's why m is pointer *SortedMapInt2Int
}
func (m SortedMapInt2Int) Get(k int) (v int, ok bool) {
	v, ok = m.m[k]
	return
}

func (m SortedMapInt2Int) MaxKey() int {
	if len(m.m) == 0 {
		panic("Max of empty map has no answer")
	}
	if m.dirty {
		m.updateSort()
		m.dirty = false
	}
	lastIdx := len(m.sortedKeys) - 1
	return m.sortedKeys[lastIdx]
}

func (m SortedMapInt2Int) MinKey() int {
	if len(m.m) == 0 {
		panic("Min of empty map has no answer")
	}
	if m.dirty {
		m.updateSort()
		m.dirty = false
	}
	return m.sortedKeys[0]
}

func (m *SortedMapInt2Int) SortedKeys() []int {
	if m.dirty {
		m.updateSort()
		m.dirty = false
	}
	return m.sortedKeys
}

// Returns the statified sorted map: [Sortindex][Key][Value]
// This could also be achieved by iterating SortedKeys() and calling Get(v).
// It may be worthwhile for external packages to save Get() calls and for repetitive iterations
func (m *SortedMapInt2Int) SortedKV() []sortedKV {
	if m.dirty {
		m.updateSort()
		m.dirty = false
	}
	m.sKV = make([]sortedKV, len(m.m))
	for k, v := range m.sortedKeys {
		m.sKV[k] = sortedKV{v, m.m[v]}
	}
	return m.sKV
}

// Again: [Sortindex][Key][Value]
// This time: Sortindex exists also for *empty* keys
// and for keys [0...MinKey)
// Motivation: In case we have some sorted data (outside),
// then we do not want to iterate to find a specific value.
// Instead we want to *pinpoint* a distinct instance.
// Therefore we need the offsets to these instances.
// These offsets can be seen as *partial sums*
func (m SortedMapInt2Int) SortedPartialSums() ([]sortedKV, []int) {
	groundedMin := m.MinKey()
	if groundedMin > 0 {
		groundedMin = 0
	}
	skv := make([]sortedKV, m.MaxKey()-groundedMin+1)
	flattened := make([]int, len(skv))
	cnt := 0
	j := 0
	for i := groundedMin; i < m.MaxKey()+1; i++ {
		if v2, ok := m.Get(i); ok {
			cnt += v2
		}
		skv[j].K = i
		skv[j].V = cnt

		flattened[j] = cnt

		j++
	}
	return skv, flattened
}
