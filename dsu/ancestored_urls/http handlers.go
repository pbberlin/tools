package ancestored_urls

import (
	"net/http"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
)

func init() {
	http.HandleFunc("/save-url/save-no-anc", loghttp.Adapter(saveURLNoAnc))
	http.HandleFunc("/save-url/save-wi-anc", loghttp.Adapter(saveURLWithAncestor))
	http.HandleFunc("/save-url/view-no-anc", loghttp.Adapter(listURLNoAnc))
	http.HandleFunc("/save-url/view-wi-anc", loghttp.Adapter(listURLWithAncestors))
	http.HandleFunc("/save-url/backend", loghttp.Adapter(backend))
	http.HandleFunc("/save-url/", loghttp.Adapter(backend))

}

func backend(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	add, tplExec := tplx.FuncTplBuilder(w, r)

	add("n_html_title", "Saving an URL into the datastore", "")
	//add("n_cont_0", tplx.PrefixLff+"body_dsu_ancestored_urls", "")
	add("n_cont_0", tplx.PrefixLff+"body_last_url", "")

	tplExec(w, r)

}
