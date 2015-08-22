package domclean2

import (
	"fmt"
	"log"

	"github.com/pbberlin/tools/net/http/dom"
	"golang.org/x/net/html"
)

var debugBreakOut = false

func searchImg(n *html.Node, fnd *html.Node, lvl int) (*html.Node, int) {

	if n.Type == html.ElementNode && n.Data == "img" {
		// log.Printf("  a has img on lvl %v\n", lvl)
		if fnd == nil {
			fnd = n
			return fnd, lvl
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		fnd, lvlfnd := searchImg(c, fnd, lvl+1)
		if fnd != nil {
			return fnd, lvlfnd
		}
	}

	return fnd, lvl
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

func breakoutImagesFromAnchorTrees(n *html.Node) {

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		breakoutImagesFromAnchorTrees(c)
	}

	if n.Type == html.ElementNode && n.Data == "a" {

		img, lvl := searchImg(n, nil, 0)

		if img != nil {

			only1Child := n.FirstChild != nil && n.FirstChild == n.LastChild
			if lvl == 1 && only1Child {
				// log.Printf("only child image lvl %v a\n", lvl)
				n.RemoveChild(img)
				n.Parent.InsertBefore(img, n.NextSibling) // "insert after; if n.NextSibling==nil => insert at the end"
				n.AppendChild(dom.Nd("text", "[was img]"))
			} else {

				if debugBreakOut {
					b0 := dom.PrintSubtree(n)
					log.Printf("\n%s\n", b0)
				}

				// log.Printf("  got it  %v\n", img.Data)
				a1 := dom.CloneNodeWithSubtree(n)
				fc1 := closureDeleter(true)
				fc1(n, 0, false)

				if debugBreakOut {
					b1 := dom.PrintSubtree(n)
					log.Printf("\n%s\n", b1)
				}

				fc2 := closureDeleter(false)
				fc2(a1, 0, false)
				if debugBreakOut {
					b2 := dom.PrintSubtree(a1)
					log.Printf("\n%s\n", b2)
					log.Printf("--------------------\n")
				}

				if true {
					n.Parent.InsertBefore(img, n.NextSibling) // "insert after; if n.NextSibling==nil => insert at the end"
					n.Parent.InsertBefore(a1, img.NextSibling)
				} else {
					// old way ; sequence corrpution if n had rightwise siblings.
					n.Parent.AppendChild(img)
					n.Parent.AppendChild(a1)

				}

			}

			// changing image to link:
			if img.Data == "img" {
				img.Data = "a"
				for i := 0; i < len(img.Attr); i++ {
					if img.Attr[i].Key == "src" {
						img.Attr[i].Key = "href"
					}
				}
				imgContent := fmt.Sprintf("[img] %v %v | ", attrX(img.Attr, "title"), attrX(img.Attr, "href"))
				img.Attr = attrSet(img.Attr, "cfrom", "img")
				nd := dom.Nd("text", imgContent)
				img.AppendChild(nd)
			}

		} else {
			// log.Printf("no img in a\n")
		}
	}

}
