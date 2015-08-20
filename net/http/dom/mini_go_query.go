// Package dom supplies simple node manipulations.
package dom

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pbberlin/tools/runtimepb"
	"golang.org/x/net/html"
)

var wpf = fmt.Fprintf

// inspired by https://github.com/PuerkitoBio/goquery/blob/master/manipulation.go

func ReplaceNode(self, dst *html.Node) {
	InsertAfter(self, dst)
	RemoveNode(self)
}

func RemoveNode(n *html.Node) {
	par := n.Parent
	if par != nil {
		par.RemoveChild(n)
	} else {
		log.Printf("\nNode to remove has no Parent\n")
		runtimepb.StackTrace(4)
	}
}

// InsertBefore inserts before itself.
// node.InsertBefore refers to its children
func InsertBefore(insPnt, toInsert *html.Node) {
	if insPnt.Parent != nil {
		insPnt.Parent.InsertBefore(toInsert, insPnt)
	} else {
		log.Printf("\nInsertBefore - insPnt has no Parent\n")
		runtimepb.StackTrace(4)
	}
}

// InsertBefore inserts at the end, when NextSibling is null.
// compare http://stackoverflow.com/questions/4793604/how-to-do-insert-after-in-javascript-without-using-a-library
func InsertAfter(insPnt, toInsert *html.Node) {
	if insPnt.Parent != nil {
		insPnt.Parent.InsertBefore(toInsert, insPnt.NextSibling)
	} else {
		log.Printf("\nInsertAfter - insPnt has no Parent\n")
		runtimepb.StackTrace(4)
	}
}

//
//
// Deep copy a node.
// The new node has clones of all the original node's
// children but none of its parents or siblingCNo
func cloneNodeWithSubtree(n *html.Node) *html.Node {
	nn := &html.Node{
		Type:     n.Type,
		DataAtom: n.DataAtom,
		Data:     n.Data,
		Attr:     make([]html.Attribute, len(n.Attr)),
	}

	copy(nn.Attr, n.Attr)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nn.AppendChild(cloneNodeWithSubtree(c)) // recursion
	}
	return nn
}

//
func PrintSubtree(n *html.Node, b *bytes.Buffer, lvl int) *bytes.Buffer {

	if b == nil {
		b = new(bytes.Buffer)
	}

	if lvl > 40 {
		log.Printf("%s", b.String())
		log.Printf("possible circular relationship\n")
		os.Exit(1)
	}

	ind := strings.Repeat(" ", lvl)
	wpf(b, "%sL%v %v", ind, lvl, NodeTypeStr(n.Type))
	wpf(b, " %v\n", n.Data)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		PrintSubtree(c, b, lvl+1) // recursion
	}

	return b
}

func Nd(ntype string, content ...string) *html.Node {

	nd0 := new(html.Node)

	if ntype == "text" {
		nd0.Type = html.TextNode
		if len(content) > 0 {
			nd0.Data = content[0]
		}
	} else {
		nd0.Type = html.ElementNode
		nd0.Data = ntype
		if len(content) > 0 {
			runtimepb.StackTrace(4)
			log.Printf("Element nodes can't have content")
		}
	}

	return nd0

}
