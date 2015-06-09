package parse2

import "golang.org/x/net/html"

func TravVertTextify(n *html.Node, lvl, horiNum int) (b []byte) {

	if n.Type == html.TextNode {
		b = append(b, n.Data...)
	}

	// Children

	for c := n.FirstChild; c != nil; c = c.NextSibling {

		if c.Type == html.TextNode {
			// Textnodes have no children
			b = append(b, n.Data...)
		} else {
			horiNum++
			b1 := TravVertTextify(c, lvl+1, horiNum)
			b = append(b, b1...)
		}

	}

	return
}
