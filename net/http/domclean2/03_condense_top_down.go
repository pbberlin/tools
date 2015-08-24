package domclean2

import (
	"strings"

	"golang.org/x/net/html"
)

func removeEmptyNodes(n *html.Node, lvl int) {

	// children
	cc := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cc = append(cc, c)
	}
	for _, c := range cc {
		removeEmptyNodes(c, lvl+1)
	}

	// processing
	// empty element nodes
	if n.Type == html.ElementNode && n.Data == "img" {
		src := attrX(n.Attr, "src")
		if src == "" {
			n.Parent.RemoveChild(n)
		}
	}

	if n.Type == html.ElementNode && n.FirstChild == nil && n.Data == "a" {
		href := attrX(n.Attr, "href")
		if href == "#" || href == "" {
			n.Parent.RemoveChild(n)
		}
	}

	if n.Type == html.ElementNode && n.FirstChild == nil &&
		(n.Data == "em" || n.Data == "strong") {
		n.Parent.RemoveChild(n)
	}

	if n.Type == html.ElementNode && n.FirstChild == nil &&
		(n.Data == "div" || n.Data == "span" || n.Data == "li" || n.Data == "p") {
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

func condenseTopDown(n *html.Node, lvl, lvlExec int) {

	// like in removeUnwanted, we first assemble children separately.
	// since "NextSibling" might be set to nil during condension
	cc := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cc = append(cc, c)
	}

	for _, c := range cc {
		condenseTopDown(c, lvl+1, lvlExec)
	}

	// position at the end => process from deepest level on upwards
	if lvl == 9 || true {
		topDownV3(n, map[string]bool{"div": true},
			map[string]bool{"div": true, "ul": true, "form": true, "li": true, "p": true,
				"a": true, "span": true})

		topDownV3(n, map[string]bool{"li": true}, map[string]bool{"div": true})

	}

	// condenseTopDown2(n, "li", map[string]bool{"a": true, "div": true}, "li")

}
