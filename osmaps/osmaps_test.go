// +build integration
// go test -tags=integration

package osmaps

import (
	"strings"
	"testing"
)

func DISABLED_TestStringKeyOMapInsertion(t *testing.T) {
	wordForWord := New()
	for _, word := range []string{"one", "Two", "THREE", "four", "Five"} {
		wordForWord.Insert(strings.ToLower(word), word)
	}
	var words []string
	fc := func(_, value string) {
		words = append(words, value)
	}
	wordForWord.Do(fc)
	actual, expected := strings.Join(words, ""), "FivefouroneTHREETwo"
	if actual != expected {
		t.Errorf("got %q != wanted %q", actual, expected)
	}
}

func TestInsertUpdateFind(t *testing.T) {
	map1 := New()
	const keyPref = "a100"
	const valPref = "val"

	keys := []string{"9", "1", "8", "2", "7", "3", "6", "4", "5", "0"}

	for _, number := range keys {
		map1.Insert(keyPref+number, number)
	}

	// UPDATE
	for _, number := range keys {
		map1.Insert(keyPref+number, valPref+number)
	}

	for _, number := range []string{"0", "1", "5", "8", "9"} {
		value, found := map1.Find(keyPref + number)
		if !found {
			t.Errorf("failed to find %v", keyPref+number)
		}
		actual, expected := value, valPref+number
		if actual != expected {
			t.Errorf("value is %v should be %v", actual, expected)
		}
	}
	for _, number := range []string{"-1", "-21", "10", "11", "148"} {
		_, found := map1.Find(keyPref + number)
		if found {
			t.Errorf("should not have found %v", keyPref+number)
		}
	}
}

func DISABLED_IntKeyOMapDelete(t *testing.T) {
	map1 := New()
	for _, number := range []string{"9", "1", "8", "2", "7", "3", "6", "4", "5", "0"} {
		map1.Insert(number, "0"+number)
	}
	if map1.Len() != 10 {
		t.Errorf("map len %d should be 10", map1.Len())
	}
	length := 9
	for i, number := range []string{"0", "1", "5", "8", "9"} {
		if deleted := map1.Delete(number); !deleted {
			t.Errorf("failed to delete %d", number)
		}
		if map1.Len() != length-i {
			t.Errorf("map len %d should be %d", map1.Len(), length-i)
		}
	}
	for _, number := range []string{"-1", "-21", "10", "11", "148"} {
		if deleted := map1.Delete(number); deleted {
			t.Errorf("should not have deleted nonexistent %d", number)
		}
	}
	if map1.Len() != 5 {
		t.Errorf("map len %d should be 5", map1.Len())
	}
}
