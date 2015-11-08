package domclean2

import (
	"fmt"
	"strings"

	"github.com/pbberlin/tools/net/http/dom"
	"github.com/pbberlin/tools/stringspb"
	"golang.org/x/net/html"
)

func closureTextNodeExists(img *html.Node) (found bool) {

	txt := attrX(img.Attr, "title")
	if len(txt) < 5 {
		return false
	}
	txt = stringspb.NormalizeInnerWhitespace(txt)
	txt = strings.TrimSpace(txt)

	// We dont search entire document, but three levels above image subtree
	grandParent := img
	for i := 0; i < 4; i++ {
		if grandParent.Parent != nil {
			grandParent = grandParent.Parent
		} else {
			// log.Printf("LevelsUp %v for %q", i, txt)
			break
		}
	}

	var recurseTextNodes func(n *html.Node)
	recurseTextNodes = func(n *html.Node) {

		if found {
			return
		}

		cc := []*html.Node{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cc = append(cc, c)
		}
		for _, c := range cc {
			recurseTextNodes(c)
		}

		if n.Type == html.TextNode {
			n.Data = stringspb.NormalizeInnerWhitespace(n.Data)
			if len(n.Data) >= len(txt) {
				// if strings.Contains(txt, "FDP") {
				// 	log.Printf("%25v     %v", stringspb.Ellipsoider(txt, 10), stringspb.Ellipsoider(n.Data, 10))
				// }
				fnd := strings.Contains(n.Data, txt)
				if fnd {
					found = true
					return
				}
			}
		}
	}
	recurseTextNodes(grandParent)

	return
}

func img2Link(img *html.Node) {

	if img.Data == "img" {

		img.Data = "a"
		for i := 0; i < len(img.Attr); i++ {
			if img.Attr[i].Key == "src" {
				img.Attr[i].Key = "href"
			}
		}

		double := closureTextNodeExists(img)
		imgContent := ""
		title := attrX(img.Attr, "title")

		if double {
			imgContent = fmt.Sprintf("[img] %v %v | ",
				"[ctdr]", // content title double removed
				urlBeautify(attrX(img.Attr, "href")))

		} else {
			imgContent = fmt.Sprintf("[img] %v %v | ",
				title,
				urlBeautify(attrX(img.Attr, "href")))
		}

		img.Attr = attrSet(img.Attr, "cfrom", "img")
		nd := dom.Nd("text", imgContent)
		img.AppendChild(nd)
	}

}

func recurseImg2Link(n *html.Node) {

	cc := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cc = append(cc, c)
	}
	for _, c := range cc {
		recurseImg2Link(c)
	}

	if n.Type == html.ElementNode && n.Data == "img" {
		img2Link(n)
	}
}
