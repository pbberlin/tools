package fetch

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pbberlin/tools/dom"
	// "github.com/pbberlin/tools/util"
	"code.google.com/p/go.net/html"
)

var fCondenseNode func(*html.Node, int) string
var fRecurse func(*html.Node)

var nums int

const emptySrc = "//:0"

// needed to get the current request into the
// "static" recursive functions
// attn: unsynchronized
var UnsyncedGlobalReq *http.Request

func ModifyHTML(r *http.Request, s string) string {

	UnsyncedGlobalReq = r

	var docRoot *html.Node
	var err error
	r1 := strings.NewReader(s)
	log.Printf("len is %v\n", len(s))

	docRoot, err = html.Parse(r1)
	if err != nil {
		panic(fmt.Sprintf("3 %v \n", err))
	}
	fRecurse(docRoot)

	var b bytes.Buffer
	err = html.Render(&b, docRoot)
	if err != nil {
		panic(fmt.Sprintf("4 %v \n", err))
	}
	log.Printf("len is %v\n", b.Len())

	return b.String()
}

func init() {

	fCondenseNode = func(n *html.Node, depth int) (ret string) {

		if n.Type == html.ElementNode && n.Data == "script" {
			ret += fmt.Sprintf(" var script%v = '[script]'; ", nums)
			nums++
			return
		}
		if n.Type == html.ElementNode && n.Data == "style" {
			ret += fmt.Sprintf(" .xxx {margin:2px;} ")
			return
		}

		if n.Type == html.ElementNode && n.Data == "img" {
			ret += fmt.Sprintf(" [img] %v %v | ", getAttrVal(n.Attr, "alt"), getAttrVal(n.Attr, "src"))
		}

		if n.Type == html.ElementNode && n.Data == "a" {
			ret += "[a]"
		}

		if n.Type == html.TextNode {
			s := n.Data
			// s = replTabsNewline.Replace(s)
			// s = strings.TrimSpace(s)
			if len(s) < 4 {
				ret += s
			} else if s != "" {
				if depth > 0 {
					ret += fmt.Sprintf(" [txt%v] %v", depth, s)
				} else {
					ret += " [txt] " + s
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			ret += fCondenseNode(c, depth+1)
		}
		return
	}

	fRecurse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "form" {

			hidFld := new(html.Node)
			hidFld.Type = html.ElementNode
			hidFld.Data = "input"
			hidFld.Attr = []html.Attribute{html.Attribute{Key: "name", Val: "redirect-to"}, html.Attribute{Key: "value", Val: absolutize(getAttrVal(n.Attr, "action"))}}
			n.AppendChild(hidFld)

			n.Attr = rewriteAttributes(n.Attr, UnsyncedGlobalReq)

		}
		if n.Type == html.ElementNode && n.Data == "script" {
			for i := 0; i < len(n.Attr); i++ {
				if n.Attr[i].Key == "src" {
					n.Attr[i].Val = emptySrc
				}
			}
		}
		if n.Type == html.ElementNode &&
			(n.Data == "a" || n.Data == "img" || n.Data == "script" || n.Data == "style") {

			s := fCondenseNode(n, 0)
			//fmt.Printf("found %v\n", s)
			textReplacement := new(html.Node)
			textReplacement.Type = html.TextNode
			textReplacement.Data = s

			if n.Data == "a" || n.Data == "img" {
				n.Attr = rewriteAttributes(n.Attr, UnsyncedGlobalReq)
			}

			// We want to remove all existing children.
			// Direct loop impossible, since "NextSibling" is set to nil by Remove().
			// Therefore first assembling separately, then removing.
			children := make(map[*html.Node]struct{})
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				children[c] = struct{}{}
			}
			for k, _ := range children {
				n.RemoveChild(k)
			}

			// we can't put our replacement "under" an image, since img cannot have children
			if n.Type == html.ElementNode && n.Data == "img" {
				// n.Parent.InsertBefore(textReplacement,n)
				dom.InsertAfter(n, textReplacement)
				dom.RemoveNode(n)

			} else {
				n.AppendChild(textReplacement)
			}

			if n.Data == "a" {
				prev := n.PrevSibling
				if prev != nil {

					breaker0 := new(html.Node)
					breaker0.Type = html.TextNode
					breaker0.Data = " || "
					n.Parent.InsertBefore(breaker0, prev)

					breaker1 := new(html.Node)
					breaker1.Type = html.ElementNode
					// breaker1.Data =  "||<br>\n"
					breaker1.Data = "br"
					n.Parent.InsertBefore(breaker1, prev)

					breaker2 := new(html.Node)
					breaker2.Type = html.TextNode
					breaker2.Data = "\n"
					n.Parent.InsertBefore(breaker2, prev)

				}
			}

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fRecurse(c)
		}
	}

}
