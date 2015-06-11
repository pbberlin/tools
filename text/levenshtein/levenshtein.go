package levenshtein

import (
	"fmt"
	"strings"

	"github.com/pbberlin/tools/util"
)

type Equaler interface {
	Equal(compare2 interface{}) bool
}

type Matrix [][]int

// LSDist returns the edit distance between
// two slices of Comparables
func LSDist(toks1, toks2 []Equaler, opt Options) int {
	return MatrixForSlices(toks1, toks2, opt).Distance()
}

// MatrixForSlices generates a 2-D array representing the dynamic programming
// table used by the Levenshtein algorithm, as described e.g. here:
// http://www.let.rug.nl/kleiweg/lev/
// The reason for putting the creation of the table into a separate function is
// that it cannot only be used for reading of the edit distance between two
// strings, but also e.g. to backtrace an edit script that provides an
// alignment between the characters of both strings.
func MatrixForSlices(rows, cols []Equaler, opt Options) Matrix {
	// Make a 2-D matrix. Rows correspond to prefixes of source, columns to
	// prefixes of target. Cells will contain edit distances.
	// Cf. http://www.let.rug.nl/~kleiweg/lev/levenshtein.html
	h := len(rows) + 1
	w := len(cols) + 1
	mx := make([][]int, h)

	// Initialize trivial distances (from/to empty string):
	// Filling the left column and the top row with row/column indices.
	for i := 0; i < h; i++ {
		mx[i] = make([]int, w)
		mx[i][0] = i
	}
	for j := 1; j < w; j++ {
		mx[0][j] = j
	}

	// Filling the remaining cells:
	// 	For each prefix pair:
	// 		Choose couple {edit history, operation} with lowest cost.
	for i := 1; i < h; i++ {
		for j := 1; j < w; j++ {
			delCost := mx[i-1][j] + opt.DelCost
			matchSubCost := mx[i-1][j-1]
			if !(rows[i-1]).Equal(cols[j-1]) {
				matchSubCost += opt.SubCost
			}
			insCost := mx[i][j-1] + opt.InsCost
			mx[i][j] = min(delCost, min(matchSubCost, insCost))
		}
	}
	PrintMatrix(rows, cols, mx)
	return mx
}

// EditScript returns an optimal edit script
// turning source into target.
func EditScript(src, dst []Equaler, opt Options) TEditScrpt {
	return backtrace(len(src), len(dst), MatrixForSlices(src, dst, opt), opt)
}

func backtrace(i int, j int, mx [][]int, opt Options) TEditScrpt {
	if i > 0 && mx[i-1][j]+opt.DelCost == mx[i][j] {
		return append(backtrace(i-1, j, mx, opt), Del)
	}
	if j > 0 && mx[i][j-1]+opt.InsCost == mx[i][j] {
		return append(backtrace(i, j-1, mx, opt), Ins)
	}
	if i > 0 && j > 0 && mx[i-1][j-1]+opt.SubCost == mx[i][j] {
		return append(backtrace(i-1, j-1, mx, opt), Sub)
	}
	if i > 0 && j > 0 && mx[i-1][j-1] == mx[i][j] {
		return append(backtrace(i-1, j-1, mx, opt), Match)
	}
	return []EditOp{}
}

func min(a int, b int) int {
	if b < a {
		return b
	}
	return a
}

// Distance returns edit distance for an existing matrix.
func (mx Matrix) Distance() int {
	return mx[len(mx)-1][len(mx[0])-1]
}

// EditScript returns an optimal edit script for an existing matrix.
func (mx Matrix) EditScript() TEditScrpt {
	return backtrace(len(mx[0])-1, len(mx)-1, mx, DefaultOptions)
}

// PrintMatrix prints a visual representation of matrix to os.Stderr
func PrintMatrix(rows, cols []Equaler, mx [][]int) {

	fp := fmt.Printf

	const cl = 11 // column length
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
	fp("\n")
}
