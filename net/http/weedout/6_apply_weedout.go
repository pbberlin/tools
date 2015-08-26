package weedout

import "golang.org/x/net/html"

func weedoutApply(n *html.Node, weedouts map[string]bool) {

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		weedoutApply(c, weedouts)
	}

	if n.Type == html.ElementNode {
		outline := attrX(n.Attr, "ol") + "."

		if weedouts[outline] {
			n.Type = html.CommentNode
			n.Data = n.Data + " replaced"
		}
	}

}
