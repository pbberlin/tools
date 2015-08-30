// Package weedout takes multiple dom instances,
// computing similar subtrees measured by levenshtein distance.
package weedout

import "fmt"

var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var pfRestore func(format string, a ...interface{}) (int, error) = fmt.Printf

var spf func(format string, a ...interface{}) string = fmt.Sprintf
var wpf = fmt.Fprintf

func pfDevNull(format string, a ...interface{}) (int, error) {
	return 0, nil // sucking void
}
