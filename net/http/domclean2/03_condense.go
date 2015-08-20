package domclean2

import (
	"strings"

	"golang.org/x/net/html"
)

func convEmptyElementLeafs(n *html.Node, lvl int) {

	// children
	cc := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cc = append(cc, c)
	}
	for _, c := range cc {
		convEmptyElementLeafs(c, lvl+1)
	}

	// processing
	// empty element nodes
	if n.Type == html.ElementNode &&
		n.FirstChild == nil &&
		(n.Data == "div" || n.Data == "span" ||
			n.Data == "li" || n.Data == "p") {
		// n.Type = html.CommentNode
		n.Parent.RemoveChild(n)
	}

	// spans with less than 2 characters inside => flatten to text
	only1Child := n.FirstChild != nil && n.FirstChild == n.LastChild
	if n.Type == html.ElementNode &&
		n.Data == "span" &&
		only1Child &&
		n.FirstChild.Type == html.TextNode &&
		len(strings.TrimSpace(n.FirstChild.Data)) < 3 {
		n.Type = html.TextNode
		n.Data = n.FirstChild.Data
		n.RemoveChild(n.FirstChild)
	}

}

func condenseNestedDivs(n *html.Node, lvl, lvlExec int) {

	// like in removeUnwanted, we first assemble children separately.
	// since "NextSibling" might be set to nil during condension
	cc := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cc = append(cc, c)
	}

	for _, c := range cc {
		condenseNestedDivs(c, lvl+1, lvlExec)
	}

	// position at the end => process from deepest level on upwards
	condenseUpwards2(n, "div", map[string]bool{"div": true, "ul": true, "form": true}, "div")

}
