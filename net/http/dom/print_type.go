package dom

import "golang.org/x/net/html"

type NodeTypeStr html.NodeType

func (n NodeTypeStr) String() string {
	switch n {
	case 0:
		return "ErroNd"
	case 1:
		return "Text  "
	case 2:
		return "DocmNd"
	case 3:
		return "Elem  "
	case 4:
		return "CommNd"
	case 5:
		return "DoctNd"
	}
	return "unknown Node type"
}
