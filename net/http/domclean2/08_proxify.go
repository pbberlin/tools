package domclean2

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/net/http/dom"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/routes"
	"golang.org/x/net/html"
)

type FuncType2 func(*html.Node)

const emptySrc = "//:0"

// r is the request to the proxy
// u is the url, that the proxy has called
func closuredProxifier(argProxyHostPort string, urlSrc *url.URL) FuncType2 {

	// needed to get the current request into the
	// "static" recursive functions
	var closProxyHostPort = argProxyHostPort // port included!

	var closRemoteHost = fetch.HostFromUrl(urlSrc)
	// log.Printf("ProxyHost %v, RemoteHost %v (%s)", closProxyHostPort, closRemoteHost, urlSrc)

	// --------------------------
	// ----------------------

	var fRecurse FuncType2
	fRecurse = func(n *html.Node) {

		switch {
		case n.Type == html.ElementNode && n.Data == "form":
			hidFld := dom.Nd("input")
			hidFld.Attr = []html.Attribute{
				html.Attribute{Key: "name", Val: "redirect-to"},
				html.Attribute{Key: "value", Val: attrX(n.Attr, "action")},
			}
			n.AppendChild(hidFld)

			submt := dom.Nd("input")
			submt.Attr = []html.Attribute{
				html.Attribute{Key: "type", Val: "submit"},
				html.Attribute{Key: "value", Val: "subm"},
				html.Attribute{Key: "accesskey", Val: "f"},
			}
			n.AppendChild(submt)

			n.Attr = attrSet(n.Attr, "method", "post")
			n.Attr = attrSet(n.Attr, "was", "rewritten")

			n.Attr = attrsAbsoluteAndProxified(n.Attr, closProxyHostPort, closRemoteHost)

		case n.Type == html.ElementNode && (n.Data == "a" || n.Data == "img"):

			if n.Data == "a" || n.Data == "img" {
				attrStore := attrsAbsoluteAndProxified(n.Attr, closProxyHostPort, closRemoteHost)
				n.Attr = attrStore
			}

		default:
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fRecurse(c)
		}
	}

	return fRecurse

}

func absolutize(scheme, host, val string) (ret string) {
	if strings.HasPrefix(val, "/") && !strings.HasPrefix(val, "//ssl.") {
		ret = fmt.Sprintf("%v://%v%v", scheme, host, val)
		// log.Printf("absolutized %v %v - %v", val, host, ret)
	} else {
		ret = val
	}
	return
}

func attrsAbsoluteAndProxified(attributes []html.Attribute, proxyHostPort, remoteHost string) []html.Attribute {

	rew := make([]html.Attribute, 0, len(attributes))

	for i := 0; i < len(attributes); i++ {

		attr := attributes[i]

		// Make all absolute
		if attr.Key == "href" || attr.Key == "src" || attr.Key == "action" { //  make absolute

			if attrX(attributes, "cfrom") == "img" {
				attr.Val = absolutize("http", remoteHost, attr.Val)
			} else {
				attr.Val = absolutize("https", remoteHost, attr.Val)
			}

		}

		if attr.Key == "href" {

			if attrX(attributes, "cfrom") == "img" {
				// dont proxif image links
			} else {
				// proxify - v1
				attr.Val = fmt.Sprintf("%v?url=%v", routes.FetchUrl, attr.Val)

				if util_appengine.IsLocalEnviron() {
					attr.Val = fmt.Sprintf("http://%v%v", proxyHostPort, attr.Val)
				} else {
					attr.Val = fmt.Sprintf("https://%v%v", proxyHostPort, attr.Val)
				}
			}

		}

		if attr.Key == "action" {

			// proxify - v2

			// Since we appended a form field, we do not need:
			// action = spf("/blob2/form-redirector?redirect-to=%v", action)
			if util_appengine.IsLocalEnviron() {
				attr.Val = fmt.Sprintf("http://%v%v", proxyHostPort, routes.FormRedirector)
			} else {
				attr.Val = fmt.Sprintf("https://%v%v", proxyHostPort, routes.FormRedirector)
			}

		}

		rew = append(rew, attr)
	}

	//
	//
	// We instrumented all forms with a field "redirect-to"
	// Now we have to make the value of this field absolute
	isRedirectInput := false
	for _, attr := range rew {
		if attr.Key == "name" && attr.Val == "redirect-to" {
			isRedirectInput = true
		}
	}
	if isRedirectInput {
		for _, attr := range rew {
			if attr.Key == "value" {
				attr.Val = absolutize("https", remoteHost, attr.Val)
			}
		}
	}

	return rew
}

func proxify(n *html.Node, ProxyHostPort string, urlRemoteSrc *url.URL) {
	fRecurser := closuredProxifier(ProxyHostPort, urlRemoteSrc)
	fRecurser(n)
}
