package parse2

import (
	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

func TraverseVert_CondenseDivStaples(n *html.Node, lvl int) {

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseVert_CondenseDivStaples(c, lvl+1)
	}

	condenseUpwards(n, []string{"div", "div"}, "div")

	condenseUpwards(n, []string{"div", "ul"}, "ul")

	condenseUpwards(n, []string{"ul", "ul"}, "ul")

}

//
func UNUSED_TraverseHori_CondenseDivStaples(lp interface{}, onlyOnLvl int) {

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
			condenseUpwards(lpn, couple, "div")
		}

		//
		// next node
		lp = queue.DeQueue()
	}
}
