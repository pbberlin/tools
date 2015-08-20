package domclean2

import (
	"fmt"

	"golang.org/x/net/html"
)

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

type NodeTypeStr html.NodeType

func (n NodeTypeStr) String() string {
	switch n {
	case 0:
		return "ErrorNode"
	case 1:
		return "TextNode"
	case 2:
		return "DocumentNode"
	case 3:
		return "ElementNode"
	case 4:
		return "CommentNode"
	case 5:
		return "DoctypeNode"
	}
	return "unknown Node type"
}
