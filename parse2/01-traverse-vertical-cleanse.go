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

func TravVertConvertEmptyLeafs(n *html.Node, lvl int) {

	// processing
	if n.FirstChild == nil &&
		n.Type == html.ElementNode &&
		(n.Data == "div" || n.Data == "span" || n.Data == "li") {
		n.Type = html.CommentNode
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TravVertConvertEmptyLeafs(c, lvl+1)
	}

}

func TraverseVert_ConvertDivDiv(n *html.Node, lvl int) {

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseVert_ConvertDivDiv(c, lvl+1)
	}

	couple := []string{"div", "div"}
	p := n.Parent
	if p == nil {
		return
	}

	cond1 := n.Type == html.ElementNode && n.Data == couple[0]
	cond2 := p.Type == html.ElementNode && p.Data == couple[1]

	noSiblings := n.PrevSibling == nil && n.NextSibling == nil
	onlyChild := n.FirstChild != nil && n.FirstChild == n.LastChild
	svrlChild := n.FirstChild != nil && n.FirstChild != n.LastChild

	if cond1 && cond2 {

		// fmt.Printf("%2v: %v/div/%v\n", lvl, p.Data, n.FirstChild.Data)

		if svrlChild || onlyChild {
			var children []*html.Node
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				children = append(children, c)
			}

			for _, c1 := range children {
				n.RemoveChild(c1)
				p.InsertBefore(c1, n.NextSibling)
			}
		}
		p.RemoveChild(n)

	}

	if false {

		if cond1 && cond2 && noSiblings {
			if onlyChild {
				fc := n.FirstChild
				p.FirstChild = fc
				// dom.RemoveNode(n)
				// fmt.Printf("<- ")
			}

		}
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
		n.Attr = append(n.Attr, html.Attribute{"", "converted-from", n.Data})
		n.Data = repl
	}
}
