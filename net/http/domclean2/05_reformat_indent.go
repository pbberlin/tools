package domclean2

import (
	"log"
	"strings"

	"github.com/pbberlin/tools/net/http/dom"
	"golang.org/x/net/html"
)

// !DOCTYPE html head
// !DOCTYPE html body
//        0    1    2
const cScaffoldLvls = 2

func reIndent(n *html.Node, lvl int) {

	if lvl > cScaffoldLvls && n.Parent == nil {
		bb := dom.PrintSubtree(n, nil, 0)
		_ = bb
		// log.Printf("%s", bb.Bytes())
		hint := ""
		if ml3[n] > 0 {
			hint = "   from ml3"
		}
		log.Print("reIndent: no parent ", hint)
		return
	}

	// Before children processing
	switch n.Type {
	case html.ElementNode:
		if lvl > cScaffoldLvls && n.Parent.Type == html.ElementNode {
			ind := strings.Repeat("\t", lvl-2)
			dom.InsertBefore(n, &html.Node{Type: html.TextNode, Data: "\n" + ind})
		}
	case html.CommentNode:
		dom.InsertBefore(n, &html.Node{Type: html.TextNode, Data: "\n"})
	case html.TextNode:
		n.Data = strings.TrimSpace(n.Data) + " "
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		reIndent(c, lvl+1)
	}

	// After children processing
	switch n.Type {
	case html.ElementNode:
		// I dont know why,
		// but this needs to happend AFTER the children
		if lvl > cScaffoldLvls && n.Parent.Type == html.ElementNode {
			ind := strings.Repeat("\t", lvl-2)
			if n.LastChild != nil {
				dom.InsertAfter(n.LastChild, &html.Node{Type: html.TextNode, Data: "\n" + ind})
			}
		}
	}

}
