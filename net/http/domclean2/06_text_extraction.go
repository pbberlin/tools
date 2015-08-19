package domclean2

import (
	"bytes"

	"github.com/pbberlin/tools/stringspb"
	"golang.org/x/net/html"
)

var textsByOutl = map[string][]byte{}
var textsByArticOutl = map[string][]SortEl{}

func textExtraction(n *html.Node, lvl int) (b []byte) {

	if lvl == 0 {
		textsByOutl = map[string][]byte{}
	}

	var cs []byte // content self
	var cc []byte // content children
	if n.Type == html.TextNode {
		cs = bytes.TrimSpace([]byte(n.Data))
		if len(cs) > 0 {
			cs = append(cs, byte(' '))
		}
	}
	if content, ok := inlineNodesToText(n); ok {
		cs = append(cs, content...)
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		var ccX []byte // content child X
		ccX = textExtraction(c, lvl+1)
		ccX = bytes.TrimSpace(ccX)
		if len(ccX) > 0 {
			ccX = append(ccX, byte(' '))
			cc = append(cc, ccX...)
		}
	}

	b = append(b, cs...)
	b = append(b, cc...)
	b = append(b, addHardBreaks(n)...)

	if lvl > cScaffoldLvls && (len(cs) > 0 || len(cc) > 0) && n.Type != html.TextNode {
		csCc := append(cs, cc...)
		ol := attrX(n.Attr, "ol")
		textsByOutl[ol] = csCc
	}
	return

}

// img and a nodes are converted into text nodes.
func inlineNodesToText(n *html.Node) (ct string, ok bool) {

	if n.Type == html.ElementNode {
		switch n.Data {

		case "br":
			ct, ok = "sbr ", true

		case "img":

			src := attrX(n.Attr, "src")
			src = stringspb.Ellipsoider(src, 5)

			alt := attrX(n.Attr, "alt")
			title := attrX(n.Attr, "title")

			if alt == "" && title == "" {
				ct = spf("[img] %v ", src)
			} else if alt == "" {
				ct = spf("[img] %v hbr %v ", src, title)
			} else {
				ct = spf("[img] %v hbr %v hbr %v ", src, title, alt)

			}
			ok = true

		case "a":

			href := attrX(n.Attr, "href")
			href = stringspb.Ellipsoider(href, 5)

			title := attrX(n.Attr, "title")
			if title == "" {
				ct = spf("[a] %v ", href)
			} else {
				ct = spf("[a] %v hbr %v ", href, title)
			}
			ok = true

		}

	}

	return

}

func addHardBreaks(n *html.Node) (s string) {

	if n.Type == html.ElementNode {
		switch n.Data {
		case "img":
			s = "hbr "
		case "p", "div":
			s = "hbr "
		}
	}
	return

}
