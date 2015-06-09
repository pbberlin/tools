package parse2

import "golang.org/x/net/html"

var (
	removes = map[string]bool{
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

	simplifies = map[string]string{
		"header":  "div",
		"footer":  "div",
		"nav":     "div",
		"section": "div",
		"article": "div",
		"aside":   "div",

		"dl": "ul",
		"dt": "li",
		"dd": "p",

		"figure":     "div",
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

	n.Attr = removeAttr(n.Attr, removeAttributes)

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseVertConvert(c, lvl+1)
	}

	NormalizeToDiv(n)
	ConvertToComment(n)

	// one time text normalization
	if n.Type == html.TextNode {
		n.Data = replNewLines.Replace(n.Data)
		n.Data = replTabs.Replace(n.Data)
		n.Data = doubleSpaces.ReplaceAllString(n.Data, " ")
	}

}

func TravVertMaxLevel(n *html.Node, lvl int) (maxLvl int) {

	maxLvl = lvl
	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret := TravVertMaxLevel(c, lvl+1)
		if ret > maxLvl {
			maxLvl = ret
		}
	}
	return
}

func TravVertConvertEmptyLeafs(n *html.Node, lvl int) {

	// processing
	if n.FirstChild == nil &&
		n.Type == html.ElementNode &&
		(n.Data == "div" || n.Data == "span" ||
			n.Data == "li" || n.Data == "p") {
		n.Type = html.CommentNode
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TravVertConvertEmptyLeafs(c, lvl+1)
	}

}

// ConvertToComment neutralizes a node.
// Note: We can not Remove() nor Replace()
// Since that breaks the recursion one step above!
// At a later stage we will emplay horizontal traversal
// to actually remove unwanted nodes.
func ConvertToComment(n *html.Node) {
	if removes[n.Data] {
		// dom.ReplaceNode(n, &html.Node{Type: html.CommentNode, Data: n.Data + " replaced"})
		// dom.RemoveNode(n)
		n.Type = html.CommentNode
		// fmt.Printf("\tunwanted %9v turned into comment\n", n.Data)
		n.Data = n.Data + " replaced"
	}
}

func NormalizeToDiv(n *html.Node) {
	if repl, ok := simplifies[n.Data]; ok {
		n.Attr = append(n.Attr, html.Attribute{"", "cfrm", n.Data})
		n.Data = repl
	}
}
