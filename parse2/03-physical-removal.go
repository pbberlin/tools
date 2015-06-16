package parse2

import (
	"github.com/pbberlin/tools/dom"
	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

// Attn: Horizontal traversal using a queue
func physicalNodeRemoval(lp interface{}) {

	var queue = util.NewQueue(10)

	for lp != nil {

		lpn := lp.(Tx).Nd
		lvl := lp.(Tx).Lvl

		for c := lpn.FirstChild; c != nil; c = c.NextSibling {
			queue.EnQueue(Tx{c, lvl + 1})
		}

		// processing
		if lpn.Type == html.CommentNode {
			dom.RemoveNode(lpn)
		}

		// extinguish textnodes that do only formatting (spaces, tabs, line breaks)
		if lpn.Type == html.TextNode && isSpacey(lpn.Data) {
			dom.RemoveNode(lpn)
		}

		// next node
		lp = queue.DeQueue()
	}
}
