package parse2

import "golang.org/x/net/html"

func weedoutApply(weedouts map[string]bool, n *html.Node) {

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		weedoutApply(weedouts, c)
	}

	if n.Type == html.ElementNode {
		outline := attrX(n.Attr, "ol")
		if weedouts[outline] {
			n.Type = html.CommentNode
			n.Data = n.Data + " replaced"
		}
	}

}
