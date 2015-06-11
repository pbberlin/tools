package levenshtein

import "fmt"

// Requires type of sl == interface{}
// Would double conversion cost.
// => We have to convert in the calling package.
func ConvertToEqualer(sl []interface{}) []Equaler {
	var ret = make([]Equaler, 0, len(sl))
	for _, v := range sl {
		cnv, ok := v.(Equaler)
		if !ok {
			panic(fmt.Sprintf("%v %T is not convertible to Equaler interface", v, v))
		}
		ret = append(ret, cnv)
	}
	return ret
}
