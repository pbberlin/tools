package parse2

import (
	"fmt"

	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

var ndStack util.Stack // filled by TraverseVert()
var skipInStack = map[string]bool{"em": true}
var stackOutp []byte

func TraverseVert(n *html.Node, lvl int) {

	// Before children
	switch n.Type {
	case html.ElementNode:

		if !skipInStack[n.Data] {
			ndStack.Push(n.Data)
			n.Attr = addIdAttr(n.Attr)
			printAttr(n.Attr, []string{"id", "bd"})
		}

		switch n.Data {
		case "a":

		case "iframe", "script", "noscript":
			return
		}

		// lvl == ndStack.Len()
		s := fmt.Sprintf("%2v: %s\n", ndStack.Len(), ndStack.StringExt(true))
		stackOutp = append(stackOutp, s...) // exceptional comfort case

	case html.TextNode:
		//
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseVert(c, lvl+1)
	}

	// After children
	switch n.Type {
	case html.ElementNode:

		if !skipInStack[n.Data] {
			ndStack.Pop()
		}
	}

}
