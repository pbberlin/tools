// Package levenshtein core computes the edit distance of two slices of tokens,
// of slim interface type Equaler; subpackages provide various granularity.
// Tokens must be of interface type <Equaler> - implementing (tok) Equal(tok) bool.
// An edit script for converting slice1 to slice2 can also be derived.
// Preference for substitution over insertion/deletion is configurable.
package levenshtein

import (
	"fmt"
	"strings"

	"github.com/pbberlin/tools/stringspb"
	"github.com/pbberlin/tools/util"
)

const cl = 11 // column length for Print funcs

// Equaler is the neccessary interface to compute the levenshtein distance.
type Equaler interface {
	Equal(compare2 interface{}) bool
}

// works also for idx == -1
func insertAfter(s []Equaler, idx int, newVal Equaler) []Equaler {
	if idx > len(s)-1 {
		panic("Cannot insert beyond existing length")
	}
	s = append(s, nil)
	copy(s[idx+2:], s[idx+1:])
	s[idx+1] = newVal
	return s
}

func deleteAt(s []Equaler, idx int) []Equaler {
	if idx > len(s)-1 {
		panic("Cannot delete beyond existing length")
	}
	copy(s[idx:], s[idx+1:])
	s = s[:len(s)-1]
	return s
}

// The internal levenshtein matrix is only exported,
// because calling packages need to declare its type.
type Matrix struct {
	mx         [][]int
	rows, cols []Equaler
	opt        Options
}

// New generates a 2-D array,
// representing the dynamic programming table
// used by the Levenshtein algorithm.
// Compare http://www.let.rug.nl/kleiweg/lev/.
// Matrix can be used for retrieval of edit distance
// and for backtrace scripts
func New(argRows, argCols []Equaler, opt Options) Matrix {

	// Make a 2-D matrix. Rows correspond to prefixes of source, columns to
	// prefixes of target. Cells will contain edit distances.
	// Cf. http://www.let.rug.nl/~kleiweg/lev/levenshtein.html
	m := Matrix{}
	m.opt = opt
	m.rows = argRows
	m.cols = argCols
	h := len(m.rows) + 1
	w := len(m.cols) + 1
	m.mx = make([][]int, h)

	// Initialize trivial distances (from/to empty string):
	// Filling the left column and the top row with row/column indices.
	for i := 0; i < h; i++ {
		m.mx[i] = make([]int, w)
		m.mx[i][0] = i
	}
	for j := 1; j < w; j++ {
		m.mx[0][j] = j
	}

	// Filling the remaining cells:
	// 	For each prefix pair:
	// 		Choose couple {edit history, operation} with lowest cost.
	for i := 1; i < h; i++ {
		for j := 1; j < w; j++ {
			delCost := m.mx[i-1][j] + opt.DelCost
			matchSubCost := m.mx[i-1][j-1]
			if !(m.rows[i-1]).Equal(m.cols[j-1]) {
				matchSubCost += opt.SubCost
			}
			insCost := m.mx[i][j-1] + opt.InsCost
			m.mx[i][j] = min(delCost, min(matchSubCost, insCost))
		}
	}

	return m
}

func min(a int, b int) int {
	if b < a {
		return b
	}
	return a
}

// Distance returns levenshtein edit distance for the two slices of tokens of m.
func (m *Matrix) Distance() (int, float64) {

	dist := m.mx[len(m.mx)-1][len(m.mx[0])-1]

	relDist := 0.0

	ls1, ls2 := len(m.mx), len(m.mx[0])

	// First relDist computation:
	// 		1.) compensated for the size
	// 		2.) related to the smaller slice
	// Can lead to Zero, when diff == dist
	diff := util.Abs(ls1 - ls2)
	if ls1 >= ls2 { // row > col
		relDist = float64(dist-diff) / float64(ls2)
	} else {
		relDist = float64(dist-diff) / float64(ls1)
	}

	// Second relDist: Simply related to the larger slice.
	// Also account for ls1 and ls2 being one larger than the practical number of tokens.
	divisor := float64(ls1)
	if ls2 > ls1 { // row > col
		divisor = float64(ls2)
	}
	divisor--
	if divisor == 0.0 {
		divisor = 1.0
	}
	relDist = float64(dist) / divisor
	if relDist == 0.25 {
		fmt.Printf("dist %v - ls1 %v - relDist %5.2v\n", dist, divisor, relDist)
	}

	return dist, relDist
}

