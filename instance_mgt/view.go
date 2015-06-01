package instance_mgt

import (
	"net/http"

	// sc "github.com/pbberlin/tools/dsu_distributed_unancestored"

	"github.com/pbberlin/tools/tpl_html"
	"github.com/pbberlin/tools/util_appengine"
	"github.com/pbberlin/tools/util_err"

	"appengine"
)

func view(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)
	_ = c

	// CIRCULAR DEPENDENCY PROBLEM

	// path := m["dir"].(string) + m["base"].(string)

	// err := sc.Increment(c, path)
	// util_err.Err_http(w, r, err, false)

	// cntr, err := sc.Count(w, r, path)
	// util_err.Err_http(w, r, err, false)
	cntr := 1

	tplAdder, tplExec := tpl_html.FuncTplBuilder(w, r)
	tplAdder("n_html_title", "Application, Module and Instance Info", nil)
	tplAdder("n_cont_1", "<pre>{{.}}</pre>", ii.String())
	tplAdder("n_cont_2", "<p>{{.}} views</p>", cntr)
	tplAdder("n_cont_0", `
		<p>On the development server, call 
		<a href='/instance-info/collect' 
		target='collect' >collect</a> first.</p>

		<p><a href='/instance-info/`+ii.InstanceID+`'>specific url</a></p>
		
		`, "")

	tplExec(w, r)

	/*
	 Requests are routed randomly accross instances

	 Following is just a futile try to register
	 an instance specific handler.
	 It is only useful, when we request an instance
	 specifically via specific hostname
	*/
	util_err.SuppressPanicUponDoubleRegistration(w, r, "/instance-info/"+ii.InstanceID,
		util_appengine.Adapter(view))

}
