package domclean2

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

// dumpXPath writes an xpath log.
// dumpXPath also collects frequency of node type data.
func dumpXPath(n *html.Node, lvl int) {

	if lvl == 0 {
		xPathDump = []byte{}
	}

	// Before children processing
	switch n.Type {
	case html.ElementNode:

		nodeDistinct[n.Data]++ // Histogram

		if !xPathSkip[n.Data] {
			xPath.Push(n.Data)
			// apart from skipped node types:  xPath.Len() == lvl
			s := fmt.Sprintf("%2v: %s\n", xPath.Len(), xPath.StringExt(false))
			xPathDump = append(xPathDump, s...) // yes, string appends to byteSlice ; http://stackoverflow.com/questions/16248241/concatenate-two-slices-in-go#
		}
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		dumpXPath(c, lvl+1)
	}

	// After children processing
	switch n.Type {
	case html.ElementNode:
		if !xPathSkip[n.Data] {
			xPath.Pop()
		}
	}

}
