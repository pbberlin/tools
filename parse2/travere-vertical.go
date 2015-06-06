package parse2

import (
	"fmt"

	"github.com/pbberlin/tools/dom"
	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

// commonly used by recursive function calls:
var (
	xPath     util.Stack
	xPathSkip = map[string]bool{"em": true, "b": true, "br": true}
	xPathDump []byte
)

func TraverseVert(n *html.Node, lvl int) {

	// Before children processing
	switch n.Type {
	case html.ElementNode:
		switch n.Data {
		case "meta":
			return
		case "link", "style":
			return
		case "iframe", "script", "noscript":
			return
		}

		if !xPathSkip[n.Data] {
			xPath.Push(n.Data)

			// n.Attr = addIdAttr(n.Attr)
			// printAttr(n.Attr, []string{"xxid"})

			// lvl == xPath.Len()
			s := fmt.Sprintf("%2v: %s\n", xPath.Len(), xPath.StringExt(true))
			xPathDump = append(xPathDump, s...) // special comfort; http://stackoverflow.com/questions/16248241/concatenate-two-slices-in-go#

			// Children
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if isSpacey(c.Data) {
					fmt.Printf("spacey tn %q\n", c.Data)
					dom.RemoveNode(c)
				}
			}

		}

	case html.TextNode:
		// if xPath.Last() == "body" {
		// 	fmt.Printf("%v \t", xPath.String())
		// 	fmt.Printf("%q\n", n.Data)
		// 	n.Data = strings.TrimSpace(n.Data)
		// }
		// n.Data = strings.TrimSpace(n.Data)

		// fmt.Printf("textnode %q\n", n.Data)
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
