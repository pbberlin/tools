package domclean2

import (
	"log"

	"github.com/pbberlin/tools/net/http/dom"
	"golang.org/x/net/html"
)

func searchImg(n *html.Node, fnd *html.Node, lvl int) *html.Node {

	if n.Type == html.ElementNode && n.Data == "img" {
		// log.Printf("  a has img on lvl %v\n", lvl)
		if fnd == nil {
			fnd = n
			return fnd
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		fnd = searchImg(c, fnd, lvl+1)
		if fnd != nil {
			return fnd
		}
	}

	return fnd
}

type DeleterFunc func(*html.Node, int, bool) bool

func closureDeleter(until bool) DeleterFunc {

	// Nodes along the path to the splitting image
	// should never not be removed in *neither* tree
	var splitPath = map[*html.Node]bool{}

	var fc DeleterFunc
	fc = func(n *html.Node, lvl int, found bool) bool {

		// fmt.Printf("found %v at l%v\n", found, lvl)
		if n.Data == "img" {
			// fmt.Printf(" found at l%v\n", lvl)
			found = true
			par := n.Parent
			for {
				if par == nil {
					break
				}
				splitPath[par] = true
				par = par.Parent
			}
		}

		// children
		cc := []*html.Node{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cc = append(cc, c)
		}
		for _, c := range cc {
			found = fc(c, lvl+1, found)
		}

		//
		// remove
		if lvl > 0 {
			if n.Data == "img" {
				n.Parent.RemoveChild(n)
			} else {
				if !until && !found && !splitPath[n] {
					n.Parent.RemoveChild(n)
				}
				if until && found && !splitPath[n] {
					n.Parent.RemoveChild(n)
				}
			}
		}

		return found

	}

	return fc

}

func splitAnchSubtreesByImage(n *html.Node) {

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		splitAnchSubtreesByImage(c)
	}

	if n.Type == html.ElementNode && n.Data == "a" {
		img := searchImg(n, nil, 0)
		if img != nil {
			b0 := dom.PrintSubtree(n)
			log.Printf("\n%s\n", b0)

			// log.Printf("  got it  %v\n", img.Data)
			a1 := dom.CloneNodeWithSubtree(n)
			fc1 := closureDeleter(true)
			fc1(n, 0, false)

			b1 := dom.PrintSubtree(n)
			log.Printf("\n%s\n", b1)

			fc2 := closureDeleter(false)
			fc2(a1, 0, false)
			b2 := dom.PrintSubtree(a1)
			log.Printf("\n%s\n", b2)
			log.Printf("--------------------\n")
		} else {
			// log.Printf("no img in a\n")
		}
	}

}
