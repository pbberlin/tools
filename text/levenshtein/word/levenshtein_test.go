package word

import (
	"fmt"

	"github.com/pbberlin/tools/pbstrings"
	ls_core "github.com/pbberlin/tools/text/levenshtein"

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

	{[]Token{"trink", "nicht", "so", "viel", "Kaffee"},
		[]Token{"nicht", "für", "Kinder", "ist", "der", "Türkentrank"}, 9},

	{[]Token{"ihn", "nicht", "der", "lassen", "kann"},
		[]Token{"nicht", "für", "Kinder", "ist", "der", "Türkentrank"}, 7},
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
				i, pbstrings.Ellipsoider(ssrc, 8), pbstrings.Ellipsoider(sdst, 8), tc.distance, got)
			t.Fail()
		}

		m.Print()
		fmt.Printf("\n")

		es := m.EditScript()
		// es.Print()
		// fmt.Printf("\n")

		got2 := m.ApplyEditScript(es)
		if !m.CompareToCol(got2) {
			t.Logf("\nwnt %v \ngot %v ", convertToCore(tc.dst), got2)
			t.Fail()
		}

		fmt.Printf("\n")
		fmt.Printf("\n")

	}

}
