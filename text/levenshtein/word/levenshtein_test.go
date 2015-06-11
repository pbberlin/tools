package word

import (
	"fmt"

	ls_core "github.com/pbberlin/tools/text/levenshtein"
	"github.com/pbberlin/tools/util"

	"testing"
)

type TestCase struct {
	src      []Token
	dst      []Token
	distance int
}

var testCases = []TestCase{

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

func init() {
	sss := [][]string{
		[]string{"word1", "word2", "up"},
		[]string{"trink", "nicht", "so", "viel", "Kaffee"},
		[]string{"nicht", "fuer", "Kinder", "ist", "Tuerkentrank"},
	}

	for i := 0; i < len(sss); i++ {
		st := make([]Token, 0, len(sss[i]))
		for j := 0; j < len(sss[i]); j++ {
			st = append(st, Token(sss[i][j]))
		}

		prev := testCases[len(testCases)-1]
		// log.Printf("%v", prev)
		testCases = append(testCases, TestCase{src: st, dst: prev.src, distance: 2})
	}

}

func TestLevenshtein(t *testing.T) {
	for _, tc := range testCases {

		got := ls_core.LSDist(
			convertToCore(tc.src),
			convertToCore(tc.dst),
			ls_core.DefaultOptions)
		ssrc := fmt.Sprintf("%v", tc.src)
		sdst := fmt.Sprintf("%v", tc.dst)
		if got != tc.distance {
			t.Logf(
				"Distance between %20v and %20v should be %v - but got %v ",
				util.Ellipsoider(ssrc, 8),
				util.Ellipsoider(sdst, 8), tc.distance, got)
			t.Fail()
		}
	}

	for _, tc := range testCases {
		ls_core.EditScript(convertToCore(tc.src), convertToCore(tc.dst), ls_core.DefaultOptions)
	}
}
