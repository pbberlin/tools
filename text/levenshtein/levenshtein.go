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
	return m.Backtrace(len(m.mx)-1, len(m.mx[0])-1)
}

func (m Matrix) Backtrace(i, j int) TEditScrpt {

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
		return append(m.Backtrace(i-1, j), eo)
	}
	if j > 0 && mx[i][j-1]+opt.InsCost == mx[i][j] {
		pf("c2")
		eo.op = Ins
		eo.src = i
		eo.dst = j - 1
		return append(m.Backtrace(i, j-1), eo)
	}
	if i > 0 && j > 0 && mx[i-1][j-1]+opt.SubCost == mx[i][j] {
		pf("c3")
		eo.op = Sub
		eo.src = i - 1
		eo.dst = j - 1
		return append(m.Backtrace(i-1, j-1), eo)
	}
	if i > 0 && j > 0 && mx[i-1][j-1] == mx[i][j] {
		pf("c4")
		eo.op = Match
		eo.src = i - 1
		eo.dst = j - 1
		return append(m.Backtrace(i-1, j-1), eo)
	}
	pf("c5")
	return []EditOpExt{}
}

// PrintTokensWithMatrix prints a visual representation
// of the slices of tokens
// and of their diff matrix
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
