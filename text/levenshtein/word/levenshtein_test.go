package word

import (
	ls_core "github.com/pbberlin/tools/text/levenshtein"

	"testing"
)

var testCases = []struct {
	src      []Token
	dst      []Token
	distance int
}{

	{[]Token{}, []Token{"word1"}, 1},
	{[]Token{"word1"}, []Token{"word1", "word1"}, 1},
	{[]Token{"word1"}, []Token{"word1", "word1", "word1"}, 2},

	{[]Token{}, []Token{}, 0},
	{[]Token{"word1"}, []Token{"word2"}, 2},
	{[]Token{"word1", "word1", "word1"}, []Token{"word1", "word2", "word1"}, 2},
	{[]Token{"word1", "word1", "word1"}, []Token{"word1", "word2"}, 3},

	{[]Token{"word1"}, []Token{"word1"}, 0},
	{[]Token{"word1", "word2"}, []Token{"word1", "word2"}, 0},
	{[]Token{"word1"}, []Token{}, 1},

	{[]Token{"word1", "word2"}, []Token{"word1"}, 1},
	{[]Token{"word1", "word2", "word3"}, []Token{"word1"}, 2},
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
