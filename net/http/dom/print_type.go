package dom

import "golang.org/x/net/html"

type NodeTypeStr html.NodeType

func (n NodeTypeStr) String() string {
	switch n {
	case 0:
		return "ErrorNode"
	case 1:
		return "TextNode"
	case 2:
		return "DocumentNode"
	case 3:
		return "ElementNode"
	case 4:
		return "CommentNode"
	case 5:
		return "DoctypeNode"
	}
	return "unknown Node type"
}
