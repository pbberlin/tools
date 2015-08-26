package weedout

import (
	"bytes"
	"sort"
	"strings"

	"github.com/pbberlin/tools/stringspb"
	"golang.org/x/net/html"
)

// one under starting node,
// one under lvl 0
func textifySubtreeBruteForce(n *html.Node, lvl int) (ret string) {

	if lvl > 0 {
		if n.Type == html.ElementNode {
			ret += spf("[%v] ", n.Data)
			for _, v := range []string{"src", "alt", "title", "name", "type", "value"} {
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
		ret += textifySubtreeBruteForce(c, lvl+1)
	}

	return
}

// img and a nodes are converted into text nodes.
func inlineNodeToText(n *html.Node) (ct string, ok bool) {

	if n.Type == html.ElementNode {
		switch n.Data {

		case "br":
			ct, ok = "sbr ", true

		case "input":
			name := attrX(n.Attr, "name")
			stype := attrX(n.Attr, "type")
			val := attrX(n.Attr, "value")
			ct = spf("[inp] %v %v %v", name, stype, val)
			ok = true

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

var sortCompactReplace = map[rune]rune{
	'.': ' ',
	',': ',',
	'-': ' ',
	':': ' ',
	'/': ' ',
	'0': ' ',
	'1': ' ',
	'2': ' ',
	'3': ' ',
	'4': ' ',
	'5': ' ',
	'6': ' ',
	'7': ' ',
	'8': ' ',
	'9': ' ',
}

func sortCompact(text []byte) (buf []byte, histo map[string]int, numTokens int) {

	// text = bytes.Replace(text, []byte(" hbr"), []byte{}, -1)
	// text = bytes.Replace(text, []byte(" sbr"), []byte{}, -1)
	text = bytes.Replace(text, []byte(`[img] `), []byte{}, -1)
	// text = bytes.Replace(text, []byte(`[a] `), []byte{}, -1)

	mapping := func(r rune) rune {
		if ret, ok := sortCompactReplace[r]; ok {
			return ret
		}
		return r
	}

	text = bytes.Map(mapping, text)

	words := bytes.Fields(text)

	histo = map[string]int{}
	for _, word := range words {
		sword := string(word)
		sword = strings.TrimSpace(sword)
		sword = strings.ToLower(sword)
		if len(words) > 3 {
			if len(sword) > 3 {
				histo[sword]++
			}
		} else {
			histo[sword]++ // no minimum length for tiny texts
		}
	}
	numTokens = len(histo)

	keys := make([]string, 0, len(histo))
	for k, _ := range histo {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	buf = []byte{32}
	for _, key := range keys {
		if len(key) > 1 {
			buf = append(buf, []byte(key)...)
			buf = append(buf, byte(32))
			// num := fmt.Sprintf("%v", mp[key])
			// buf = append(buf, []byte(num)...)
			// buf = append(buf, byte(32))
		}
	}

	buf = bytes.TrimSpace(buf)

	return
}
