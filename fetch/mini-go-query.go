package fetch

import (
	"code.google.com/p/go.net/html"
)


// inspired by https://github.com/PuerkitoBio/goquery/blob/master/manipulation.go


func ReplaceNode(src, dst *html.Node) {
	InsertAfter(src, dst)
	RemoveNode(src)
}

func RemoveNode(n *html.Node) {
	if n.Parent != nil {
		n.Parent.RemoveChild(n)
	}
}

// InsertBefore inserts at the end, when NextSibling is null.
// compare http://stackoverflow.com/questions/4793604/how-to-do-insert-after-in-javascript-without-using-a-library
func InsertAfter(src, dst *html.Node) {
	if src.Parent != nil {
		src.Parent.InsertBefore(dst, src.NextSibling)
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
