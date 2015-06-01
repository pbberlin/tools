package subsort

import "sort"

type SortedByStringVal struct {
	IdxOrig int    // index of the original slice
	Val     string // dynamically filled value to be sorted by
}

type sliceSortableStringAsc []SortedByStringVal
type sliceSortableStringDesc []SortedByStringVal

func (s sliceSortableStringAsc) Len() int {
	return len(s)
}

func (s sliceSortableStringAsc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sliceSortableStringAsc) Less(i, j int) bool {
	if s[i].Val < s[j].Val {
		return true
	}
	return false
}
func (s sliceSortableStringDesc) Len() int {
	return len(s)
}

func (s sliceSortableStringDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sliceSortableStringDesc) Less(i, j int) bool {
	if s[i].Val > s[j].Val {
		return true
	}
	return false
}

//
// A second way only requires the callee
// only to submit a func, that extracts the sort value by index:
type ExtractStringFielder func(i int) string

//
// An interface ExtractStringFielder, would require all callee slices to be wrapped.
func SortByStringValAsc(size int, f ExtractStringFielder) []SortedByStringVal {

	copyOfSubset := make([]SortedByStringVal, size)
	for i := 0; i < size; i++ {
		copyOfSubset[i].IdxOrig = i
		copyOfSubset[i].Val = f(i)
	}

	wrapperSortable := sliceSortableStringAsc(copyOfSubset)
	sort.Sort(wrapperSortable)
	unwrap := []SortedByStringVal(wrapperSortable)
	return unwrap

}
func SortByStringValDesc(size int, f ExtractStringFielder) []SortedByStringVal {

	copyOfSubset := make([]SortedByStringVal, size)
	for i := 0; i < size; i++ {
		copyOfSubset[i].IdxOrig = i
		copyOfSubset[i].Val = f(i)
	}

	wrapperSortable := sliceSortableStringDesc(copyOfSubset)
	sort.Sort(wrapperSortable)
	unwrap := []SortedByStringVal(wrapperSortable)
	return unwrap

}
