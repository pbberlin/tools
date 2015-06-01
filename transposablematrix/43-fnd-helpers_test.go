// +build abundant
// go test -tags=abundant
package transposablematrix

import (
	"fmt"
	"testing"
)

func Test_abundantHeightMatch(t *testing.T) {

	inputs := []int{5, 6, 5, 7, 8}
	ar := NewReservoir()
	ar.GenerateSpecificAmorphs(inputs)

	for i, _ := range ar.Amorphs {
		ar.Amorphs[i].AestheticValue = 0
	}

	fmt.Println("Test_abundantHeightMatch")
}
