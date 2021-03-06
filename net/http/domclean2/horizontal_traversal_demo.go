package domclean2

import (
	"fmt"

	"github.com/pbberlin/tools/util"
	"golang.org/x/net/html"
)

// TraverseHori traverses the tree horizontally.
// It uses a queue. A FiFo structure.
// Inspired by www.geeksforgeeks.org/level-order-tree-traversal/
func traverseHori(lp interface{}) {

	var queue = util.NewQueue(10)

	lvlPrev := 0
	for lp != nil {

		lpn := lp.(NdX).Nd
		lvl := lp.(NdX).Lvl

		// print current
		if lvl != lvlPrev { // new level => newline
			fmt.Printf("\n%2v:\t", lvl)
			lvlPrev = lvl
		}
		fmt.Printf("%8s  ", lpn.Data)

		// enqueue all children
		for c := lpn.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				queue.EnQueue(NdX{c, lvl + 1})
			}
		}
		lp = queue.DeQueue()
	}
}
