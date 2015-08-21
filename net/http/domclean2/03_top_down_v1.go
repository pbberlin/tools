package domclean2

import (
	"github.com/pbberlin/tools/net/http/dom"
	"golang.org/x/net/html"
)

/*
   div                     div
       div                     p
           p         TO        img
           img                 p
           p


	Operates from the *middle* div.
	Saves all children in inverted slice.
	Removes each child and reattaches it one level higher.
	Finally the intermediary, now childless div is removed.




   \                  /
    \       /\       /
     \_____/  \_____/

     \              /
      \_____/\_____/

       \__________/     => Breaks are gone


       \p1___p2___/     => Wrapping preserves breaks




*/
func topDownV1(n *html.Node, couple []string, parentType string) {

	if noParent(n) {
		return
	}
	p := n.Parent

	parDiv := p.Type == html.ElementNode && p.Data == couple[0] // Parent is a div
	iAmDiv := n.Type == html.ElementNode && n.Data == couple[1] // I am a div

	noSiblings := n.PrevSibling == nil && n.NextSibling == nil

	only1Child := n.FirstChild != nil && n.FirstChild == n.LastChild
	svrlChildn := n.FirstChild != nil && n.FirstChild != n.LastChild
	noChildren := n.FirstChild == nil

	_, _ = noSiblings, noChildren

	if parDiv && iAmDiv {

		if only1Child || svrlChildn {

			var children []*html.Node
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				children = append([]*html.Node{c}, children...) // order inversion
			}

			insertionPoint := n.NextSibling
			for _, c1 := range children {

				n.RemoveChild(c1)

				if c1.Type == html.TextNode || c1.Data == "a" {
					// pf("wrapping %v\n", NodeTypeStr(c1.Type))
					wrap := html.Node{Type: html.ElementNode, Data: "p",
						Attr: []html.Attribute{html.Attribute{Key: "cfrm", Val: "div"}}}
					wrap.FirstChild = c1
					p.InsertBefore(&wrap, insertionPoint)
					c1.Parent = &wrap
					insertionPoint = &wrap

				} else {
					p.InsertBefore(c1, insertionPoint)
					insertionPoint = c1
				}

			}
			p.RemoveChild(n)
			if p.Data != parentType {
				p.Data = parentType
			}

		}

	}

}

func noParent(n *html.Node) bool {

	p := n.Parent
	if p == nil {
		if n.Type == html.DoctypeNode || n.Type == html.DocumentNode {
			return true
		}
		pf("parent is nil\n")
		b := dom.PrintSubtree(n)
		pf("%s", b)
		return true
	}

	return false

}
