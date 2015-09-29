package dedup

import "golang.org/x/net/html"

func dedupApply(n *html.Node, dedups map[string]bool) {

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		dedupApply(c, dedups)
	}

	if n.Type == html.ElementNode {
		outline := attrX(n.Attr, "ol") + "."

		if dedups[outline] {
			n.Type = html.CommentNode
			n.Data = n.Data + " replaced"
		}
	}

}
