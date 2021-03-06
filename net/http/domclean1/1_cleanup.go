// Package domclean1 normalizes html dom trees in a primitive way.
package domclean1

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pbberlin/tools/net/http/dom"
	"github.com/pbberlin/tools/net/http/fetch"
	"golang.org/x/net/html"
)

var fCondenseNode func(*html.Node, int) string
var fRecurse func(*html.Node)

const emptySrc = "//:0"

// r is the request to the proxy
// u is the url, that the proxy has called
func ModifyHTML(r *http.Request, u *url.URL, s string) string {

	var nums int // counter

	// needed to get the current request into the
	// "static" recursive functions
	var PackageProxyHost = r.Host // port included!
	var PackageRemoteHost = fetch.HostFromUrl(u)

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

	// --------------------------
	// ----------------------

	fRecurse = func(n *html.Node) {

		if n.Type == html.ElementNode && n.Data == "form" {
			hidFld := new(html.Node)
			hidFld.Type = html.ElementNode
			hidFld.Data = "input"
			hidFld.Attr = []html.Attribute{
				html.Attribute{Key: "name", Val: "redirect-to"},
				html.Attribute{Key: "value", Val: absolutize(getAttrVal(n.Attr, "action"), PackageRemoteHost)},
			}
			n.AppendChild(hidFld)

			submt := new(html.Node)
			submt.Type = html.ElementNode
			submt.Data = "input"
			submt.Attr = []html.Attribute{
				html.Attribute{Key: "type", Val: "submit"},
				html.Attribute{Key: "value", Val: "subm"},
				html.Attribute{Key: "accesskey", Val: "f"},
			}
			n.AppendChild(submt)

			n.Attr = rewriteAttributes(n.Attr, PackageProxyHost, PackageRemoteHost)

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

			attrStore := []html.Attribute{}
			if n.Data == "a" || n.Data == "img" {
				attrStore = rewriteAttributes(n.Attr, PackageProxyHost, PackageRemoteHost)
			}
			if n.Data == "img" {
				n.Data = "a"
			}
			if n.Data == "a" {
				n.Attr = attrStore
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

			// Insert a  || and a newline before every <a...>
			if n.Data == "a" {
				prev := n

				breaker0 := dom.Nd("text", "||")
				n.Parent.InsertBefore(breaker0, prev)

				breaker1 := dom.Nd("br")
				n.Parent.InsertBefore(breaker1, prev)

				breaker2 := dom.Nd("text", "\n")
				n.Parent.InsertBefore(breaker2, prev)
			}

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fRecurse(c)
		}
	}

	// --------------------------
	// ----------------------
	var docRoot *html.Node
	var err error
	rdr := strings.NewReader(s)
	docRoot, err = html.Parse(rdr)
	if err != nil {
		panic(fmt.Sprintf("3 %v \n", err))
	}

	fRecurse(docRoot)

	var b bytes.Buffer
	err = html.Render(&b, docRoot)
	if err != nil {
		panic(fmt.Sprintf("4 %v \n", err))
	}
	// log.Printf("len is %v\n", b.Len())

	return b.String()
}

func init() {

}
