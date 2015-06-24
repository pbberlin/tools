package parse2

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

func condenseNestedDivs(n *html.Node, lvl int) {

	// Children
	if true {
		// like in removeDirect, we first assemble children separately.
		// since "NextSibling" might be set to nil during condension
		cc := []*html.Node{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cc = append(cc, c)
		}
		for _, c := range cc {
			condenseNestedDivs(c, lvl+1)
		}
	} else {
		// this also worked,
		// but required 45 (!) repetitive traversals...
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			condenseNestedDivs(c, lvl+1)
		}
	}

	condenseUpwards(n, []string{"div", "div"}, "div")

	// condenseUpwards(n, []string{"div", "ul"}, "ul")  // incompatible with separate assembly, move to separate traversal

	// condenseUpwards(n, []string{"ul", "ul"}, "ul")  // incompatible with separate assembly, move to separate traversal

	// condenseUpwards(n, []string{"li", "div"}, "li") // questionable

}
