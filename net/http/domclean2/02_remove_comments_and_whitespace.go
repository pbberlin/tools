package domclean2

import (
	"github.com/pbberlin/tools/net/http/dom"
	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

// NdX is a html.node, extended by its level.
// It's used since the horizontal traversal with
// a queue has no recursion and therefore
// keeps no depth information.
type NdX struct {
	Nd  *html.Node
	Lvl int
}

// removeCommentsAndIntertagWhitespace employs horizontal traversal using a queue
func removeCommentsAndIntertagWhitespace(lp interface{}) {

	var queue = util.NewQueue(10)

	for lp != nil {

		lpn := lp.(NdX).Nd
		lvl := lp.(NdX).Lvl

		for c := lpn.FirstChild; c != nil; c = c.NextSibling {
			queue.EnQueue(NdX{c, lvl + 1})
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
