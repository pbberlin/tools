package domclean2

import (
	"fmt"
	"strings"

	"github.com/pbberlin/tools/stringspb"
	"golang.org/x/net/html"
)

func attrsX(attributes []html.Attribute, keys []string) (s string) {
	for _, a := range attributes {
		for i := 0; i < len(keys); i++ {
			if keys[i] == a.Key {
				s += fmt.Sprintf("%v is %v\n", a.Key, a.Val)
			}
		}
	}
	return
}

func attrX(attributes []html.Attribute, key string) (s string) {
	for _, a := range attributes {
		if key == a.Key {
			s = a.Val
			break
		}
	}
	return
}

func attrSet(attrs []html.Attribute, key, val string) []html.Attribute {
	for i, a := range attrs {
		if a.Key == key {
			attrs[i].Val = val
			return attrs
		}
	}
	// attr does not exist => append it
	attrs = append(attrs, html.Attribute{Key: key, Val: val})
	return attrs
}

func removeAttr(attributes []html.Attribute, removeKeys map[string]bool) []html.Attribute {

	ret := []html.Attribute{}
	var alt, title string

	for _, a := range attributes {
		a.Key = strings.TrimSpace(a.Key)
		a.Val = strings.TrimSpace(a.Val)
		a.Val = stringspb.NormalizeInnerWhitespace(a.Val) // having encountered title or alt values with newlines
		if removeKeys[a.Key] || strings.HasPrefix(a.Key, "data") {
			//
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

func Unused_addIdAttr(attributes []html.Attribute, id string) []html.Attribute {
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
