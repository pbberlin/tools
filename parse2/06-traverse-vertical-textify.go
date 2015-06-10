package parse2

import (
	"bytes"
	"fmt"

	"golang.org/x/net/html"
)

var mp = map[string][]byte{}

func TravVertTextify(n *html.Node, lvl, argHoriNum int) (b []byte, horiNum int) {

	if lvl == 0 {
		mp = map[string][]byte{}
	}

	horiNum = argHoriNum

	id := fmt.Sprintf("%v-%v", lvl, horiNum)
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

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		var ccX []byte // content child X
		ccX, horiNum = TravVertTextify(c, lvl+1, horiNum+1)
		ccX = bytes.TrimSpace(ccX)
		if len(ccX) > 0 {
			ccX = append(ccX, byte(' '))
			cc = append(cc, ccX...)
		}
	}

	b = append(b, cs...)
	b = append(b, cc...)

	if lvl > cScaffoldLvls && (len(cs) > 0 || len(cc) > 0) && n.Type != html.TextNode {
		csCc := append(cs, cc...)
		idMap := fmt.Sprintf("%v-% 5v", id, len(csCc))
		mp[idMap] = csCc
	}

	return
}
