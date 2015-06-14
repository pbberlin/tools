package subsort

// File contains examples of other strategies for flexible sorting.
// I considered and rejected these examples.
// Then I developed my own versatile sorting helper.

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

//--------------------------------------------
// Precursor 1

// A type, containing function headers for the sort interface.
type sortI struct {
	l    int // not len(), but its value
	less func(int, int) bool
	swap func(int, int)
}

// The methods of the sort interface
// are now satisfied, using the *members* of the struct.
func (s *sortI) Len() int {
	return s.l
}

func (s *sortI) Less(i, j int) bool {
	return s.less(i, j)
}

func (s *sortI) Swap(i, j int) {
	s.swap(i, j)
}

// SortI can now be used as follows:
//   sort.Sort( &sortI{} )

// SortF wraps the construction
// and usage of a sortI instance.
func SortF(Len int, Less func(int, int) bool, Swap func(int, int)) {
	si := &sortI{l: Len, less: Less, swap: Swap}
	sort.Sort(si)
}

// We can now take any int slice
// and construct two-and-a-half closures with it
// and pass the closures (with our int slice implicitly piggypacked) to SortF
func TestSortI_demo(t *testing.T) {
	ints := []int{3, 4, 1, 7, 0}
	SortF(len(ints), func(i, j int) bool {
		return ints[i] < ints[j]
	}, func(i, j int) {
		ints[i], ints[j] = ints[j], ints[i]
	})
	want := fmt.Sprintf("%#v", []int{0, 1, 3, 4, 7})
	got := fmt.Sprintf("%#v", ints)
	if want != got {
		t.Errorf("wanted vs got: \n%v \n%v", want, got)
	}
}

// This approach still needs typespecific variations of sortI.
// And it requires the closure notation.

//
//--------------------------------------------

// precursor2 - using interface{}
// copyAndSort() first produces a copy,
// the copy containing only the desired data.
// Then this subset copy is sorted.
// Sadly, the argument of type []interface{} is expensive to create
// from calling packages. It mostly involves previous copying element by element
// thus genericism causes *two* rounds of copying.
// For this reason we prefer the SortByVal() func,
// width the only downside, that the preparation of the subset slice
// needs to be done by the calling package.
// Otherwise we are left with only *one* round of copying
// and without the need for reflection.
func demo__copyAndSort(sArg []interface{}, fieldname string) []SortedByStringVal {

	copyOfSubset := []SortedByStringVal{}
	for i := 0; i < len(sArg); i++ {
		lp := sArg[i]
		immutable := reflect.ValueOf(lp)
		// reflect.Value.String() does not panic upon non-strings.
		// Instead it returns "<type value>" for non-string
		dynVal := immutable.FieldByName(fieldname).String()
		copyOfSubset = append(copyOfSubset, SortedByStringVal{IdxOrig: i, Val: dynVal})
	}

	wrapperSortable := sliceSortableStringAsc(copyOfSubset)
	sort.Sort(wrapperSortable)
	unwrap := []SortedByStringVal(wrapperSortable)
	return unwrap

}

//
//
func demo__GetFieldValueByName() {
	type MyStruct struct {
		SortBy string
	}
	myStruct := MyStruct{"001"}
	immutable := reflect.ValueOf(myStruct)
	val := immutable.FieldByName("SortBy").String()
	_ = val
}
