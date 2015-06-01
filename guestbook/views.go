package guestbook

import (
	"net/http"

	"appengine"

	sc "github.com/pbberlin/tools/dsu_distributed_unancestored"
	"github.com/pbberlin/tools/util_err"

	gbp "github.com/pbberlin/tools/dsu_ancestored_gb_entries" // guest book persistence
	"github.com/pbberlin/tools/tpl_html"
)

const c_view_gbe = `
	{{range .}}
		{{$a := .Date}}
		{{$b := .Date  | df }}
		{{$c := df .Date}}
			<p>
		{{with .Author}}
			<b>{{.}}</b> wrote on {{$c}}<br>
		{{else}}
			An anonymous person wrote:   <br>
		{{end}}
			<span style='display:block; 
			 white-space: pre;
			 max-width:300px;font-size:12px;' >{{.Content}}</span>
		</p>
	{{end}}
`

const c_new_gbe = `
	Try number {{.}}
	<form action="/guest-save" method="post">
		<div><textarea name="content" rows="3" cols="60"></textarea></div>
		<div><input type="submit" value="Save Entry"></div>
	</form>
`

func guestEntry(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)
	path := m["dir"].(string) + m["base"].(string)

	err := sc.Increment(c, path)
	util_err.Err_http(w, r, err, false)

	cntr, err := sc.Count(w, r, path)
	util_err.Err_http(w, r, err, false)

	tplAdder, tplExec := tpl_html.FuncTplBuilder(w, r)
	tplAdder("n_html_title", "New guest book entry", nil)
	tplAdder("n_cont_0", c_new_gbe, cntr)
	tplExec(w, r)

}

func guestSave(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	contnt := r.FormValue("content")
	gbp.SaveEntry(w, r, map[string]interface{}{"content": contnt})
	http.Redirect(w, r, "/guest-view", http.StatusFound)

}

func guestView(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)
	path := m["dir"].(string) + m["base"].(string)

	err := sc.Increment(c, path)
	util_err.Err_http(w, r, err, false)

	cntr, err := sc.Count(w, r, path)
	util_err.Err_http(w, r, err, false)

	gbEntries, report := gbp.ListEntries(w, r)

	tplAdder, tplExec := tpl_html.FuncTplBuilder(w, r)
	tplAdder("n_html_title", "List of guest book entries", nil)
	tplAdder("n_cont_0", c_view_gbe, gbEntries)
	tplAdder("n_cont_1", "<pre>{{.}}</pre>", report)
	tplAdder("n_cont_2", "Visitors: {{.}}", cntr)
	tplExec(w, r)

}
