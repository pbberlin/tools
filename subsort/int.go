package subsort

import "sort"

type SortedByIntVal struct {
	IdxOrig int // index of the original slice
	Val     int // dynamically filled value to be sorted by
}

type sliceSortableIntAsc []SortedByIntVal
type sliceSortableIntDesc []SortedByIntVal

func (s sliceSortableIntAsc) Len() int {
	return len(s)
}

func (s sliceSortableIntAsc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sliceSortableIntAsc) Less(i, j int) bool {
	if s[i].Val < s[j].Val {
		return true
	}
	return false
}
func (s sliceSortableIntDesc) Len() int {
	return len(s)
}

func (s sliceSortableIntDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sliceSortableIntDesc) Less(i, j int) bool {
	if s[i].Val > s[j].Val {
		return true
	}
	return false
}

//
// A second way only requires the callee
// only to submit a func, that extracts the sort value by index:
type ExtractIntFielder func(i int) int

//
// An interface ExtractIntFielder, would require all callee slices to be wrapped.
func SortByIntValAsc(size int, f ExtractIntFielder) []SortedByIntVal {

	copyOfSubset := make([]SortedByIntVal, size)
	for i := 0; i < size; i++ {
		copyOfSubset[i].IdxOrig = i
		copyOfSubset[i].Val = f(i)
	}

	wrapperSortable := sliceSortableIntAsc(copyOfSubset)
	sort.Sort(wrapperSortable)
	unwrap := []SortedByIntVal(wrapperSortable)
	return unwrap

}
func SortByIntValDesc(size int, f ExtractIntFielder) []SortedByIntVal {

	copyOfSubset := make([]SortedByIntVal, size)
	for i := 0; i < size; i++ {
		copyOfSubset[i].IdxOrig = i
		copyOfSubset[i].Val = f(i)
	}

	wrapperSortable := sliceSortableIntDesc(copyOfSubset)
	sort.Sort(wrapperSortable)
	unwrap := []SortedByIntVal(wrapperSortable)
	return unwrap

}
