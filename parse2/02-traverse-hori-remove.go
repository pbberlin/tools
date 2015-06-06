package parse2

import (
	"github.com/pbberlin/tools/dom"
	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

func TraverseHoriRemoveNodes(lp interface{}) {

	var queue = util.NewQueue(10)

	for lp != nil {

		lpn := lp.(Tx).Nd
		lvl := lp.(Tx).Lvl

		if lpn.Type == html.CommentNode {
			// fmt.Printf("comment removed\n")
			dom.RemoveNode(lpn)
		}

		// enqueue all children
		for c := lpn.FirstChild; c != nil; c = c.NextSibling {
			// if c.Type == html.ElementNode || c.Type == html.CommentNode {
			queue.EnQueue(Tx{c, lvl + 1})
			// }
		}
		lp = queue.DeQueue()
	}
}
