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

		attr := html.Attribute{"", "ol", s}

		newAttrs := make([]html.Attribute, 0, len(n.Attr)+2) // make space for outline now - and id later
		newAttrs = append(newAttrs, attr)
		newAttrs = append(newAttrs, n.Attr...)
		n.Attr = newAttrs

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

			attr := html.Attribute{"", "id", spf("%v", num)}

			prep := []html.Attribute{attr}
			n.Attr = append(prep, n.Attr...)

		}
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		num = addIdAttr(c, lvl+1, num+1)
	}

	return
}
