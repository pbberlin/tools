package word

import (
	"fmt"

	"github.com/pbberlin/tools/stringspb"
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

	{"Ich ging im Walde so vor mich hin",
		"Ich ging im Forst so vor mich her", []int{4, 228}},

	{"Ich ging im Walde so vor mich hin",
		"Ich ging im Walde so vor mich her hin", []int{1, 228}},
}

func TestLevenshteinA(t *testing.T) {

	cases := testCasesBasic
	cases = append(cases, testCasesAdv1...)
	cases = append(cases, testCasesMoved...)
	// cases = testCasesMoved[3:4]

	{
		inner(t, &cases, 0, ls_core.DefaultOptions, false)

		// opt2 := ls_core.DefaultOptions
		// opt2.SubCost = 1
		// inner(t, &cases, 1, opt2, false)
		// inner(t, &cases, 1, opt2, true)
	}

}

func inner(t *testing.T, cases *[]TestCase, wantIdx int, opt ls_core.Options, sortIt bool) {

	for i, tc := range *cases {

		m := ls_core.New(WrapAsEqualer(tc.src, sortIt), WrapAsEqualer(tc.dst, sortIt), opt)
		got, relDist := m.Distance()
		_ = relDist
		// fmt.Printf("%v %v\n", got, relDist)

		ssrc := fmt.Sprintf("%v", tc.src)
		sdst := fmt.Sprintf("%v", tc.dst)
		if got != tc.distances[wantIdx] {
			t.Logf(
				"%2v: Distance between %20v and %20v should be %v - but got %v (sorted %v)",
				i, stringspb.Ellipsoider(ssrc, 8), stringspb.Ellipsoider(sdst, 8), tc.distances[wantIdx], got, sortIt)
			t.Fail()
		}

		m.Print()
		fmt.Printf("\n")

		es := m.EditScript()

		got2 := m.ApplyEditScript(es)
		if !m.CompareToCol(got2) {
			t.Logf("\nwnt %v \ngot %v ", WrapAsEqualer(tc.dst, sortIt), got2)
			t.Fail()
		}

		fmt.Printf("\n")
		fmt.Printf("\n")

	}

}