// EditScript returns an optimal edit script for an existing matrix.
func (m *Matrix) EditScript() TEditScrpt {
	return m.backtrace(len(m.mx)-1, len(m.mx[0])-1)
}

// backtrace is recursive.
// It starts bottom right and steps left/top/lefttop
func (m *Matrix) backtrace(i, j int) TEditScrpt {

	pf := func(str string) {}
	// pf := fmt.Printf

	mx := m.mx
	opt := m.opt
	eo := EditOpExt{}

	if i > 0 && mx[i-1][j]+opt.DelCost == mx[i][j] {
		pf("c1")
		eo.op = Del
		eo.src = i - 1
		eo.dst = j
		return append(m.backtrace(i-1, j), eo)
	}
	if j > 0 && mx[i][j-1]+opt.InsCost == mx[i][j] {
		pf("c2")
		eo.op = Ins
		eo.src = i
		eo.dst = j - 1
		return append(m.backtrace(i, j-1), eo)
	}
	if i > 0 && j > 0 && mx[i-1][j-1]+opt.SubCost == mx[i][j] {
		pf("c3")
		eo.op = Sub
		eo.src = i - 1
		eo.dst = j - 1
		return append(m.backtrace(i-1, j-1), eo)
	}
	if i > 0 && j > 0 && mx[i-1][j-1] == mx[i][j] {
		pf("c4")
		eo.op = Match
		eo.src = i - 1
		eo.dst = j - 1
		return append(m.backtrace(i-1, j-1), eo)
	}
	pf("c5")
	return []EditOpExt{}
}

// Print prints a visual representation
// of the slices of tokens and their distance matrix
func (m *Matrix) Print() {

	rows, cols := m.rows, m.cols
	mx := m.mx

	fp := fmt.Printf

	fmt2 := fmt.Sprintf("%s-%vd", "%", cl)

	fp(strings.Repeat(" ", 2*cl))
	for _, col := range cols {
		scol := fmt.Sprintf("%v", col)
		fp("%v ", stringspb.ToLen(scol, cl-1)) // at least one space right
	}
	fp("\n")

	fp(strings.Repeat(" ", cl))
	fp(fmt2, mx[0][0])
	for j, _ := range cols {
		fp(fmt2, mx[0][j+1])
	}
	fp("\n")

	//
	for i, row := range rows {
		srow := fmt.Sprintf("%v", row)
		fp("%v ", stringspb.ToLen(srow, cl-1)) // at least one space right
		fp(fmt2, mx[i+1][0])
		for j, _ := range cols {
			fp(fmt2, mx[i+1][j+1])
		}
		fp("\n")
	}
	// fp("\n")
}

// ApplyEditScript applies the given Editscript
// to the first slice of tokens of m.
// The returned slice should be equal
// to the second slice of tokens of m.
func (m *Matrix) ApplyEditScript(es TEditScrpt) []Equaler {

	sumIns := 0
	sumDel := 0
	fmt2 := fmt.Sprintf("%s-%vv", "%", cl)

	rows2 := make([]Equaler, 0, len(m.rows))
	for _, v := range m.rows {
		rows2 = append(rows2, v)
	}

	const offs = 1
	fmt.Printf("%v", strings.Repeat(" ", 2*cl))
	for _, v := range es {

		s := fmt.Sprintf("%v-%v-%v", v.op, offs+v.src+sumIns-sumDel, offs+v.dst)
		fmt.Printf(fmt2, s)

		pos := v.src + sumIns - sumDel

		if v.op == Ins {
			// rows2 = insertAfter(rows2, util.Min(pos, len(rows2)-1), m.cols[v.dst])
			rows2 = insertAfter(rows2, pos-1, m.cols[v.dst])
			sumIns++
		}

		if v.op == Del {
			rows2 = deleteAt(rows2, pos)
			sumDel++
		}

		if v.op == Sub {
			rows2[pos] = m.cols[v.dst]
		}
	}
	fmt.Printf("\n")

	fmt.Printf("%v", strings.Repeat(" ", 2*cl))
	for _, row := range rows2 {
		fmt.Printf(fmt2, row)
	}
	return rows2

}

// CompareToCol takes a slice of Equaler-Tokens
// and compares them against the second matrix slice.
func (m *Matrix) CompareToCol(col2 []Equaler) bool {
	equal := true
	for idx, v := range m.cols {
		if v != col2[idx] {
			equal = false
			break
		}
	}
	return equal
}
