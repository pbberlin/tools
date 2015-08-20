// Package dom supplies simple node manipulations.
package dom

import (
	"bytes"
	"fmt"
	"strings"

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
	}
}

// InsertBefore inserts before itself.
// node.InsertBefore refers to its children
func InsertBefore(self, dst *html.Node) {
	if self.Parent != nil {
		self.Parent.InsertBefore(dst, self)
	}
}

// InsertBefore inserts at the end, when NextSibling is null.
// compare http://stackoverflow.com/questions/4793604/how-to-do-insert-after-in-javascript-without-using-a-library
func InsertAfter(self, dst *html.Node) {
	if self.Parent != nil {
		self.Parent.InsertBefore(dst, self.NextSibling)
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

	ind := strings.Repeat(" ", lvl)
	wpf(b, "%sL%v", ind, lvl)
	wpf(b, "T%v", n.Type)
	wpf(b, " D%v\n", n.Data)
	lvl++

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		PrintSubtree(c, b, lvl) // recursion
	}

	return b
}
