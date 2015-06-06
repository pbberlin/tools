package parse2

import (
	"fmt"
	"strings"

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

func removeAttr(attributes []html.Attribute, removeKeys map[string]bool) []html.Attribute {
	ret := []html.Attribute{}
	var alt, title string
	for _, a := range attributes {
		if removeKeys[strings.TrimSpace(a.Key)] ||
			strings.HasPrefix(a.Key, "data") {
		} else {
			if a.Key == "alt" {
				alt = a.Val
			}
			if a.Key == "title" {
				title = a.Val
			}
			attrDistinct[a.Key]++
			ret = append(ret, a)
		}
	}
	if alt != "" && alt == title {
		for i := 0; i < len(ret); i++ {
			if ret[i].Key == "alt" {
				ret[i].Key = ""
				ret[i].Val = ""
				// fmt.Printf(" double alt,title %v\n", title)
				break
			}
		}

	}

	return ret
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
