package parse2

import (
	"bytes"
	"fmt"

	"github.com/pbberlin/tools/pbstrings"
	"golang.org/x/net/html"
)

const cMinLen = 10

var mp = map[string][]byte{}

func textExtraction(n *html.Node, lvl, argHoriNum int) (b []byte, horiNum int) {

	if lvl == 0 {
		mp = map[string][]byte{}
	}

	horiNum = argHoriNum

	id := fmt.Sprintf("%v-%2v", lvl, horiNum)
	if lvl > cScaffoldLvls {
		if n.Type == html.ElementNode {
			n.Attr = addIdAttr(n.Attr, id)
		}
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
		// cs = append([]byte(content), cs...)
		cs = append(cs, content...)
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		var ccX []byte // content child X
		ccX, horiNum = textExtraction(c, lvl+1, horiNum+1)
		ccX = bytes.TrimSpace(ccX)
		if len(ccX) > 0 {
			ccX = append(ccX, byte(' '))
			cc = append(cc, ccX...)
		}
	}

	// b = append(b, "slf:"...)
	b = append(b, cs...)
	// b = append(b, "chn:"...)
	b = append(b, cc...)
	b = append(b, addHardBreaks(n)...)

	if lvl > cScaffoldLvls && (len(cs) > 0 || len(cc) > 0) && n.Type != html.TextNode {
		csCc := append(cs, cc...)
		idMap := fmt.Sprintf("%v-% 5v", id, len(csCc))
		if len(csCc) > cMinLen {
			mp[idMap] = csCc
		}
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

			href := attrX(n.Attr, "href")
			href = pbstrings.Ellipsoider(href, 5)

			alt := attrX(n.Attr, "alt")
			title := attrX(n.Attr, "title")

			if alt == "" && title == "" {
				ct = spf("[img] %v ", href)
			} else if alt == "" {
				ct = spf("[img] %v hbr %v ", title, href)
			} else {
				ct = spf("[img] %v hbr %v hbr %v ", title, alt, href)

			}

			ok = true
		case "a":
			href := attrX(n.Attr, "href")
			href = pbstrings.Ellipsoider(href, 5)

			title := attrX(n.Attr, "title")

			if title == "" {
				ct = spf("[a] %v ", href)
			} else {
				ct = spf("[a] %v hbr %v ", title, href)
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
