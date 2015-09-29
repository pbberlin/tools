package dedup

import (
	"github.com/pbberlin/tools/net/http/dom"
	"github.com/pbberlin/tools/stringspb"
	"golang.org/x/net/html"
)

func textifyBruteForce(n *html.Node) {

	cc := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cc = append(cc, c)
	}
	for _, c := range cc {
		textifyBruteForce(c)
	}

	textifyNodeSubtree(n)

}

func textifyNodeSubtree(n *html.Node) {

	if n.Type == html.ElementNode {

		nd := dom.Nd("text")
		nd.Data = textifySubtreeBruteForce(n, 0)
		nd.Data = stringspb.NormalizeInnerWhitespace(nd.Data)

		cc := []*html.Node{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cc = append(cc, c)
		}
		for _, c := range cc {
			n.RemoveChild(c)
		}

		n.AppendChild(nd)

		nd2 := dom.Nd("br")
		dom.InsertAfter(n, nd2)

	}

}
