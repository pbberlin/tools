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

	// normalize on title
	if alt != "" && alt == title {
		ret1 := []html.Attribute{}
		for i := 0; i < len(ret); i++ {
			if ret[i].Key != "alt" {
				ret1 = append(ret1, ret[i])
			}
		}
		ret = ret1
	}

	// remove both
	if alt == "" && alt == title {
		ret1 := []html.Attribute{}
		for i := 0; i < len(ret); i++ {
			if ret[i].Key != "alt" && ret[i].Key != "title" {
				ret1 = append(ret1, ret[i])
			}
		}
		ret = ret1
	}

	return ret
}

func addIdAttr(attributes []html.Attribute, id string) []html.Attribute {
	hasId := false
	for _, a := range attributes {
		if a.Key == "id" {
			hasId = true
			break
		}
	}
	if !hasId {
		attributes = append(attributes, html.Attribute{"", "id", id})
	}
	return attributes
}
