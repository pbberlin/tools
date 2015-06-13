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

	// Edit Script Test Cases
	{[]Token{"wd1", "wd2", "wd1"}, []Token{"wd1", "wd3", "wd1"}, 2},
	{[]Token{"wd2", "wd1", "wd1", "wd1", "wd2", "wd1"}, []Token{"wd1", "wd1"}, 4},

	//
	{[]Token{}, []Token{"wd1"}, 1},
	{[]Token{"wd1"}, []Token{"wd1", "wd1"}, 1},
	{[]Token{"wd1"}, []Token{"wd1", "wd1", "wd1"}, 2},

	{[]Token{}, []Token{}, 0},
	{[]Token{"wd1"}, []Token{"wd2"}, 2},
	{[]Token{"wd1", "wd1", "wd1"}, []Token{"wd1", "wd2", "wd1"}, 2},
	{[]Token{"wd1", "wd1", "wd1"}, []Token{"wd1", "wd2"}, 3},

	{[]Token{"wd1"}, []Token{"wd1"}, 0},
	{[]Token{"wd1", "wd2"}, []Token{"wd1", "wd2"}, 0},
	{[]Token{"wd1"}, []Token{}, 1},

	{[]Token{"wd1", "wd2"}, []Token{"wd1"}, 1},
	{[]Token{"wd1", "wd2", "wd3"}, []Token{"wd1"}, 2},

	{[]Token{"wd1", "wd1", "wd1"}, []Token{"wd1", "wd2", "wd1", "wd3"}, 3},
}

func init() {
	sss := [][]string{
		[]string{"wd1", "wd2", "up"},
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
	for i, tc := range testCases {

		m := ls_core.New(convertToCore(tc.src), convertToCore(tc.dst), ls_core.DefaultOptions)
		got := m.Distance()

		ssrc := fmt.Sprintf("%v", tc.src)
		sdst := fmt.Sprintf("%v", tc.dst)
		if got != tc.distance {
			t.Logf(
				"%2v: Distance between %20v and %20v should be %v - but got %v ",
				i, util.Ellipsoider(ssrc, 8), util.Ellipsoider(sdst, 8), tc.distance, got)
			t.Fail()
		}

		if i > 1 {
			// continue
		}

		m.Print()
		fmt.Printf("\n")

		es := m.EditScript()
		// es.Print()
		// fmt.Printf("\n")

		got2 := m.ApplyEditScript(es)
		if !m.CompareToCol(got2) {
			t.Logf(
				"wnt %v \ngot %v ", convertToCore(tc.dst), got2)
			t.Fail()

		}

		fmt.Printf("\n")
		fmt.Printf("\n")

	}

}
