package domclean2

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/pbberlin/tools/net/http/dom"
	"golang.org/x/net/html"
)

func flattenSubtreeV2(n *html.Node, b *bytes.Buffer, depth int, tpar *html.Node) (*bytes.Buffer, *html.Node) {

	if b == nil {
		b = new(bytes.Buffer)
	}
	if tpar == nil {
		tpar = &html.Node{
			Type:     n.Type,
			DataAtom: n.DataAtom,
			Data:     n.Data,
			Attr:     make([]html.Attribute, len(n.Attr)),
		}
		copy(tpar.Attr, n.Attr)
	}

	switch {
	case n.Type == html.ElementNode && n.Data == "a":
		n.Parent.RemoveChild(n)
		tpar.AppendChild(n)
		// wpf(b, "[a] ")
	case n.Type == html.ElementNode && n.Data == "img":
		// img2Link(n)
		n.Parent.RemoveChild(n)
		tpar.AppendChild(n)
	case n.Data == "em" || n.Data == "strong":
		wpf(b, "[%v l%v] ", n.Data, depth)
		n.Parent.RemoveChild(n)
		tpar.AppendChild(n)
	case n.Data == "label" || n.Data == "input" || n.Data == "textarea":
		n.Parent.RemoveChild(n)
		tpar.AppendChild(n)
	case n.Data == "p" || n.Data == "div" || n.Data == "li" || n.Data == "ol" || n.Data == "h1" || n.Data == "h2" || n.Data == "ul":
		n.Parent.RemoveChild(n)
		tpar.AppendChild(n)
	case n.Data == "span":
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			n.RemoveChild(c)
			tpar.AppendChild(c)
		}
		n.Parent.RemoveChild(n)
	case n.Type == html.TextNode && n.Data != "":
		n.Data = strings.TrimSpace(n.Data)
		n.Data += " "
		wpf(b, n.Data)
		n.Parent.RemoveChild(n)
		tpar.AppendChild(n)
	default:
		log.Printf("unhandled %s %s\n", dom.NodeTypeStr(n.Type), n.Data)
		n.Parent.RemoveChild(n)
	}

	//
	//
	children := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// fmt.Printf("still has children %v\n", c.Data)
		children = append(children, c) //  assembling separately, before removing.
	}
	for _, c := range children {
		flattenSubtreeV2(c, b, depth+1, tpar)
	}

	return b, tpar
}

func condenseBottomUpV2(n *html.Node, lvl, lvlDo int, types map[string]bool) {

	if lvl < lvlDo {

		cs := []*html.Node{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cs = append(cs, c)
		}
		for _, c := range cs {
			condenseBottomUpV2(c, lvl+1, lvlDo, types)
		}

	} else {

		// log.Printf("action on %v %v\n", lvl, lvlDo)

		switch {

		case n.Type == html.ElementNode && types[n.Data]:

			oldPar := n.Parent
			if oldPar == nil {
				return
			}

			b, newPar := flattenSubtreeV2(n, nil, 0, nil)

			// placeholder := dom.Nd("div")
			// par := n.Parent
			// par.InsertBefore(placeholder, n.NextSibling)
			// par.RemoveChild(n)
			// par.InsertBefore(n2, placeholder)

			for c := oldPar.FirstChild; c != nil; c = c.NextSibling {
				oldPar.RemoveChild(c)
			}

			for c := newPar.FirstChild; c != nil; c = c.NextSibling {
				newPar.RemoveChild(c)
				oldPar.AppendChild(c)
			}

			if lvlDo > 4 {
				bx := dom.PrintSubtree(newPar)
				fmt.Printf("%s", bx)
			}

			// n = n2

			nodeRepl := dom.Nd("text", b.String())

			if false {

				// Remove all existing children.
				// Direct loop impossible, since "NextSibling" is set to nil by Remove().
				children := []*html.Node{}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					children = append(children, c) //  assembling separately, before removing.
				}
				for _, c := range children {
					log.Printf("c %4v rem from %4v ", c.Data, n.Data)
					n.RemoveChild(c)
				}

				// we can't put our replacement "under" an image, since img cannot have children
				if n.Type == html.ElementNode && n.Data == "img" {
					n.Parent.InsertBefore(nodeRepl, n.NextSibling) // if n.NextSibling==nil => insert at the end
					n.Parent.RemoveChild(n)
				} else {
					n.AppendChild(nodeRepl)
				}

				// Insert a  || and a newline before every <a...>
				// if n.Data == "a" {
				// 	n.Parent.InsertBefore(dom.Nd("text", " || "), n)
				// }
			}

		default:
		}

	}

}
