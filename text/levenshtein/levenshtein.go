package levenshtein

import (
	"fmt"
	"os"
)

// type Token interface {
// 	Matches(cmp1, cmp2 interface{}) bool
// }

type Token interface {
	Matches(cmp1, cmp2 interface{}) bool
}

// DistanceOfSlices returns the edit distance between
// slices of tokens
func DistanceOfSlices(tokens1, tokens2 []Token, opt Options) int {
	return DistanceForMatrix(MatrixForStrings(tokens1, tokens2, opt))
}

// DistanceForMatrix reads the edit distance off the given Levenshtein matrix.
func DistanceForMatrix(matrix [][]int) int {
	return matrix[len(matrix)-1][len(matrix[0])-1]
}

// MatrixForStrings generates a 2-D array representing the dynamic programming
// table used by the Levenshtein algorithm, as described e.g. here:
// http://www.let.rug.nl/kleiweg/lev/
// The reason for putting the creation of the table into a separate function is
// that it cannot only be used for reading of the edit distance between two
// strings, but also e.g. to backtrace an edit script that provides an
// alignment between the characters of both strings.
func MatrixForStrings(rows, cols []Token, opt Options) [][]int {
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
			if !(rows[i-1]).Matches(rows[i-1], cols[j-1]) {
				matchSubCost += opt.SubCost
			}
			insCost := mx[i][j-1] + opt.InsCost
			mx[i][j] = min(delCost, min(matchSubCost, insCost))
		}
	}
	//LogMatrix(rows, cols, mx)
	return mx
}

// EditScriptForStrings returns an optimal edit script
// turning source into target.
func EditScriptForStrings(src, dst []Token, opt Options) EditScript {
	return backtrace(len(src), len(dst), MatrixForStrings(src, dst, opt), opt)
}

// EditScriptForMatrix returns an optimal edit script based on the given
// Levenshtein matrix.
func EditScriptForMatrix(matrix [][]int, opt Options) EditScript {
	return backtrace(len(matrix[0])-1, len(matrix)-1, matrix, opt)
}

// LogMatrix prints a visual representation of matrix to os.Stderr
func LogMatrix(src, dst []Token, mx [][]int) {

	fp := func(format string, args ...interface{}) { fmt.Fprintf(os.Stderr, format, args) }

	fp("    ")
	for _, dstX := range dst {
		fp("  %c", dstX)
	}
	fp("\n")
	fp("  %2d", mx[0][0])
	for j, _ := range dst {
		fp(" %2d", mx[0][j+1])
	}
	fp("\n")
	for i, srcX := range src {
		fp("%c %2d", srcX, mx[i+1][0])
		for j, _ := range dst {
			fp(" %2d", mx[i+1][j+1])
		}
		fp("\n")
	}
}

func backtrace(i int, j int, mx [][]int, opt Options) EditScript {
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
