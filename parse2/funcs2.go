package parse2

import (
	"strings"

	"github.com/pbberlin/tools/dom"
	"golang.org/x/net/html"
)

var replTabsNewline = strings.NewReplacer("\r\n", " ", "\r", " ", "\n", " ")

func isSpacey(sarg string) bool {
	s := sarg
	s = replTabsNewline.Replace(s)
	s = strings.TrimSpace(s)
	if s == "" {
		// fmt.Printf("\t\t\tspacey: %q\n", sarg)
		return true
	}
	return false

}

func cleanseSpaceyTextNodes(n *html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if isSpacey(c.Data) {
			// fmt.Printf("spacey tn %q\n", c.Data)
			dom.RemoveNode(c)
		}
	}
}
