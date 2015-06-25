package tplx

import (
	"net/http"

	"github.com/pbberlin/tools/net/http/loghttp"
)

type GBEntry struct {
	Author  string
	Content string
}

const c_tpl_gbentry = `
		{{range .}}
			- <b>{{.Author}}</b> wrote: 
			{{.Content}} <br>
		{{end}}
`

func templatesDemo2(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	gbe1 := GBEntry{
		Content: "It crawls into a man's bowels ...",
		Author:  "John Dos Passos",
	}
	gbe2 := GBEntry{
		Content: "Praised be the man ...",
		Author:  "T.S. Elliot",
	}
	vGbe := []GBEntry{gbe1, gbe2}

	myTemplateAdder, myTplExec := FuncTplBuilder(w, r)

	myTemplateAdder("n_html_title", "What authors think", "") // no dyn data
	myTemplateAdder("n_cont_0", "Header <i>{{.}}</i><br><br>", "Deep thought:")
	myTemplateAdder("n_cont_0", PrefixLff+"reloaded_template", "dyn data for file")
	myTemplateAdder("n_cont_1", c_tpl_gbentry, vGbe)
	myTemplateAdder("n_cont_2", "<br>end of thoughts", "")

	myTplExec(w, r)

}

func init() {
	http.HandleFunc("/tpl/demo2", loghttp.Adapter(templatesDemo2))
}
