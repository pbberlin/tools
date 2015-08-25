package weedout

import (
	"bytes"

	"golang.org/x/net/html"
)

func BubbledUpTextExtraction(n *html.Node, fnKey string) ([]*TextifiedTree, []byte) {

	// reset
	mp := []*TextifiedTree{}

	bts, mp := textExtract(n, fnKey, 0, mp)

	return mp, bts
}

func textExtract(n *html.Node, fnKey string, lvl int, mp []*TextifiedTree) ([]byte, []*TextifiedTree) {

	var cs []byte // content self
	var cc []byte // content children

	if n.Type == html.TextNode {
		cs = bytes.TrimSpace([]byte(n.Data))
		if len(cs) > 0 {
			cs = append(cs, byte(' '))
		}
	} else if n.Type == html.ElementNode {

		for _, v := range []string{"alt", "title"} {
			val := attrX(n.Attr, v)
			if len(val) > 0 {
				cs = append(cs, val...)
				cs = append(cs, byte(32))
			}
		}

	}
	// if content, ok := inlineNodeToText(n); ok {
	// 	cs = append(cs, content...)
	// }

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		var cChX []byte // content child X
		cChX, mp = textExtract(c, fnKey, lvl+1, mp)
		if len(cChX) > 0 {
			cChX = append(cChX, byte(' '))
			cc = append(cc, cChX...)
		}
	}

	if lvl > cScaffoldLvls && (len(cs) > 0 || len(cc) > 0) && n.Type != html.TextNode {
		csCc := append(cs, cc...)
		ol := attrX(n.Attr, "ol")
		compacted, numTokens := sortCompact(csCc)
		tt := &TextifiedTree{}
		tt.SourceID = fnKey
		tt.Lvl = lvl - cScaffoldLvls
		tt.Outline = ol
		tt.NumTokens = numTokens
		tt.Text = compacted
		mp = append(mp, tt)
	}

	b := new(bytes.Buffer)
	b.Write(cs)
	b.Write(cc)
	// b.WriteString(addHardBreaks(n))

	return b.Bytes(), mp

}
