// Package parse2 normalizes html dom trees; structure and formatting are simplified.
// Lots of iterations/traversals are required before all nested structs are flattened.
package parse2

import "fmt"

var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var spf func(format string, a ...interface{}) string = fmt.Sprintf
