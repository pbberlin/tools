package util

import (
	"fmt"
	"reflect"
	"sort"
)

// StringKeysToSortedArray() is a helper to print maps in a sorted key order.
// We extract the keys and sort those into an array.
// For printing we loop over the !array! instead over the map.
// Deplorably, our map values need type interface
//
// But we do not need this generify stuff anyway!
// Simply wanna sort a map?
//		sKeys := make([]int, 0, len(mapAny))
//		for k := range mapAny {
//			sKeys = append(sKeys, k)
//		}
//		sort.Ints(sKeys)
//		for _, k := range sKeys {
//			doSomething  := mapAny[k]
//		}
func StringKeysToSortedArray(m map[string]interface{}) (vKeys []string) {
	vKeys = make([]string, len(m))
	i := 0
	for key, _ := range m {
		vKeys[i] = key
		i++
	}
	sort.Strings(vKeys)
	return
}

// IntKeysToSortedArray is like StringKeysToSortedArray, but more idiomatic.
// It too is inapplicable to any other typed map
func IntKeysToSortedArray(m map[int]interface{}) (keys []int) {
	keys = make([]int, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	return
}

// only for debug
func MapGenerifyType32168(m map[int][]int) (ret map[string]interface{}) {
	ret = make(map[string]interface{})
	for k, v := range m {
		ret[fmt.Sprintf("%02v", k)] = v
	}
	return ret
}

// generic approach - I dont like it because of the O(n) reflecion
func IntegerFieldToSortedArray(m []interface{}, fieldName string) (vKeys []int) {
	if len(m) > 1000 {
		panic("this uses reflection - not for large structs")
	}
	vKeys = make([]int, len(m))
	for i, iface := range m {
		vKeys[i] = GetIntField(iface, fieldName)
	}
	sort.Ints(vKeys)
	return
}

// fStudy is based on this
// http://stackoverflow.com/questions/6395076/in-golang-using-reflect-how-do-you-set-the-value-of-a-struct-field
func GetIntField(myStruct interface{}, fieldName string) (ret int) {

	ps := reflect.ValueOf(&myStruct) // pointer to a struct => addressable
	s := ps.Elem()
	if s.Kind() == reflect.Struct {
		// exported field
		f := s.FieldByName(fieldName)
		if f.IsValid() {
			// A Value can be changed only if it is
			// addressable and was not obtained by
			// the use of unexported struct fields.
			if f.CanSet() { // instead of CanAdr
				if f.Kind() == reflect.Int {
					// the "set" case is incommented:
					// x := int64(7)
					// if !f.OverflowInt(x) {
					// 	f.SetInt(x)
					// }
					ret = int(f.Int())
				} else {
					panic("not an int")
				}
			} else {
				panic("field can not set")
			}
		} else {
			panic("field invalid")
		}
	} else {
		panic("not a struct")
	}
	return
}

// DoubleSliceCap doubles the capacity of a slice
func DoubleSliceCap(ba []byte) (newBa []byte) {
	newBa = make([]byte, len(ba), 2*cap(ba))
	copy(newBa, ba)
	fmt.Printf("len: %v -  newcap: %v <br>\n", len(newBa), cap(newBa))
	return newBa
}

// https://github.com/golang/go/wiki/SliceTricks
func InsertAfter(s []int, idx int, newVal int) []int {
	if idx > len(s)-1 {
		panic("Cannot insert beyond existing length")
	}
	s = append(s, 0)
	copy(s[idx+2:], s[idx+1:])
	s[idx+1] = newVal
	return s
}

//
// When having multiple URL params of the same name,
func StringSliceToMapKeys(s []string) map[string]bool {
	ret := map[string]bool{}
	for _, v := range s {
		ret[v] = true
	}
	return ret
}

// Intslice2Int condenses several integers into a single int.
// Typically to use it as a map key.
// Intslice2Int unifies permutations.
// []int{0,3,1} and []int{1,3,0} are sorted to []int{0,1,3} and mapped to 30100
// []int{4,3,1} is sorted to []int{1,3,4} and mapped to 40301
func Intslice2Int(arg []int) (rv int) {
	if len(arg) > 10 {
		panic("only for small sets!")
	}
	sorted := make([]int, len(arg))
	copy(sorted, arg)
	sort.Ints(sorted)
	fact := 1
	for i := 0; i < len(sorted); i++ {
		rv += sorted[i] * fact
		fact = fact * 100
	}
	return
}

func IntsEqualsInts(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}

	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
