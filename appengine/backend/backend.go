// Package backend ties together
// the visual controls and utils.
package backend

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/pbberlin/tools/appengine/login"
	"github.com/pbberlin/tools/net/http/coinbase"
	"github.com/pbberlin/tools/net/http/dedup"
	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/repo"
	"github.com/pbberlin/tools/net/http/routes"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/net/http/upload"
	// _ "github.com/pbberlin/tools/os/fsi/dsfs"
	"github.com/pbberlin/tools/os/fsi/webapi"
	"github.com/pbberlin/tools/stringspb"
	"github.com/pbberlin/tools/util"
)

var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var pfRestore func(format string, a ...interface{}) (int, error) = fmt.Printf
var spf func(format string, a ...interface{}) string = fmt.Sprintf
var wpf func(w io.Writer, format string, a ...interface{}) (int, error) = fmt.Fprintf

func init() {
	upload.InitHandlers()
	webapi.InitHandlers()
	repo.InitHandlers()
	dedup.InitHandlers()
	coinbase.InitHandlers()
	tplx.InitHandlers()
	login.InitHandlers()

}

func backend(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if ok, _, msg := login.CheckForAdminUser(r); !ok {
		w.Write([]byte(msg))
		return
	}

	b1 := new(bytes.Buffer)
	b1.WriteString(tplx.ExecTplHelper(tplx.Head, map[string]interface{}{"HtmlTitle": "Backend V1"}))

	htmlfrag.Wb(b1, "Debug pprof", "/debug/pprof")

	htmlfrag.Wb(b1, "Diverse", "nobr")
	htmlfrag.Wb(b1, "Schreib-Methoden", "/write-methods")
	htmlfrag.Wb(b1, "Letzte Email", "/email-view")
	htmlfrag.Wb(b1, "Blob List", "/blob2")

	htmlfrag.Wb(b1, "fetch via proxy", routes.ProxifyURI)
	htmlfrag.Wb(b1, "Instance Info", "/instance-info/view")

	htmlfrag.Wb(b1, "Fulltext put", "/fulltext-search/put")
	htmlfrag.Wb(b1, "Fulltext get", "/fulltext-search/get")

	htmlfrag.Wb(b1, "datastore object view quoted printabe", "/dsu/show")

	htmlfrag.Wb(b1, "Statistics", "/_ah/stats")

	htmlfrag.Wb(b1, "Request Images ", "")
	htmlfrag.Wb(b1, "WrapBlob from Datastore", "/image/img-from-datastore?p=chart1")
	htmlfrag.Wb(b1, "base64 from Datastore", "/image/base64-from-datastore?p=chart1")
	htmlfrag.Wb(b1, "base64 from Variable", "/image/base64-from-var?p=1")
	htmlfrag.Wb(b1, "base64 from File", "/image/base64-from-file?p=static/pberg1.png")
	htmlfrag.Wb(b1, "Drawing a static chart", "/image/draw-lines-example")

	htmlfrag.Wb(b1, "Big Query ...", "")
	htmlfrag.Wb(b1, "Get real data", "/big-query/query-into-datastore")
	htmlfrag.Wb(b1, "Get mocked data", "/big-query/mock-data-into-datastore")
	htmlfrag.Wb(b1, "  &nbsp; &nbsp; &nbsp; ... with Chart", "")
	htmlfrag.Wb(b1, "Process Data 1 (mock=1)", "/big-query/regroup-data-01?mock=0")
	htmlfrag.Wb(b1, "Process Data 2", "/big-query/regroup-data-02?f=table")
	htmlfrag.Wb(b1, "Show as Table", "/big-query/show-table")
	htmlfrag.Wb(b1, "Show as Chart", "/big-query/show-chart")
	htmlfrag.Wb(b1, "As HTML", "/big-query/html")

	htmlfrag.Wb(b1, "Namespaces + Task Queues", "")
	htmlfrag.Wb(b1, "Increment", "/namespaced-counters/increment")
	htmlfrag.Wb(b1, "Read", "/namespaced-counters/read")
	htmlfrag.Wb(b1, "Push to task-queue", "/namespaced-counters/queue-push")

	htmlfrag.Wb(b1, "URLs with/without ancestors", "nobr")
	htmlfrag.Wb(b1, "Backend", "/save-url/backend")

	htmlfrag.Wb(b1, "Guest Book", "")
	htmlfrag.Wb(b1, "Eintrag hinzufügen", "/guest-entry")
	htmlfrag.Wb(b1, "Einträge auflisten", "/guest-view")
	htmlfrag.Wb(b1, "Einträge auflisten - paged - serialized cursor", "/guest-view-cursor")

	b1.WriteString("<hr>\n")

	uiDsFs := webapi.BackendUIRendered()
	b1.Write(uiDsFs.Bytes())

	b1.WriteString("<hr>\n")

	b1.Write(upload.BackendUIRendered().Bytes())

	b1.Write(repo.BackendUIRendered().Bytes())

	b1.Write(dedup.BackendUIRendered().Bytes())

	b1.Write(coinbase.BackendUIRendered().Bytes())

	b1.Write(tplx.BackendUIRendered().Bytes())

	b1.Write(login.BackendUIRendered().Bytes())

	b1.WriteString("<br>\n")
	b1.WriteString("<hr>\n")

	urlLocalAdmin := fmt.Sprintf("http://localhost:%v/mail", routes.DevAdminPort())
	ancLocalAdmin := fmt.Sprintf(" &nbsp; &nbsp; <a target='_gae' href='%v' >local app console</a><br>\n", urlLocalAdmin)
	b1.WriteString(ancLocalAdmin)

	urlConsole := fmt.Sprintf("https://console.developers.google.com/project/%v", routes.AppID())
	ancConsole := fmt.Sprintf("<a target='_gae' href='%v' ><b>global</b> developer console</a>\n", urlConsole)
	b1.WriteString(ancConsole)

	urlOldAdmin := fmt.Sprintf("https://appengine.google.com/settings?&app_id=s~%v", routes.AppID())
	ancOldAdmin := fmt.Sprintf(" &nbsp; &nbsp; <a target='_gae' href='%v' >old admin UI</a><br>\n ", urlOldAdmin)
	b1.WriteString(ancOldAdmin)

	b1.WriteString(` &nbsp; &nbsp; <a target='_gae' 
			href='http://go-lint.appspot.com/github.com/pbberlin/tools/dsu' 
			>lint a package</a><br>`)

	dir := m["dir"].(string)
	base := m["base"].(string)
	b1.WriteString("<br>\n")
	b1.WriteString("Dir: --" + dir + "-- &nbsp; &nbsp; &nbsp; &nbsp;   Base: --" + base + "-- <br>\n")

	b1.WriteString("<br>\n")
	s := fmt.Sprintf("IntegerSequenes a, b: %v %v %v<br>\n", util.MyIntSeq01(), util.MyIntSeq01(), util.MyIntSeq02())
	b1.WriteString(s)

	// b1.WriteString("<br>\n")
	// b1.WriteString(fmt.Sprintf("Temp dir is %s<br>\n", os.TempDir()))

	b1.WriteString("<br>\n")

	io.WriteString(b1, "Date: "+util.TimeMarker()+"  - ")
	b1.WriteString(fmt.Sprintf("Last Month %q - 24 Months ago is %q<br>\n", util.MonthsBack(0),
		util.MonthsBack(24)))

	b1.WriteString("<br>\n")
	x1 := " z" + stringspb.IncrementString("--z")
	x2 := " Z" + stringspb.IncrementString("--Z")
	x3 := " 9" + stringspb.IncrementString("--9")
	x4 := stringspb.IncrementString(" --Peter")
	sEnc := "Łódź <  " + stringspb.IncrementString("Łódź") + x1 + x2 + x3 + x4
	b1.WriteString(fmt.Sprint(string([]byte(sEnc)), "<br>"))

	b1.WriteString(tplx.Foot)

	w.Write(b1.Bytes())

}

func init() {
	http.HandleFunc("/backend", loghttp.Adapter(backend))
}
