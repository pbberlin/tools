package gitkit1

/*
	This is a TWO step template creation process

	Hugo template is 'executed' into hugo-frame with gitkit body.
	Result is then again 'executed' to a template containing gitkit params.

*/

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
)

// Templates file path.
const (
	Headers = "\t" + `<script type="text/javascript"            src="//www.gstatic.com/authtoolkit/js/gitkit.js"></script>
	<link   type="text/css" rel="stylesheet" href="//www.gstatic.com/authtoolkit/css/gitkit.css">`
)

func GetHomeTpl(w http.ResponseWriter, r *http.Request, title, body string) *template.Template {

	if body == "" {
		body = IDCardHTML + UserInfoHTML
	}

	lg, _ := loghttp.BuffLoggerUniversal(w, r)

	bstpl := tplx.TemplateFromHugoPage(w, r)

	b := new(bytes.Buffer)

	fmt.Fprintf(b, tplx.ExecTplHelper(bstpl, map[string]interface{}{
		// "HtmlTitle":       "{{ .HtmlTitle }}", // this seems to cause problems sometimes.
		"HtmlTitle":       title,
		"HtmlDescription": "", // reminder
		"HtmlHeaders":     template.HTML(Headers),
		"HtmlContent":     template.HTML(body),
	}))

	intHomeTemplate, err := template.New("home").Parse(b.String())
	lg(err)

	return intHomeTemplate

}

func GetIDCardTpl(w http.ResponseWriter, r *http.Request, u *User) string {

	b := new(bytes.Buffer)

	fmt.Fprintf(b, tplx.ExecTplHelper(IDCardHTML, map[string]interface{}{
		"WidgetURL":  WidgetSigninAuthorizedRedirectURL,
		"SignOutURL": signOutURL,
		"User":       u,
		// "CookieDump": template.HTML(htmlfrag.CookieDump(r)),
	}))

	return b.String()
}

func GetWidgetTpl(w http.ResponseWriter, r *http.Request, title string) *template.Template {

	lg, _ := loghttp.BuffLoggerUniversal(w, r)

	bstpl := tplx.TemplateFromHugoPage(w, r) // the jQuery irritates
	// bstpl := tplx.HugoTplNoScript

	b := new(bytes.Buffer)
	fmt.Fprintf(b, tplx.ExecTplHelper(bstpl, map[string]interface{}{
		// "HtmlTitle":       "{{ .HtmlTitle }}",  // it DOES cause some eternal loop. But why only here?
		"HtmlTitle":       title,
		"HtmlDescription": "", // reminder
		"HtmlHeaders":     template.HTML(Headers),
		"HtmlContent":     template.HTML(widgetHTML),
	}))

	intGitkitTemplate, err := template.New("widg").Parse(b.String())
	lg(err)

	return intGitkitTemplate

}
