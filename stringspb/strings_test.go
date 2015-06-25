package stringspb

// go test -bench=. -benchtime=6ms

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

// Sadly, conversion of []byte to string costs O(n)
func BenchmarkByteToString(b *testing.B) {

	var tests = []struct {
		inp []byte
		wnt string
	}{
		{[]byte("C A F F E E trink nicht so viel Kaffee"), "C A F F E E trink nicht so viel Kaffee"},
		{[]byte("änne zylinder"), "änne zylinder"},
	}

	b.ResetTimer() // Don't time creation and population

	b.StartTimer()
	idx := 0
	for i := 0; i < b.N; i++ {
		idx = i % len(tests)
		_ = string(tests[idx].inp)
	}

}
