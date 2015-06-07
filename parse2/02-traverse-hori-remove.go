package parse2

import (
	"github.com/pbberlin/tools/dom"
	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

func TravHoriRemoveCommentAndSpaces(lp interface{}) {

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

		if lpn.Type == html.TextNode && isSpacey(lpn.Data) {
			dom.RemoveNode(lpn)
		}

		// next node
		lp = queue.DeQueue()
	}
}
