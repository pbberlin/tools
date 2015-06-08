package parse2

import (
	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

func TraverseVert_ConvertDivDiv(n *html.Node, lvl int) {

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseVert_ConvertDivDiv(c, lvl+1)
	}

	couple := []string{"div", "div"}
	condenseUpwards(n, couple)

}

func TraverseHori_ConvertDivDiv(lp interface{}, onlyOnLvl int) {

	var queue = util.NewQueue(10)

	for lp != nil {

		lpn := lp.(Tx).Nd
		lvl := lp.(Tx).Lvl

		for c := lpn.FirstChild; c != nil; c = c.NextSibling {
			queue.EnQueue(Tx{c, lvl + 1})
		}

		// processing
		if lvl == onlyOnLvl {
			couple := []string{"div", "div"}
			condenseUpwards(lpn, couple)
		}

		//
		// next node
		lp = queue.DeQueue()
	}
}

func condenseUpwards(n *html.Node, couple []string) {

	p := n.Parent
	if p == nil {
		return
	}

	iAmDiv := n.Type == html.ElementNode && n.Data == couple[0] // I am a div
	parDiv := p.Type == html.ElementNode && p.Data == couple[1] // Parent is a div

	noSiblings := n.PrevSibling == nil && n.NextSibling == nil

	only1Child := n.FirstChild != nil && n.FirstChild == n.LastChild
	svrlChildn := n.FirstChild != nil && n.FirstChild != n.LastChild
	noChildren := n.FirstChild == nil

	_, _ = noSiblings, noChildren

	if iAmDiv && parDiv {

		// fmt.Printf("%2v: %v/div/%v\n", lvl, p.Data, n.FirstChild.Data)

		if svrlChildn || only1Child {
			var children []*html.Node
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				children = append([]*html.Node{c}, children...) // inversion
				// children = append(children, c)
			}

			insertionPoint := n.NextSibling
			for _, c1 := range children {
				n.RemoveChild(c1)
				p.InsertBefore(c1, insertionPoint)
				insertionPoint = c1
			}
			p.RemoveChild(n)
		}

	}

}
