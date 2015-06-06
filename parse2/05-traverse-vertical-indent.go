package parse2

import (
	"strings"

	"github.com/pbberlin/tools/dom"
	"golang.org/x/net/html"
)

func TraverseVertIndent(n *html.Node, lvl int) {

	// Before children processing
	switch n.Type {
	case html.ElementNode:
		if lvl > 2 && n.Parent.Type == html.ElementNode {
			indent := strings.Repeat("\t", lvl-2)
			dom.InsertBefore(n, &html.Node{Type: html.TextNode, Data: "\n" + indent})
		}
	case html.CommentNode:
		dom.InsertBefore(n, &html.Node{Type: html.TextNode, Data: "\n"})
	case html.TextNode:

		// if strings.HasSuffix(n.Data, "\n") {
		// }
		// if strings.HasSuffix(n.Data, " ") {
		// }

		n.Data = strings.TrimSpace(n.Data) + " "
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseVertIndent(c, lvl+1)
	}

	// After children processing
	switch n.Type {
	case html.ElementNode:

		// I dont know why,
		// but this needs to happend AFTER the children
		if lvl > 2 && n.Parent.Type == html.ElementNode {
			indent := strings.Repeat("\t", lvl-2)
			if n.LastChild != nil {
				dom.InsertAfter(n.LastChild, &html.Node{Type: html.TextNode, Data: "\n" + indent})
			}
		}
	}

}
