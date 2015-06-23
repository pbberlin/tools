package parse2

import "fmt"

var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var pfRestore func(format string, a ...interface{}) (int, error) = fmt.Printf
var spf func(format string, a ...interface{}) string = fmt.Sprintf

// type pft func(format string, a ...interface{}) (int, error)

func pfDevNull(format string, a ...interface{}) (int, error) {
	return 0, nil // sucking void
}

// Disable all printing
// within a function and callees:
func exampleUsage() {
	pf = pfDevNull
	defer func() { pf = pfRestore }()
}
