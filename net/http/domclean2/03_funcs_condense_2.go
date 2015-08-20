package domclean2

import (
	"github.com/pbberlin/tools/net/http/dom"
	"golang.org/x/net/html"
)

/*
	<div>
		<div>
			<p>paragr 1</p>
			<p>paragr 2</p>
		</div>
		<div>
			simple line
		</div>
		<div>
			simple line
		</div>
	</div>


		<div>
			<p>paragr 1</p>
			<p>paragr 2</p>
		</div>
		<div>
			simple line
		</div>
		<div>
			simple line
		</div>


*/
func condenseUpwards2(l1 *html.Node, l2Type string, l3Types map[string]bool, newL2Type string) {

	if l1.Type != html.ElementNode {
		return // cannot assign to textnode
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
	l3Div := true // all l3 must have l3Type
	for _, l2 := range l2s {
		l2Div = l2Div && l2.Type == html.ElementNode && l2.Data == l2Type // l2 is a div
		for l3 := l2.FirstChild; l3 != nil; l3 = l3.NextSibling {
			l3Div = l3Div && (l3.Type == html.ElementNode && l3Types[l3.Data]) // l3 is a div or ul or form
		}
	}

	// act
	if l2Div && l3Div {
		for _, l2 := range l2s {

			var l3s []*html.Node
			for l3 := l2.FirstChild; l3 != nil; l3 = l3.NextSibling {
				l3s = append(l3s, l3)
				// l3s = append([]*html.Node{l3}, l3s...) // order inversion
			}

			for _, l3 := range l3s {
				l2.RemoveChild(l3)
			}
			l1.RemoveChild(l2)
			for _, l3 := range l3s {
				l1.InsertBefore(l3, nil) // insert at end
			}
			l2.Data = newL2Type
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
