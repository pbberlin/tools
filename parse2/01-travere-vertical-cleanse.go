package parse2

import (
	"strings"

	"golang.org/x/net/html"
)

var (
	removeTypes = map[string]bool{
		"meta":     true,
		"link":     true,
		"style":    true,
		"iframe":   true,
		"script":   true,
		"noscript": true,
		"canvas":   true,
	}
)

func TraverseVertCleanse(n *html.Node, lvl int) {

	RemoveUnwanted(n)

	switch n.Type {
	case html.ElementNode:
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

// RemoveUnwanted neutralizes a node.
// Note: We can not Remove() nor Replace()
// Since that breaks the recursion one step above!
func RemoveUnwanted(n *html.Node) {
	if removeTypes[n.Data] {
		// dom.ReplaceNode(n, &html.Node{Type: html.CommentNode, Data: n.Data + " replaced"})
		// dom.RemoveNode(n)
		n.Type = html.CommentNode
		// fmt.Printf("\tunwanted %9v turned into comment\n", n.Data)
		n.Data = n.Data + " replaced"
	}
}
