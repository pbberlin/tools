package pbstrings

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {

	var tests = []struct {
		Input string
		Want  []string
		//Output []string
	}{
		{"a b", []string{"a", "b"}},
		{"änne zylinder", []string{"änne", "zylinder"}},
		{"änne    zylinder", []string{"änne", "zylinder"}},
	}

	for i0, v := range tests {
		gots := SplitByWhitespace(v.Input)

		var msg string = fmt.Sprintf(
			"Test Nr. %q split(%q) == %q - want %q",
			i0, v.Input, gots, v.Want)

		for i1, _ := range gots {
			got := gots[i1]
			if got == v.Want[i1] {
				// fmt.Printf("fine - found %q as expected\n", got)
			} else {
				fmt.Println("ERROR", msg)
				t.Error("got %q - want %q", got, v.Want[i1])
			}

		}

	}
}
