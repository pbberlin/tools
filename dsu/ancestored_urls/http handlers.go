package ancestored_urls

import (
	"net/http"

	"github.com/pbberlin/tools/net/http/tplx"
)

func init() {
	http.HandleFunc("/save-url/save-no-anc", saveURLNoAnc)
	http.HandleFunc("/save-url/save-wi-anc", saveURLWithAncestor)
	http.HandleFunc("/save-url/view-no-anc", listURLNoAnc)
	http.HandleFunc("/save-url/view-wi-anc", listURLWithAncestors)
	http.HandleFunc("/save-url/backend", backend)
	http.HandleFunc("/save-url/", backend)

}

func backend(w http.ResponseWriter, r *http.Request) {

	add, tplExec := tplx.FuncTplBuilder(w, r)

	add("n_html_title", "Saving an URL into the datastore", "")
	//add("n_cont_0", tplx.PrefixLff+"body_dsu_ancestored_urls", "")
	add("n_cont_0", tplx.PrefixLff+"body_last_url", "")

	tplExec(w, r)

}
