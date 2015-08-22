package domclean2

import (
	"bytes"

	"github.com/pbberlin/tools/net/http/dom"
	"golang.org/x/net/html"
)

func condenseDetail(n *html.Node, b *bytes.Buffer, depth int) *bytes.Buffer {

	if b == nil {
		b = new(bytes.Buffer)
	}

	switch {
	// case n.Type == html.ElementNode && n.Data == "img":
	// 	wpf(b, fmt.Sprintf("[img] %v %v | ", attrX(n.Attr, "title"), attrX(n.Attr, "src")))
	case n.Type == html.ElementNode && n.Data == "a":
		wpf(b, "[a] ")
	case n.Type == html.TextNode && n.Data != "":
		if len(n.Data) < 4 {
			wpf(b, n.Data)
		} else {
			wpf(b, "[txt l%v] %v", depth, n.Data)
		}
	default:
		if n.Type == html.ElementNode {
			wpf(b, " [%v] ", n.Data)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		condenseDetail(c, b, depth+1)
	}

	return b
}

func condenseBottomUp(n *html.Node) {

	switch {

	case n.Type == html.ElementNode && (n.Data == "a" || n.Data == "img"):

		b := condenseDetail(n, nil, 0)
		nodeRepl := dom.Nd("text", b.String())

		// We want to remove all existing children.
		// Direct loop impossible, since "NextSibling" is set to nil by Remove().
		// Therefore first assembling separately, then removing.
		children := make(map[*html.Node]struct{})
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			children[c] = struct{}{}
		}
		for k, _ := range children {
			n.RemoveChild(k)
		}

		// we can't put our replacement "under" an image, since img cannot have children
		if n.Type == html.ElementNode && n.Data == "img" {
			dom.InsertAfter(n, nodeRepl)
			dom.RemoveNode(n)

		} else {
			n.AppendChild(nodeRepl)
		}

		// Insert a  || and a newline before every <a...>
		if n.Data == "a" {
			prev := n

			breaker0 := dom.Nd("text", " || ")
			n.Parent.InsertBefore(breaker0, prev)

			// breaker1 := dom.Nd("br")
			// n.Parent.InsertBefore(breaker1, prev)

			// breaker2 := dom.Nd("text", "\n")
			// n.Parent.InsertBefore(breaker2, prev)
		}

	default:
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		condenseBottomUp(c)
	}
}
