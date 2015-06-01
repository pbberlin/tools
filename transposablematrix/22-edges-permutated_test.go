// +build edges
// go test -tags=edges

package transposablematrix

import (
	"fmt"
	"sort"
	"testing"

	"github.com/pbberlin/tools/util"
)

func Test_edges(t *testing.T) {

	inputs := []int{5, 6}
	wants := []string{
		"1-1-1; 1-1-3; 1-2-2; 1-3-1; 2-1-1; ",
		"1-1-4; 1-2-1; 1-2-3; 1-3-2; 1-4-1; 2-1-2; ",
	}

	/*
		             XX
				XX   X                   |  XXX
				XX   X    XXXX  XXX      |  X
				X    X    X     XX       |  X
				111  131  113   211      |  122

				111; 113; 131; 211;


				XX     XXXXX  XXXX  XX   |  XXXX  XXX
				X      X      XX    XX   |  X     X
				X                   X    |  X     X
				X                   X    |        X
				X                        |
				141    114    212   121  |  123   132


				Keep in Mind:
				x1;y;x2  equivalent x2;-y;x1

				This is handled in Enc()

	*/

	for i, inp := range inputs {

		ar := NewReservoir()
		ar.GenerateSpecificAmorphs([]int{inp})

		keys := make([]int, 0, len(ar.Edge3))
		for k, _ := range ar.Edge3 {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		got := ""
		for _, k := range keys {
			_, _, _, ks := Dec(k)
			got = fmt.Sprintf("%v%v; ", got, ks)
		}

		if got != wants[i] {
			t.Errorf("\ngot  %v\nwant %v\nfrom %v", got, wants[i], inp)
		}

	}
	fmt.Println("test PermutateEdges")

}

func Test_Factorize(t *testing.T) {

	inputs := []int{4, 5, 6, 8, 12}
	wants := [][]int{[]int{2}, []int{}, []int{2, 3}, []int{2, 4}, []int{2, 3, 4, 6}}

	for i, inp := range inputs {

		got := Factorize(inp)

		if !util.IntsEqualsInts(got, wants[i]) {
			t.Errorf("\ngot  %v\nwant %v\nfrom %v", got, wants[i], inp)
		}

	}
	fmt.Println("test factorize")

}
