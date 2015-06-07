package parse2

import (
	"github.com/pbberlin/tools/dom"
	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

func TravHoriRemoveWhitesp(lp interface{}) {

	var queue = util.NewQueue(10)

	for lp != nil {

		lpn := lp.(Tx).Nd
		lvl := lp.(Tx).Lvl

		// enqueue all children
		for c := lpn.FirstChild; c != nil; c = c.NextSibling {
			// if c.Type == html.ElementNode || c.Type == html.CommentNode {
			queue.EnQueue(Tx{c, lvl + 1})
			// }
		}

		// processing
		if lpn.Type == html.TextNode {
			if isSpacey(lpn.Data) {
				dom.RemoveNode(lpn)
			}
		}

		// next node
		lp = queue.DeQueue()
	}
}

func TravHoriRemoveEmptyDiv(lp interface{}) {

	var queue = util.NewQueue(10)

	for lp != nil {

		lpn := lp.(Tx).Nd
		lvl := lp.(Tx).Lvl

		// processing
		if lpn.FirstChild == nil &&
			// lpn.PrevSibling == nil &&
			// lpn.NextSibling == nil &&
			lpn.Type == html.ElementNode &&
			lpn.Data == "div" {

			lpn.Type = html.CommentNode
			// dom.RemoveNode(lpn)

		}

		// enqueue all children
		for c := lpn.FirstChild; c != nil; c = c.NextSibling {
			// if c.Type == html.ElementNode || c.Type == html.CommentNode {
			queue.EnQueue(Tx{c, lvl + 1})
			// }
		}

		// next node
		lp = queue.DeQueue()
	}
}

func TravHoriRemoveComment(lp interface{}) {

	var queue = util.NewQueue(10)

	for lp != nil {

		lpn := lp.(Tx).Nd
		lvl := lp.(Tx).Lvl

		// enqueue all children
		for c := lpn.FirstChild; c != nil; c = c.NextSibling {
			// if c.Type == html.ElementNode || c.Type == html.CommentNode {
			queue.EnQueue(Tx{c, lvl + 1})
			// }
		}

		// processing
		if lpn.Type == html.CommentNode {
			dom.RemoveNode(lpn)
		}

		// next node
		lp = queue.DeQueue()
	}
}
