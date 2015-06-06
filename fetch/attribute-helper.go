package fetch

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pbberlin/tools/util_appengine"

	"golang.org/x/net/html"
)

// type Attribute struct {
//     Namespace, Key, Val string
// }
func rewriteAttributes(attributes []html.Attribute, r *http.Request) []html.Attribute {

	rew := make([]html.Attribute, 0, len(attributes))

	for i := 0; i < len(attributes); i++ {
		attr := attributes[i]

		if attr.Key == "class" || attr.Key == "style" {
			continue
		}

		if attr.Key == "href" || attr.Key == "src" || attr.Key == "action" { //  make absolute
			attr.Val = absolutize(attr.Val)
		}

		if attr.Key == "href" || attr.Key == "src" {
			if attr.Key == "href" {
				attr.Val = fmt.Sprintf("/%v?url=%v", fetchURL, attr.Val)
			}
		}

		if attr.Key == "action" {
			if attr.Key == "href" || attr.Key == "action" {
				// attr.Val = fmt.Sprintf("/blob2/form-redirector?redirect-to=%v", attr.Val) // appended as form field, thus not needed here
				if util_appengine.IsLocalEnviron() {
					attr.Val = fmt.Sprintf("http://%v/blob2/form-redirector", r.Host)
				} else {
					attr.Val = fmt.Sprintf("https://%v/blob2/form-redirector", r.Host)
				}

			}
		}

		if attr.Key == "method" {
			attr.Val = "post"
		}

		rew = append(rew, attr)
	}

	rew = append(rew, html.Attribute{Key: "was", Val: "rewritten"})
	rew = append(rew, html.Attribute{Key: "method", Val: "post"})

	return rew
}

func getAttrVal(attributes []html.Attribute, key string) string {
	for i := 0; i < len(attributes); i++ {
		attr := attributes[i]
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func absolutize(val string) string {
	if strings.HasPrefix(val, "/") && !strings.HasPrefix(val, "//ssl.") {
		val = fmt.Sprintf("https://%v%v", host, val)
	}
	return val
}
