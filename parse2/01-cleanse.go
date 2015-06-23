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

func maxTreeDepth(n *html.Node, lvl int) (maxLvl int) {

	maxLvl = lvl
	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret := maxTreeDepth(c, lvl+1)
		if ret > maxLvl {
			maxLvl = ret
		}
	}
	return
}

func cleanseDom(n *html.Node, lvl int) {

	n.Attr = removeAttr(n.Attr, removeAttributes)

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cleanseDom(c, lvl+1)
	}

	if true {
		removeDirect(n)
	} else {
		convertToComment(n)
	}

	// ---

	normalizeToDiv(n)

	// one time text normalization
	if n.Type == html.TextNode {
		n.Data = textNormalize(n.Data)
	}

}

// convertToComment neutralizes a node.
// Note: We can not Remove() nor Replace()
// Since that breaks the recursion one step above!
// At a later stage we will employ horizontal traversal
// to actually remove unwanted nodes.
func convertToComment(n *html.Node) {
	if removes[n.Data] {
		n.Type = html.CommentNode
		n.Data = n.Data + " replaced"
	}
}

// We want to remove some children.
// A direct loop is impossible,
// since "NextSibling" is set to nil during Remove().
// Therefore:
//   First assemble children separately.
//   Then remove them.
func removeDirect(n *html.Node) {
	cc := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cc = append(cc, c)
	}
	for _, c := range cc {
		if removes[c.Data] {
			n.RemoveChild(c)
		}
	}
}

func normalizeToDiv(n *html.Node) {
	if repl, ok := simplifies[n.Data]; ok {
		n.Attr = append(n.Attr, html.Attribute{"", "cfrm", n.Data})
		n.Data = repl
	}
}
