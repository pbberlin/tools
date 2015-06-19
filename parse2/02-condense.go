package parse2

import (
	"strings"

	"golang.org/x/net/html"
)

func convEmptyElementLeafs(n *html.Node, lvl int) {

	// processing

	// empty element nodes
	if n.Type == html.ElementNode &&
		n.FirstChild == nil &&
		(n.Data == "div" || n.Data == "span" ||
			n.Data == "li" || n.Data == "p") {
		n.Type = html.CommentNode
	}

	// spans with only 2 characters inside => remove
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

	// children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		convEmptyElementLeafs(c, lvl+1)
	}

}

func condenseNestedDivs(n *html.Node, lvl int) {

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		condenseNestedDivs(c, lvl+1)
	}

	condenseUpwards(n, []string{"div", "div"}, "div")

	// condenseUpwards(n, []string{"div", "ul"}, "ul")

	// condenseUpwards(n, []string{"ul", "ul"}, "ul")

	// condenseUpwards(n, []string{"li", "div"}, "li") // questionable

}
