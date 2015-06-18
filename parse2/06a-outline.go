package parse2

import (
	"strings"

	"golang.org/x/net/html"
)

var outlFmt = strings.NewReplacer("[", "", "]", "", " ", ".")

func computeOutline(n *html.Node, lvl int, argOutline []int) (outline []int) {

	outline = argOutline

	if n.Type == html.ElementNode && lvl > cScaffoldLvls {

		outline[len(outline)-1]++
		s := spf("%v", outline)
		s = outlFmt.Replace(s)
		n.Attr = append(n.Attr, html.Attribute{"", "ol", s})

		outline = append(outline, 0) // add children lvl
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline = computeOutline(c, lvl+1, outline)
	}

	if n.Type == html.ElementNode && lvl > cScaffoldLvls {
		outline = outline[:len(outline)-1] // reset children lvl
	}

	return
}

func nodeCountHoriz(n *html.Node, lvl int, argNum int) (num int) {

	num = argNum

	if lvl > cScaffoldLvls {
		if n.Type == html.ElementNode {
			n.Attr = append(n.Attr, html.Attribute{"", "id", spf("%v", num)})
		}
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		num = nodeCountHoriz(c, lvl+1, num+1)
	}

	return
}
