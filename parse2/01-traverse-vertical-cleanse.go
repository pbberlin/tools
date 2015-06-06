package parse2

import "golang.org/x/net/html"

var (
	removeTypes = map[string]bool{
		"meta":     true,
		"link":     true,
		"style":    true,
		"iframe":   true,
		"script":   true,
		"noscript": true,

		"canvas": true,
		"object": true,

		"wbr": true,
	}

	replaceNodeTypeBy = map[string]string{
		"section": "div",
		"article": "div",
		"header":  "div",
		"footer":  "div",
		"nav":     "div",
		"aside":   "div",

		"dl":     "div",
		"figure": "div",

		"dd":         "p",
		"dt":         "p",
		"figcaption": "p",
	}

	removeAttributes = map[string]bool{
		"style": true,
		"class": true,
		// "alt":                 true,
		// "title":               true,

		"align":       true,
		"placeholder": true,

		"target":   true,
		"id":       true,
		"rel":      true,
		"tabindex": true,
		"headline": true,

		"onload":      true,
		"onclick":     true,
		"onmousedown": true,
		"onerror":     true,

		"readonly":       true,
		"accept-charset": true,

		"itemprop":  true,
		"itemtype":  true,
		"itemscope": true,

		"datetime":               true,
		"current-time":           true,
		"fb-iframe-plugin-query": true,
		"fb-xfbml-state":         true,

		"frameborder":       true,
		"async":             true,
		"charset":           true,
		"http-equiv":        true,
		"allowtransparency": true,
		"allowfullscreen":   true,
		"scrolling":         true,
		"ftghandled":        true,
		"ftgrandomid":       true,
		"marginwidth":       true,
		"marginheight":      true,
		"vspace":            true,
		"hspace":            true,
		"seamless":          true,
		"aria-hidden":       true,
		"gapi_processed":    true,
		"property":          true,
		"media":             true,

		"content":  true,
		"language": true,

		"role": true,
	}

	nodeDistinct = map[string]int{}
	attrDistinct = map[string]int{}
)

func TraverseVertConvert(n *html.Node, lvl int) {

	PrepareUnwantedForDeletion(n)
	ConvertToDiv(n)

	switch n.Type {
	case html.ElementNode:
	}
	n.Attr = removeAttr(n.Attr, removeAttributes)

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseVertConvert(c, lvl+1)
	}

}

// PrepareUnwantedForDeletion neutralizes a node.
// Note: We can not Remove() nor Replace()
// Since that breaks the recursion one step above!
// At a later stage we will emplay horizontal traversal
// to actually remove unwanted nodes.
func PrepareUnwantedForDeletion(n *html.Node) {
	if removeTypes[n.Data] {
		// dom.ReplaceNode(n, &html.Node{Type: html.CommentNode, Data: n.Data + " replaced"})
		// dom.RemoveNode(n)
		n.Type = html.CommentNode
		// fmt.Printf("\tunwanted %9v turned into comment\n", n.Data)
		n.Data = n.Data + " replaced"
	}
}

func ConvertToDiv(n *html.Node) {
	if repl, ok := replaceNodeTypeBy[n.Data]; ok {
		n.Attr = append(n.Attr, html.Attribute{"", "converted-from", n.Data})
		n.Data = repl
	}
}
