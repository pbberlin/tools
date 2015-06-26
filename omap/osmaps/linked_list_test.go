// +build ll
// go test -tags=ll

package osmaps

import "testing"

func TestInsertUpdateFind(t *testing.T) {

	nd := &node{}
	nd.kv = make([]kvt, 0, 100)
	fct1 := func(a, b string) bool {
		return a < b
	}

	keys := []string{"c", "a", "f", "e", "T", "r", "i", "n", "k", "i", "t"}
	keys = []string{"c", "a", "f", "e", "b", "T", "Q", "L"}
	keys = []string{"ca", "ff", "ee", "tr", "in", "k ", "ni", "ch", "t ", "so", "vi", "el"}

	for _, key := range keys {
		insertLinkedList(nd, fct1, key, "val_of_"+key)
	}

}
