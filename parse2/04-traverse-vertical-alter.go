package parse2

import (
	"fmt"

	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

// commonly used by recursive function calls:
var (
	xPath     util.Stack
	xPathSkip = map[string]bool{"em": true, "b": true, "br": true}
	xPathDump []byte

	removeAttributes = map[string]bool{
		"style": true,
		"class": true,
		// "alt":                  true,
		// "title":                  true,
		"target":                 true,
		"id":                     true,
		"rel":                    true,
		"headline":               true,
		"onload":                 true,
		"onclick":                true,
		"onmousedown":            true,
		"onerror":                true,
		"readonly":               true,
		"itemprop":               true,
		"itemtype":               true,
		"itemscope":              true,
		"datetime":               true,
		"current-time":           true,
		"fb-iframe-plugin-query": true,
		"fb-xfbml-state":         true,
	}
	attrDistinct = map[string]int{}
)

func TraverseVert(n *html.Node, lvl int) {

	if lvl == 0 {
		xPathDump = []byte{}
	}

	// Before children processing
	switch n.Type {
	case html.ElementNode:

		if !xPathSkip[n.Data] {
			xPath.Push(n.Data)

			// n.Attr = addIdAttr(n.Attr)
			// printAttr(n.Attr, []string{"xxid"})

			// lvl == xPath.Len()
			s := fmt.Sprintf("%2v: %s\n", xPath.Len(), xPath.StringExt(true))
			xPathDump = append(xPathDump, s...) // special comfort; http://stackoverflow.com/questions/16248241/concatenate-two-slices-in-go#

		}
		n.Attr = removeAttr(n.Attr, removeAttributes)

	case html.TextNode:
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseVert(c, lvl+1)
	}

	// After children processing
	switch n.Type {
	case html.ElementNode:
		if !xPathSkip[n.Data] {
			xPath.Pop()
		}
	}

}
