package domclean2

import (
	"github.com/pbberlin/tools/net/http/dom"
	"golang.org/x/net/html"
)

func flattenTraverse(n *html.Node) {

	checkForFlattening(n)

	cc := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cc = append(cc, c)
	}
	for _, c := range cc {
		flattenTraverse(c)
	}

}

func checkForFlattening(n *html.Node) {

	if n.Type == html.ElementNode && n.Data == "a" {

		nd := &html.Node{Type: html.TextNode}
		nd.Data = flattenNodeBelow(n, 0)
		nd.Data = textNormalize(nd.Data)

		cc := []*html.Node{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cc = append(cc, c)
		}
		for _, c := range cc {
			n.RemoveChild(c)
		}

		n.AppendChild(nd)

		nd2 := &html.Node{Type: html.ElementNode, Data: "br"}
		dom.InsertAfter(n, nd2)

	}

}

// one under starting node,
// one under lvl 0
func flattenNodeBelow(n *html.Node, lvl int) (ret string) {

	if lvl > 0 {
		if n.Type == html.ElementNode {
			ret += spf("[%v] ", n.Data)
			for _, v := range []string{"src", "alt", "title"} {
				av := attrX(n.Attr, v)
				if len(av) > 0 {
					ret += spf("%v ", av)
					// ret += spf("%v ", stringspb.Ellipsoider(av, 5))
				}
			}
		} else if n.Type == html.TextNode {
			ret += n.Data
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += flattenNodeBelow(c, lvl+1)
	}

	return
}
