package levenshtein

import (
	"fmt"
	"strings"

	"github.com/pbberlin/tools/util"
)

const cl = 11 // column length for Print funcs

type Equaler interface {
	Equal(compare2 interface{}) bool
}

// works also for idx == -1
func InsertAfter(s []Equaler, idx int, newVal Equaler) []Equaler {
	if idx > len(s)-1 {
		panic("Cannot insert beyond existing length")
	}
	s = append(s, nil)
	copy(s[idx+2:], s[idx+1:])
	s[idx+1] = newVal
	return s
}

func Delete(s []Equaler, idx int) []Equaler {
	if idx > len(s)-1 {
		panic("Cannot delete beyond existing length")
	}
	copy(s[idx:], s[idx+1:])
	s = s[:len(s)-1]
	return s
}

type Matrix struct {
	mx         [][]int
	rows, cols []Equaler
	opt        Options
}

// New generates a 2-D array,
// representing the dynamic programming table
// used by the Levenshtein algorithm.
// Compare http://www.let.rug.nl/kleiweg/lev/
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

// Distance returns edit distance for an existing matrix.
func (m Matrix) Distance() int {
	return m.mx[len(m.mx)-1][len(m.mx[0])-1]
}

// EditScript returns an optimal edit script for an existing matrix.
func (m Matrix) EditScript() TEditScrpt {
	return m.backtrace(len(m.mx)-1, len(m.mx[0])-1)
}

// Backtrace is recursive
// It starts bottom right and steps left/top/lefttop
func (m Matrix) backtrace(i, j int) TEditScrpt {

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

// PrintTokensWithMatrix prints a visual representation
// of the slices of tokens and their distance matrix
func (m Matrix) Print() {

	rows, cols := m.rows, m.cols
	mx := m.mx

	fp := fmt.Printf

	fmt2 := fmt.Sprintf("%s-%vd", "%", cl)

	fp(strings.Repeat(" ", 2*cl))
	for _, col := range cols {
		scol := fmt.Sprintf("%v", col)
		fp("%v ", util.ToLen(scol, cl-1)) // at least one space right
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
		fp("%v ", util.ToLen(srow, cl-1)) // at least one space right
		fp(fmt2, mx[i+1][0])
		for j, _ := range cols {
			fp(fmt2, mx[i+1][j+1])
		}
		fp("\n")
	}
	// fp("\n")
}

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
			// rows2 = InsertAfter(rows2, util.Min(pos, len(rows2)-1), m.cols[v.dst])
			rows2 = InsertAfter(rows2, pos-1, m.cols[v.dst])
			sumIns++
		}

		if v.op == Del {
			rows2 = Delete(rows2, pos)
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
// and compares them against the column tokens of m.
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
