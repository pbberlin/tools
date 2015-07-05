// Package instance_info contains provides html views; the core is in instance_mgt.
package instance_info

import (
	"fmt"
	"net/http"

	"appengine"

	// sc "github.com/pbberlin/tools/dsu/distributed_unancestored"

	"github.com/pbberlin/tools/appengine/instance_mgt"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
)

func view(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	cntr := 1

	tplAdder, tplExec := tplx.FuncTplBuilder(w, r)
	tplAdder("n_html_title", "Application, Module and Instance Info", nil)
	tplAdder("n_cont_1", "<pre>{{.}}</pre>", instance_mgt.GetStatic().String())
	tplAdder("n_cont_2", "<p>{{.}} views</p>", cntr)
	tplAdder("n_cont_0", `
		<p>On the development server, call 
		<a href='/instance-info/collect' 
		target='collect' >collect</a> first.</p>

		<p><a href='/instance-info/`+instance_mgt.GetStatic().InstanceID+`'>specific url</a></p>
		
		`, "")

	tplExec(w, r)

	/*
	 Requests are routed randomly accross instances

	 Following is just a futile try to register
	 an instance specific handler.
	 It is only useful, when we request an instance
	 specifically via specific hostname
	*/
	SuppressPanicUponDoubleRegistration(
		w, r, "/instance-info/"+instance_mgt.GetStatic().InstanceID, loghttp.Adapter(view))

}

// SuppressPanicUponDoubleRegistration registers
// a request hanlder for a route.
//
//
// Because of asynchronicity we need to
// catch the ensuing panic for repeated registration
// of the same handler
func SuppressPanicUponDoubleRegistration(w http.ResponseWriter, r *http.Request,
	urlPattern string, handler func(http.ResponseWriter, *http.Request)) string {
	defer func() {
		panicSignal := recover()
		if panicSignal != nil {
			w.Write([]byte(fmt.Sprintf("panic caught:\n\n %s", panicSignal)))
		}
	}()

	http.HandleFunc(urlPattern, handler)
	return urlPattern

}

func collectInfo(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	instance_mgt.Get(appengine.NewContext(r), m)
}

func init() {

	// InstanceId := appengine.InstanceID() // does not during init, only after a few seconds

	http.HandleFunc("/instance-info/view", loghttp.Adapter(view))
	http.HandleFunc("/instance-info/collect", loghttp.Adapter(collectInfo))

}
