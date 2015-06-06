package parse2

import (
	"fmt"

	"golang.org/x/net/html"
)

func printAttr(attributes []html.Attribute, keys []string) {
	for _, a := range attributes {
		for i := 0; i < len(keys); i++ {
			if keys[i] == a.Key {
				fmt.Printf("id is %v\n", a.Val)
			}
		}
	}
}

var idCntr = 0

func addIdAttr(attributes []html.Attribute) []html.Attribute {
	hasId := false
	for _, a := range attributes {
		if a.Key == "id" {
			hasId = true
			break
		}
	}
	if !hasId {
		attributes = append(attributes, html.Attribute{"", "id", fmt.Sprintf("d%v", idCntr)})
		idCntr++
	}
	return attributes
}

func printLvl(n *html.Node, col int) {
	if n.Type == html.ElementNode {
		fmt.Printf("%2v: %2v ", col, n.Data)
	}
}
