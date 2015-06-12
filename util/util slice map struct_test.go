package util

import (
	"fmt"
	"testing"
)

func Test_Intslice2Int(t *testing.T) {

	inputs := [][]int{
		[]int{1},
		[]int{1, 2},
		[]int{2, 1},
		[]int{1, 3, 3},
		[]int{3, 1, 3},
		[]int{0, 3, 0},
		[]int{4, 3, 1},
	}
	wants := []int{
		1,
		201,
		201,
		30301,
		30301,
		30000,
		40301,
	}

	for i, inp := range inputs {
		got := Intslice2Int(inp)
		if got != wants[i] {
			t.Errorf("got %v - want %v - from %v", got, wants[i], inp)
		}

	}
}

func Test_InsertAfter(t *testing.T) {

	{
		inp := []int{1, 2, 3, 4}
		got := InsertAfter(inp, 2, 17)
		want := []int{1, 2, 3, 17, 4}

		if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", want) {
			t.Errorf("got %v - want %v", got, want)
		}

	}

	{
		inp := []int{1, 2, 3, 4}
		got := InsertAfter(inp, 3, 17)
		want := []int{1, 2, 3, 4, 17}

		if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", want) {
			t.Errorf("got %v - want %v", got, want)
		}

	}

	{
		inp := []int{1, 2, 3, 4}
		got := InsertAfter(inp, -1, 17)
		want := []int{17, 1, 2, 3, 4}

		if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", want) {
			t.Errorf("got %v - want %v", got, want)
		}

	}

}

func Test_Delete(t *testing.T) {

	{
		inp := []int{1, 2, 3, 4}
		got := Delete(inp, 2)
		want := []int{1, 2, 4}

		if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", want) {
			t.Errorf("got %v - want %v", got, want)
		}

	}

	{
		inp := []int{1, 2, 3, 4}
		got := Delete(inp, 0)
		want := []int{2, 3, 4}

		if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", want) {
			t.Errorf("got %v - want %v", got, want)
		}

	}

	{
		inp := []int{1, 2, 3, 4}
		got := Delete(inp, 3)
		want := []int{1, 2, 3}

		if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", want) {
			t.Errorf("got %v - want %v", got, want)
		}

	}
}
