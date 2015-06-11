package rune

import (
	ls_core "github.com/pbberlin/tools/text/levenshtein"

	"testing"
)

var testCases = []struct {
	src      []Token
	dst      []Token
	distance int
}{
	{[]Token{}, []Token{'a'}, 1},
	{[]Token{'a'}, []Token{'a', 'a'}, 1},
	{[]Token{'a'}, []Token{'a', 'a', 'a'}, 2},

	{[]Token{}, []Token{}, 0},
	{[]Token{'a'}, []Token{'b'}, 2},
	{[]Token{'a', 'a', 'a'}, []Token{'a', 'b', 'a'}, 2},
	{[]Token{'a', 'a', 'a'}, []Token{'a', 'b'}, 3},

	{[]Token{'a'}, []Token{'a'}, 0},
	{[]Token{'a', 'b'}, []Token{'a', 'b'}, 0},
	{[]Token{'a'}, []Token{}, 1},

	{[]Token{'a', 'a'}, []Token{'a'}, 1},
	{[]Token{'a', 'a', 'a'}, []Token{'a'}, 2},

	// {[]Token{'a'}, []Token{'a'}, 220},
}

func TestLevenshtein(t *testing.T) {
	for _, testCase := range testCases {

		got := ls_core.LSDist(
			convertToCore([]Token(testCase.src)),
			convertToCore([]Token(testCase.dst)),
			ls_core.DefaultOptions)
		if got != testCase.distance {
			t.Logf(
				"Distance between %v and %v should be %v - but got %v ",
				testCase.dst, testCase.src, testCase.distance, got)
			t.Fail()
		}
	}
}
