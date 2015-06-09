package parse2

import "golang.org/x/net/html"

func condenseUpwards(n *html.Node, couple []string, parentType string) {

	p := n.Parent
	if p == nil {
		return
	}

	parDiv := p.Type == html.ElementNode && p.Data == couple[0] // Parent is a div
	iAmDiv := n.Type == html.ElementNode && n.Data == couple[1] // I am a div

	noSiblings := n.PrevSibling == nil && n.NextSibling == nil

	only1Child := n.FirstChild != nil && n.FirstChild == n.LastChild
	svrlChildn := n.FirstChild != nil && n.FirstChild != n.LastChild
	noChildren := n.FirstChild == nil

	_, _ = noSiblings, noChildren

	if iAmDiv && parDiv {

		if svrlChildn || only1Child {
			var children []*html.Node
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				children = append([]*html.Node{c}, children...) // order inversion
			}

			insertionPoint := n.NextSibling
			for _, c1 := range children {

				n.RemoveChild(c1)

				if c1.Type != html.TextNode {
					p.InsertBefore(c1, insertionPoint)
					insertionPoint = c1
				} else {
					wrap := html.Node{Type: html.ElementNode, Data: "p",
						Attr: []html.Attribute{html.Attribute{Key: "cfrm", Val: "div"}}}
					wrap.FirstChild = c1
					p.InsertBefore(&wrap, insertionPoint)
					insertionPoint = &wrap
				}

			}
			p.RemoveChild(n)
			if p.Data != parentType {
				// fmt.Printf("\t%v=>%v\n", p.Data, parentType)
				p.Data = parentType
			}

		}

	}

}
