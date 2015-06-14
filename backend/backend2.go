package backend

import (
	"net/http"

	sc "github.com/pbberlin/tools/dsu/distributed_unancestored"
	"github.com/pbberlin/tools/tpl_html"
	"github.com/pbberlin/tools/util_appengine"
	"github.com/pbberlin/tools/util_err"

	"appengine"
)

type link struct {
	Linktext string
	Url      string
	Target   string
}

var links1 = []interface{}{
	link{Linktext: "Login", Url: "/login"},
	link{Linktext: "Schreib-Methoden", Url: "/write-methods"},
	link{Linktext: "Letzte Email", Url: "/email-view"},
	link{Linktext: "Blob List", Url: "/blob/list"},
	link{Linktext: "Template Demo 1", Url: "/tpl/demo1"},
	link{Linktext: "Template Demo 2", Url: "/tpl/demo2"},
	link{Linktext: "Http fetch", Url: "/fetch-url"},
	link{Linktext: "Instance Info", Url: "/instance-info/view"},
	link{Linktext: "Gob encode decode", Url: "/big-query/test-gob-codec"},

	link{Linktext: "JSON encode", Url: "/json-encode"},
	link{Linktext: "JSON decode", Url: "/json-decode"},

	link{Linktext: "Fulltext put", Url: "/fulltext-search/put"},
	link{Linktext: "Fulltext get", Url: "/fulltext-search/get"},
}
var links2 = []interface{}{
	link{Linktext: "Eintrag hinzufügen", Url: "/guest-entry"},
	link{Linktext: "Einträge auflisten", Url: "/guest-view"},
	link{Linktext: "Einträge auflisten - paged - serialized cursor", Url: "/guest-view-cursor"},
}

var blocks2 = map[string]interface{}{
	"01 Diverse": []interface{}{
		link{Linktext: "Login", Url: "/login"},
		link{Linktext: "Schreib-Methoden", Url: "/write-methods"},
		link{Linktext: "Letzte Email", Url: "/email-view"},
		link{Linktext: "Blob List", Url: "/blob/list"},
		link{Linktext: "Template Demo 1", Url: "/tpl/demo1"},
		link{Linktext: "Template Demo 2", Url: "/tpl/demo2"},
		link{Linktext: "Http fetch", Url: "/fetch-url"},
		link{Linktext: "Instance Info", Url: "/instance-info/view"},
		link{Linktext: "Gob encode decode", Url: "/big-query/test-gob-codec"},

		link{Linktext: "JSON encode", Url: "/json-encode"},
		link{Linktext: "JSON decode", Url: "/json-decode"},

		link{Linktext: "Fulltext put", Url: "/fulltext-search/put"},
		link{Linktext: "Fulltext get", Url: "/fulltext-search/get"},
	},
	"02 Guestbook": []interface{}{
		link{Linktext: "Eintrag hinzufügen", Url: "/guest-entry"},
		link{Linktext: "Einträge auflisten", Url: "/guest-view"},
		link{Linktext: "Einträge auflisten - paged - serialized cursor", Url: "/guest-view-cursor"},
	},
	"03 Drawing": []interface{}{
		link{Linktext: "Drawing a static chart", Url: "/image/draw-lines-example"},
	},
	"04 Big Query": []interface{}{
		link{Linktext: "Get real data", Url: "/big-query/query-into-datastore"},
		link{Linktext: "Get mocked data", Url: "/big-query/mock-data-into-datastore"},
	},
	"05 ... with Chart": []interface{}{
		link{Linktext: "Process Data 1 (mock=1},", Url: "/big-query/regroup-data-01?mock=0"},
		link{Linktext: "Process Data 2", Url: "/big-query/regroup-data-02?f=table"},
		link{Linktext: "Show as Table", Url: "/big-query/show-table"},
		link{Linktext: "Show as Chart", Url: "/big-query/show-chart"},
		link{Linktext: "As HTML", Url: "/big-query/html"},
	},
	"06 Request Images": []interface{}{
		link{Linktext: "WrapBlob from Datastore", Url: "/image/img-from-datastore?p=chart1"},
		link{Linktext: "base64 from Datastore", Url: "/image/base64-from-datastore?p=chart1"},
		link{Linktext: "base64 from Variable", Url: "/image/base64-from-var?p=1"},
		link{Linktext: "base64 from File", Url: "/image/base64-from-file?p=static/pberg1.png"},
	},
	"07 Namespaces + Task Queues": []interface{}{
		link{Linktext: "Increment", Url: "/namespaced-counters/increment"},
		link{Linktext: "Read", Url: "/namespaced-counters/read"},
		link{Linktext: "Push to task-queue", Url: "/namespaced-counters/queue-push"},
	},
	"08 URLs with/without ancestors": []interface{}{
		link{Linktext: "Backend", Url: "/save-url/backend"},
	},
}

func backend2(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	path := m["dir"].(string) + m["base"].(string)

	err := sc.Increment(c, path)
	util_err.Err_http(w, r, err, false)

	cntr, err := sc.Count(w, r, path)
	util_err.Err_http(w, r, err, false)

	add, tplExec := tpl_html.FuncTplBuilder(w, r)
	add("n_html_title", "Backend", nil)
	//add("n_cont_0", c_link, links)
	add("n_cont_0", tpl_html.PrefixLff+"backend_body", blocks2)
	add("tpl_legend", tpl_html.PrefixLff+"backend_body_embed01", "")

	//add("n_cont_1", "<pre>{{.}}</pre>", "pure text")
	add("n_cont_2", "<p>{{.}} views</p>", cntr)

	tplExec(w, r)

}

func init() {
	http.HandleFunc("/backend2", util_appengine.Adapter(backend2))

}
