package domclean2

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func addOutlineAttr(n *html.Node, lvl int, argOutline []int) (outline []int) {

	outline = argOutline

	if n.Type == html.ElementNode && lvl > cScaffoldLvls {

		outline[len(outline)-1]++

		s := ""
		for _, v := range outline {
			s = fmt.Sprintf("%v%v.", s, v)
		}
		if strings.HasSuffix(s, ".") {
			s = s[:len(s)-1]
		}
		n.Attr = append(n.Attr, html.Attribute{"", "ol", s})

		outline = append(outline, 0) // add children lvl
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline = addOutlineAttr(c, lvl+1, outline)
	}

	if n.Type == html.ElementNode && lvl > cScaffoldLvls {
		outline = outline[:len(outline)-1] // reset children lvl
	}

	return
}

func addIdAttr(n *html.Node, lvl int, argNum int) (num int) {

	num = argNum

	if lvl > cScaffoldLvls {
		if n.Type == html.ElementNode {
			n.Attr = append(n.Attr, html.Attribute{"", "id", spf("%v", num)})
		}
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		num = addIdAttr(c, lvl+1, num+1)
	}

	return
}
