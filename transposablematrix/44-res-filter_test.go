// +build abundant
// go test -tags=abundant

package transposablematrix

import (
	"fmt"
	"testing"
)

func Test_abundantHeightMatch(t *testing.T) {

	inputs := []int{}

	pf = fmt.Printf

	inputs = []int{2, 8, 12, 13, 14}
	ar1 := NewReservoir()
	ar1.GenerateSpecificAmorphs(inputs)

	inputs = []int{5, 6, 7, 8, 15, 14}
	ar2 := NewReservoir()
	ar2.GenerateSpecificAmorphs(inputs)
	for i := 0; i < len(ar2.Amorphs); i++ {
		ar2.Amorphs[i].IdxArticle += 10
	}

	amorphBlocks := [][]Amorph{ar1.Amorphs, ar2.Amorphs}

	f := NewFusion()
	f.w = append(f.w, 2)
	f.e = append(f.e, 2)
	f.xyx = append(f.xyx, 2, 1, 2)
	f.pm = append(f.pm, 4, 4, 6)

	chos := abundantHeightMatch(amorphBlocks, f)

	if chos.Rows != 3 {
		t.Errorf("\ngot  %v\nwant %v\n", chos.Rows, 3)
	}
	// fmt.Printf("ArtIdx%v: H%vx%v Els%v\n", chos.IdxArticle, chos.Rows, chos.Cols, chos.NElements)

	fmt.Println("Test_abundantHeightMatch")
}

func Test_SwitchHoriz(t *testing.T) {

	w := westw
	w.SwitchHoriz()
	if w != eastw {
		t.Errorf("\ngot  %v\nwant %v\n", w, eastw)
	}

	w1 := eastw
	w1.SwitchHoriz()
	if w1 != westw {
		t.Errorf("\ngot  %v\nwant %v\n", w1, westw)
	}

}
