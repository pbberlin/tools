package parse2

import (
	"fmt"

	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

// vars used by all recursive function calls:
var (
	xPath     util.Stack
	xPathSkip = map[string]bool{"em": true, "b": true, "br": true}
	xPathDump []byte
)

// TraverseVertAlter1 writes an xpath log.
// TraverseVertAlter1 cleans up the attributes
func TraverseVertAlter1(n *html.Node, lvl int) {

	if lvl == 0 {
		xPathDump = []byte{}
	}

	// Before children processing
	switch n.Type {
	case html.ElementNode:

		nodeDistinct[n.Data]++

		if !xPathSkip[n.Data] {
			xPath.Push(n.Data)

			// lvl == xPath.Len()
			s := fmt.Sprintf("%2v: %s\n", xPath.Len(), xPath.StringExt(true))
			xPathDump = append(xPathDump, s...) // yes, string appends to byteSlice ; http://stackoverflow.com/questions/16248241/concatenate-two-slices-in-go#

		}
	}
	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseVertAlter1(c, lvl+1)
	}

	// After children processing
	switch n.Type {
	case html.ElementNode:
		if !xPathSkip[n.Data] {
			xPath.Pop()
		}
	}

}
