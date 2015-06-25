// Package tplx makes composing templates more comfortable.
package tplx

import (
	"fmt"
	"html"
	tt "html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pbberlin/tools/net/http/loghttp"
)

/*
We a assume a fixed html page
with embedded detail templates.


We give http handler functions a tiny interface
to construct and execute html templates.


Three and a half prefixes:
   n_ ... is the string *name*    of a template
   c_ ... is the string *content* of a template
		PrefixLff prompts loading of the template content
		from directory static/[x].html
   t_ ... is a *parsed template*



Usage:

type GBEntry struct {
    Author  string
    Content string
}

const c_tpl_gbentry = `
		{{range .}}
			<b>{{.Author}}</b> wrote:
			{{.Content}} <br>
		{{end}}
`

	gbe1 := GBEntry{
		Content: "It crawls into a man's bowels ...",
		Author:   "John Dos Passos",
	}
	gbe2 := GBEntry{
		Content: "Praised be the man ...",
		Author:   "T.S. Elliot",
	}
	vGbe := []GBEntry{gbe1,gbe2}

	myTemplateAdder, myTplExec := FuncTplBuilder(w,r)

	myTemplateAdder("n_html_title","What authors think","")  // no dyn data
	myTemplateAdder("n_cont_0","<i>{{.}}</i><br><br>","Deep thought:")
	myTemplateAdder("n_cont_0", PrefixLff + "reloaded_template"," dyn data for file ")
	myTemplateAdder("n_cont_1",c_tpl_gbentry ,vGbe)
	myTemplateAdder("n_cont_2","<br><br>end of thoughts","")

	myTplExec(w,r)



*/

const FurtherDefinition = "Further_Definition_"

var dbg bool = false

const c_page_scaffold_01 = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
	<link rel="icon" href="data:;base64,=">
	<title>{{template "n_html_title"  }}</title>
	<link rel="stylesheet" type="text/css" href="http://fonts.googleapis.com/css?family=Open Sans">
  </head>
  <link rel="stylesheet" href="/static/basic.css" media="screen" type="text/css" />
  <body>
	<div class="wmax" style='margin: 0px auto; '>
		{{template "n_cont_0" .n_cont_0}}
		{{template "n_cont_1" .n_cont_1 }}
		{{template "n_cont_2" .n_cont_2 }}
		{{template "n_cont_3" .n_cont_3 }}
		<p><a href='/'>Back to root</a></p>
	</div>
  </body>
</html>`

// the defautls map must contain all subtemplates
// demanded by c_page_scaffold_01
var map_default map[string]string = map[string]string{
	"n_html_title": "",
	"n_cont_0":     "",
	"n_cont_1":     "",
	"n_cont_2":     "",
	"n_cont_3":     "",
}

const c_contentpl_extended = `
	--{{.}}--
`

const PrefixLff string = "load_file_" // prefix load from file
var lenLff int = len(PrefixLff)

//

var t_base *tt.Template = nil

// a clone factory
//  as t_base is a "app scope" variable
//  it will live as long as the application runs
func cloneFromBase(w http.ResponseWriter, r *http.Request) *tt.Template {

	funcMap := tt.FuncMap{
		"unescape":         html.UnescapeString,
		"escape":           html.EscapeString,
		"fMult":            fMult,
		"fAdd":             fAdd,
		"fChop":            fChop,
		"fAccessElement":   fAccessElement,
		"fAccessElementB2": fAccessElementB2,
		"fMakeRange":       fMakeRange,
		"fNumCols":         fNumCols,
		"df": func(g time.Time) string {
			return g.Format("2006-01-02 (Jan 02)")
		},
	}

	if t_base == nil {
		t_base = tt.Must(tt.New("n_page_scaffold_01").Funcs(funcMap).Parse(c_page_scaffold_01))
	}

	t_derived, err := t_base.Clone()
	loghttp.E(w, r, err, false)

	return t_derived
}

// m contains defaults or overwritten template contents
func templatesExtend(w http.ResponseWriter, r *http.Request, subtemplates map[string]string) *tt.Template {

	var err error = nil
	tder := cloneFromBase(w, r)

	//for k,v := range furtherDefinitions{
	//	tder, err = tder.Parse( `{{define "` + k  +`"}}`   + v + `{{end}}` )
	//	loghttp.E(w,r,err,false)
	//}

	for k, v := range subtemplates {

		if len(v) > lenLff && v[0:lenLff] == PrefixLff {
			fileName := v[lenLff:]
			fcontent, err := ioutil.ReadFile("templates/" + fileName + ".html")
			loghttp.E(w, r, err, false, "could not open static template file")
			v = string(fcontent)
		}

		tder, err = tder.Parse(`{{define "` + k + `"}}` + v + `{{end}}`)
		loghttp.E(w, r, err, false)

	}

	return tder
}

// returns two functions
//   first  function   collects detail templates and detail content data
//   second function   combines all those and executes them
// example usage
// myTemplateAdder, myTplExec := FuncTplBuilder()
// myTemplateAdder("n_content","--{{.}}--","some dyn data string")
// myTplExec(w,r)

func FuncTplBuilder(w http.ResponseWriter, r *http.Request) (f1 func(string, string, interface{}),
	f2 func(http.ResponseWriter, *http.Request)) {

	// prepare collections
	mtc := map[string]string{}      // map template content
	mtd := map[string]interface{}{} // map template data

	// template key - template content, template data
	f1 = func(tk string, tc string, td interface{}) {
		mtc[tk] = tc
		mtd[tk] = td
	}

	f2 = func(w http.ResponseWriter, r *http.Request) {

		// merge arguments with defaults
		map_result := map[string]string{}
		for k, v := range map_default {
			if _, ok := mtc[k]; ok {
				map_result[k] = mtc[k]
				delete(mtc, k)
			} else {
				map_result[k] = v
			}
			if dbg {
				w.Write([]byte(fmt.Sprintf("  %q  %q \n", map_result[k], v)))
			}
		}

		// additional templates beyond the default
		for k, v := range mtc {
			map_result[k] = v
			if dbg {
				w.Write([]byte(fmt.Sprintf("  %q  %q \n", map_result[k], v)))
			}
		}

		tpl_extended := templatesExtend(w, r, map_result)

		err := tpl_extended.ExecuteTemplate(w, "n_page_scaffold_01", mtd)
		loghttp.E(w, r, err, false)

	}

	return f1, f2
}
