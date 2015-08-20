package domclean2

import (
	"github.com/pbberlin/tools/net/http/dom"
	"golang.org/x/net/html"
)

// Condense upwards builds a three-levels subtree
// starting from param node l1
// l2 and l3 nodes need to comply by type
//
// Then l3 is moved under l1; l2 is eliminated
//
// For <a> or "text" l3 nodes, we could introduce wrappers
//
// l2Types so far always is "div".
// Multiple l2Types are possible, but difficult to imagine.
//
// l1 type could be changed - from div to ul for instance, but I found no use for that
//
// Implementation yields similar result as condenseUpwards1
// but the "all-or-nothing" logic is clearer
func condenseUpwards2(l1 *html.Node, l2Types map[string]bool, l3Types map[string]bool) {

	if l1.Type != html.ElementNode || l1.Data == "span" || l1.Data == "a" {
		return // cannot assign to textnode ; want not assign to
	}

	// dig two levels deeper

	// isolate l2
	var l2s []*html.Node
	for l2 := l1.FirstChild; l2 != nil; l2 = l2.NextSibling {
		l2s = append(l2s, l2)
		// l2s = append([]*html.Node{l2}, l2s...) // order inversion
	}

	// measure types
	l2Div := true

	// note that *all* l3 must have l3Type, not just those those of one l2 element
	// otherwise we get only partial restructuring - and therefore sequence errors
	l3Div := true

	for _, l2 := range l2s {
		l2Div = l2Div && l2.Type == html.ElementNode && l2Types[l2.Data] // l2 is a div
		for l3 := l2.FirstChild; l3 != nil; l3 = l3.NextSibling {
			l3Div = l3Div && (l3.Type == html.ElementNode && l3Types[l3.Data]) // l3 is a div or ul or form
		}
	}

	// act
	if l2Div && l3Div {
		for _, l2 := range l2s {

			// isolate l3
			var l3s []*html.Node
			for l3 := l2.FirstChild; l3 != nil; l3 = l3.NextSibling {
				l3s = append(l3s, l3)
				// l3s = append([]*html.Node{l3}, l3s...) // order inversion
			}

			// detach l3 from l2
			for _, l3 := range l3s {
				l2.RemoveChild(l3)
			}
			l1.RemoveChild(l2) // detach l2 from l1

			for _, l3 := range l3s {
				// attach l3 to l1, possible wrapper of <a> or <span>
				l1.InsertBefore(l3, nil) // insert at end

				// wrap := html.Node{Type: html.ElementNode, Data: "p", Attr: []html.Attribute{html.Attribute{Key: "cfrm", Val: "div"}}}
				// wrap.FirstChild = c1
				// l1.InsertBefore(&wrap, nil)

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
		b := dom.PrintSubtree(n, nil, 0)
		pf("%s", b)
		return true
	}

	return false

}
