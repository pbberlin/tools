package word

import (
	lscore "github.com/pbberlin/tools/text/levenshtein"

	"testing"
)

var testCases = []struct {
	src      []Token
	dst      []Token
	distance int
}{
	{[]Token{"word1", "word2"}, []Token{"word1"}, 1},
	{[]Token{"word1", "word2", "word3"}, []Token{"word1"}, 2},
}

func TestLevenshtein(t *testing.T) {
	for _, testCase := range testCases {

		for k, v := range testCase.src {
			if v == "word1" {
				another := Token("word1")
				testCase.src[k] = another
			}
		}

		got := lscore.DistanceOfSlices(
			convertToCore([]Token(testCase.src)),
			convertToCore([]Token(testCase.dst)),
			lscore.DefaultOptions)
		if got != testCase.distance {
			t.Logf(
				"Distance between %v and %v should be %v - but got %v ",
				testCase.dst, testCase.src, testCase.distance, got)
			t.Fail()
		}
	}
}
