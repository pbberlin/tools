package parse

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/pbberlin/tools/util"

	"code.google.com/p/go.net/html"
)

var fNodeModify func(*html.Node) string
var fRecurse func(*html.Node)
var replTabsNewline = strings.NewReplacer("\n", " ", "\t", "")

func ParseHtmlFiles() {

	testDataDir := "./"
	testFiles, err := filepath.Glob(testDataDir + "test*.html")
	if err != nil {
		pf("%v \n", err)
	}

	for _, tf := range testFiles {
		pf("%v\n", tf)

		f, err := os.Open(tf)
		if err != nil {
			pf("1 %v \n", err)
		}
		defer f.Close()
		r1 := bufio.NewReader(f)

		var docRoot *html.Node
		docRoot, err = html.Parse(r1)
		if err != nil {
			pf("3 %v \n", err)
		}

		fRecurse = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				s := strings.TrimSpace(fNodeModify(n))
				//pf("found %v\n", s)
				nNew := new(html.Node)
				nNew.Type = html.TextNode
				nNew.Data = s

				// We want to remove all children.
				// Direct loop impossible, since "NextSibling" is set to nil 
				// 		during Remove().
				// Therefore first assembling separately, then removing.
				children := map[*html.Node]string{}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					children[c] = "xx"
				}
				for k, _ := range children {
					n.RemoveChild(k)
					// pf("  removed  %q\n", strings.TrimSpace(k.Data))
				}
				n.AppendChild(nNew)

			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				fRecurse(c)
			}
		}
		fRecurse(docRoot)

		var b bytes.Buffer
		html.Render(&b, docRoot)
		util.WriteBytesToFilename("yy_"+tf, &b)
		//fixedHtml := b.String()

		//fmt.Printf("%s \n", spew.Sdump(docRoot))

	}
}

func init() {

	fNodeModify = func(n *html.Node) (ret string) {
		if n.Type == html.ElementNode && n.Data == "img" {
			ret += spf(" [img] %v ", getAttrVal(n.Attr, "alt"))
		}
		if n.Type == html.TextNode {
			s := replTabsNewline.Replace(n.Data)
			s = strings.TrimSpace(s)
			if len(s) < 4 {
				ret += s
			} else if s != "" {
				ret += " [txt] " + s
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			ret += fNodeModify(c)
		}
		return
	}

}
