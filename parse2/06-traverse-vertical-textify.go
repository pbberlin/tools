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

	var self []byte // conent self
	var b1 []byte   // content children
	if n.Type == html.TextNode {
		self = bytes.TrimSpace([]byte(n.Data))
		if len(self) > 0 {
			self = append(self, byte(' '))
		}
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {

		b1, horiNum = TravVertTextify(c, lvl+1, horiNum+1)

		b1 = bytes.TrimSpace(b1)
		if len(b1) > 0 {
			b1 = append(b1, byte(' '))
		}
	}

	b = append(b, self...)
	b = append(b, b1...)

	if lvl > cScaffoldLvls {
		idMap := fmt.Sprintf("%v-%4v", id, len(b1))
		mp[idMap] = append(self, b1...)
	}

	return
}
