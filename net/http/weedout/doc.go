// Package weedout takes multiple dom instances,
// computing similar subtrees measured by levenshtein distance.
package weedout

import (
	"fmt"

	"golang.org/x/net/html"
)

var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var pfRestore func(format string, a ...interface{}) (int, error) = fmt.Printf

var spf func(format string, a ...interface{}) string = fmt.Sprintf
var wpf = fmt.Fprintf

func pfDevNull(format string, a ...interface{}) (int, error) {
	return 0, nil // sucking void
}

// !DOCTYPE html head
// !DOCTYPE html body
//        0    1    2
const cScaffoldLvls = 2

func attrX(attributes []html.Attribute, key string) (s string) {
	for _, a := range attributes {
		if key == a.Key {
			s = a.Val
			break
		}
	}
	return
}
