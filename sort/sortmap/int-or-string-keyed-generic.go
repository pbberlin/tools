// Package sortmap sorts map[string]int by int; as required in histograms;
// it also contains useless map[string]interface{} and map[int]interface{} sorting.
package sortmap

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
// But we do not need this generify stuff anyway.
// Just inline it:
//		keys := make([]int, 0, len(mapAny))
//		for k := range mapAny {
//			keys = append(keys, k)
//		}
//		sort.Ints(keys) // or Strings...
//		for _, k := range keys {
//			doSomething  := mapAny[k]
//		}
func StringKeysToSortedArray(m map[string]interface{}) (keys []string) {
	keys = make([]string, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}

// IntKeysToSortedArray is like StringKeysToSortedArray, but for int
func IntKeysToSortedArray(m map[int]interface{}) (keys []int) {
	keys = make([]int, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	return
}

// MapGenerifyType32168 converts some special map into a map[string]interface{}.
// The function argument is too specific.
// We need to rewrite it for every other slice type.
func MapGenerifyType32168(m map[int][]int) (ret map[string]interface{}) {
	ret = make(map[string]interface{})
	for k, v := range m {
		ret[fmt.Sprintf("%02v", k)] = v
	}
	return ret
}

// IntegerFieldToSortedArray tries a generic approach to extract
// an int out of any slice type.
// I dont like it because of the O(n) reflecion.
// Also: The argument slice once again needs to be created first.
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
