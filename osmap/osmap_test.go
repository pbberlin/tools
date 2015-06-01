package osmap

import (
	"strings"
	"testing"
)

func TestStringKeyOMapInsertion(t *testing.T) {
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
		t.Errorf("%q != %q", actual, expected)
	}
}

func TestIntKeyOMapFind(t *testing.T) {
	intMap := New()
	for _, number := range []string{"9", "1", "8", "2", "7", "3", "6", "4", "5", "0"} {

		intMap.Insert(number, "0"+number)
	}
	for _, number := range []string{"0", "1", "5", "8", "9"} {
		value, found := intMap.Find(number)
		if !found {
			t.Errorf("failed to find %d", number)
		}
		actual, expected := value, "0"+number
		if actual != expected {
			t.Errorf("value is %d should be %d", actual, expected)
		}
	}
	for _, number := range []string{"-1", "-21", "10", "11", "148"} {
		_, found := intMap.Find(number)
		if found {
			t.Errorf("should not have found %d", number)
		}
	}
}

func TestIntKeyOMapDelete(t *testing.T) {
	intMap := New()
	for _, number := range []string{"9", "1", "8", "2", "7", "3", "6", "4", "5", "0"} {
		intMap.Insert(number, "0"+number)
	}
	if intMap.Len() != 10 {
		t.Errorf("map len %d should be 10", intMap.Len())
	}
	length := 9
	for i, number := range []string{"0", "1", "5", "8", "9"} {
		if deleted := intMap.Delete(number); !deleted {
			t.Errorf("failed to delete %d", number)
		}
		if intMap.Len() != length-i {
			t.Errorf("map len %d should be %d", intMap.Len(), length-i)
		}
	}
	for _, number := range []string{"-1", "-21", "10", "11", "148"} {
		if deleted := intMap.Delete(number); deleted {
			t.Errorf("should not have deleted nonexistent %d", number)
		}
	}
	if intMap.Len() != 5 {
		t.Errorf("map len %d should be 5", intMap.Len())
	}
}
