package parse2

import (
	"strings"

	"golang.org/x/net/html"
)

func TraverseVertCleanse(n *html.Node, lvl int) {

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
		cleanseSpaceyTextNodes(n)

	case html.TextNode:
		n.Data = strings.TrimSpace(n.Data) + " "
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseVertCleanse(c, lvl+1)
	}

	// After children processing
	switch n.Type {
	case html.ElementNode:
	}

}
