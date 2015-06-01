//x +build stairy
// go test -tags=stairy

package transposablematrix

import (
	"fmt"
	"sort"
	"testing"
)

func Test_stairyEdge(t *testing.T) {

	inputs := []int{5, 6, 5, 7, 8}
	ar := NewReservoir()
	ar.GenerateSpecificAmorphs(inputs)

	for i, _ := range ar.Amorphs {
		ar.Amorphs[i].AestheticValue = 0
	}

	// pf = fmt.Printf

	// ar.PrintAmorphSummary("Amorph Summary")
	// DumpAmorphs(ar.Amorphs, 2)

	// PrintEdgesRestrictTo("edge3", ar.Edge3, map[int]bool{})
	// PrintEdgesRestrictTo("no-sl", ar.EdgesSlackless, map[int]bool{})
	// PrintEdgesRestrictTo("numel", ar.MElements, map[int]bool{})

	want1 := [][]Amorph{
		[]Amorph{ar.Amorphs[1]},
		[]Amorph{ar.Amorphs[3]},
		[]Amorph{ar.Amorphs[4]},
	}

	want2 := [][]Amorph{
		[]Amorph{ar.Amorphs[1]},
		[]Amorph{ar.Amorphs[0], ar.Amorphs[2], ar.Amorphs[4]},
	}

	want3 := [][]Amorph{
		[]Amorph{ar.Amorphs[1]},
		[]Amorph{ar.Amorphs[3]},
		[]Amorph{ar.Amorphs[0], ar.Amorphs[2], ar.Amorphs[4]},
		[]Amorph{ar.Amorphs[4]},
	}

	wants := [][][]Amorph{want1, want2, want3}
	shrinks := []GrowShrink{grow, shrink, growShrink}

	for i := 0; i < 3; i++ {

		got, _ := ar.ByStairyEdge(2, 1, 2, 2, eastw, shrinks[i])
		if !cmpareSl(got, wants[i]) {
			t.Errorf("\ngot  %v\nwant %v\n", got, wants[i])
		}
		// for i := 0; i < len(got); i++ {
		// 	DumpAmorphs(got[i], 2)
		// }
	}
	fmt.Println("test ByStairyEdge")

}

func cmpareSl(a, b [][]Amorph) bool {
	// return reflect.DeepEqual(a, b)

	if len(a) != len(b) {
		fmt.Println("disticnt1")
		return false
	}

	for i := 0; i < len(a); i++ {
		if len(a[i]) != len(b[i]) {
			fmt.Println("disticnt2")
			return false
		}

		ai, bi := make([]int, len(a[i])), make([]int, len(b[i]))
		for j := 0; j < len(a[i]); j++ {
			ai[j] = a[i][j].IdxArticle
			bi[j] = b[i][j].IdxArticle
		}

		sort.Ints(ai)
		sort.Ints(bi)

		as := fmt.Sprintf("%v", ai)
		bs := fmt.Sprintf("%v", bi)

		if as != bs {
			fmt.Println("disticnt3")
			return false
		}

	}
	return true
}
