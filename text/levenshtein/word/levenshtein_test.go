package word

import (
	"fmt"

	"github.com/pbberlin/tools/pbstrings"
	ls_core "github.com/pbberlin/tools/text/levenshtein"

	"testing"
)

type TestCase struct {
	src       string
	dst       string
	distances []int
}

var testCasesBasic = []TestCase{

	// Edit Script Test Cases
	{"", "wd1", []int{1, 1}},
	{"wd1", "wd1 wd1", []int{1, 1}},
	{"wd1", "wd1 wd1 wd1", []int{2, 2}},

	{"", "", []int{0, 0}},
	{"wd1", "wd2", []int{2, 2}},
	{"wd1 wd1 wd1", "wd1 wd2 wd1", []int{2, 2}},
	{"wd1 wd1 wd1", "wd1 wd2", []int{3, 3}},

	{"wd1", "wd1", []int{0, 0}},
	{"wd1 wd2", "wd1 wd2", []int{0, 0}},
	{"wd1", "", []int{1, 1}},

	{"wd1 wd2", "wd1", []int{1, 1}},
	{"wd1 wd2 wd3", "wd1", []int{2, 2}},

	{"wd1 wd1 wd1", "wd1 wd2 wd1 wd3", []int{3, 3}},
}

var testCasesAdv1 = []TestCase{
	{"trink nicht so viel Kaffee",
		"nicht für Kinder ist der Türkentrank", []int{9, 9}},

	{"ihn nicht der lassen kann",
		"nicht für Kinder ist der Türkentrank", []int{7, 7}},
}

var testCasesMoved = []TestCase{
	{"Ich ging im Walde so vor mich hin",
		"im Walde Ich ging so vor mich hin", []int{4, 4}},

	{"Ich ging im Walde so vor mich hin",
		"so vor mich hin im Walde Ich ging", []int{8, 8}},
}

func TestLevenshteinA(t *testing.T) {

	cases := &testCasesBasic
	cases = &testCasesAdv1
	cases = &testCasesMoved

	{
		inner(t, cases, 0, ls_core.DefaultOptions, false)
		opt2 := ls_core.DefaultOptions
		opt2.SubCost = 1
		inner(t, cases, 1, opt2, false)
		inner(t, cases, 1, opt2, true)
	}

}

func inner(t *testing.T, cases *[]TestCase, wantIdx int, opt ls_core.Options, sortIt bool) {

	for i, tc := range *cases {

		m := ls_core.New(wrapAsEqualer(tc.src, sortIt), wrapAsEqualer(tc.dst, sortIt), opt)
		got := m.Distance()

		ssrc := fmt.Sprintf("%v", tc.src)
		sdst := fmt.Sprintf("%v", tc.dst)
		if got != tc.distances[wantIdx] {
			t.Logf(
				"%2v: Distance between %20v and %20v should be %v - but got %v ",
				i, pbstrings.Ellipsoider(ssrc, 8), pbstrings.Ellipsoider(sdst, 8), tc.distances[wantIdx], got)
			t.Fail()
		}

		m.Print()
		fmt.Printf("\n")

		es := m.EditScript()

		got2 := m.ApplyEditScript(es)
		if !m.CompareToCol(got2) {
			t.Logf("\nwnt %v \ngot %v ", wrapAsEqualer(tc.dst, sortIt), got2)
			t.Fail()
		}

		fmt.Printf("\n")
		fmt.Printf("\n")

	}

}
