package tplx

import (
	"github.com/pbberlin/tools/shared_structs"
)

func fMult(x, y int) int {
	return x * y
}
func fAdd(x, y int) int {
	return x + y
}

func fMakeRange(num int) []int {
	sl := make([]int, num)
	for i, _ := range sl {
		sl[i] = i
	}
	return sl
}

// use "index entity idx" instead - see
//	http://stackoverflow.com/questions/12701452/golang-html-template-how-to-index-a-slice-element
func fAccessElement(v []interface{}, i int) interface{} {
	return v[i]
}

func fAccessElementB2(v []shared_structs.B2, i int) shared_structs.B2 {
	return v[i]
}

func fChop(s string, i int) string {
	if len(s) > i {
		return s[i:]
	}
	return s
}

func fNumCols(ncols, perRow int) []int {

	nrows := ncols/perRow + 1

	var ret = make([]int, nrows)

	for i, _ := range ret {
		if i < nrows-1 {
			ret[i] = perRow

		} else {
			ret[i] = ncols % perRow
		}
	}
	return ret
}
