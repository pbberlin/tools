// Package dom supplies simple node manipulations.
package dom

import "golang.org/x/net/html"

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
