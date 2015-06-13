package backend

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/pbberlin/tools/conv"
	htmlpb "github.com/pbberlin/tools/pbhtml"
	"github.com/pbberlin/tools/pbstrings"
	"github.com/pbberlin/tools/util"
	"github.com/pbberlin/tools/util_appengine"
	"github.com/pbberlin/tools/util_err"
)

func backend(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	b1 := new(bytes.Buffer)

	htmlpb.Wb(b1, "Diverse", "")
	htmlpb.Wb(b1, "Login", "/login")
	htmlpb.Wb(b1, "Schreib-Methoden", "/write-methods")
	htmlpb.Wb(b1, "Letzte Email", "/email-view")
	htmlpb.Wb(b1, "Blob List", "/blob/list")
	htmlpb.Wb(b1, "Template Demo 1", "/tpl/demo1")
	htmlpb.Wb(b1, "Template Demo 2", "/tpl/demo2")
	htmlpb.Wb(b1, "Http fetch", "/fetch-url")
	htmlpb.Wb(b1, "Instance Info", "/instance-info/view")
	htmlpb.Wb(b1, "Gob encode decode", "/big-query/test-gob-codec")

	htmlpb.Wb(b1, "JSON encode", "/json-encode")
	htmlpb.Wb(b1, "JSON decode", "/json-decode")

	htmlpb.Wb(b1, "Fulltext put", "/fulltext-search/put")
	htmlpb.Wb(b1, "Fulltext get", "/fulltext-search/get")

	htmlpb.Wb(b1, "datastore object view quoted printabe", "/dsu/show")

	htmlpb.Wb(b1, "Guest Book", "")
	htmlpb.Wb(b1, "Eintrag hinzufügen", "/guest-entry")
	htmlpb.Wb(b1, "Einträge auflisten", "/guest-view")
	htmlpb.Wb(b1, "Einträge auflisten - paged - serialized cursor", "/guest-view-cursor")

	htmlpb.Wb(b1, " ", "")
	htmlpb.Wb(b1, "Drawing a static chart", "/image/draw-lines-example")

	htmlpb.Wb(b1, "Big Query ...", "")
	htmlpb.Wb(b1, "Get real data", "/big-query/query-into-datastore")
	htmlpb.Wb(b1, "Get mocked data", "/big-query/mock-data-into-datastore")
	htmlpb.Wb(b1, "  &nbsp; &nbsp; &nbsp; ... with Chart", "")
	htmlpb.Wb(b1, "Process Data 1 (mock=1)", "/big-query/regroup-data-01?mock=0")
	htmlpb.Wb(b1, "Process Data 2", "/big-query/regroup-data-02?f=table")
	htmlpb.Wb(b1, "Show as Table", "/big-query/show-table")
	htmlpb.Wb(b1, "Show as Chart", "/big-query/show-chart")
	htmlpb.Wb(b1, "As HTML", "/big-query/html")

	htmlpb.Wb(b1, "Request Images ", "")
	htmlpb.Wb(b1, "WrapBlob from Datastore", "/image/img-from-datastore?p=chart1")
	htmlpb.Wb(b1, "base64 from Datastore", "/image/base64-from-datastore?p=chart1")
	htmlpb.Wb(b1, "base64 from Variable", "/image/base64-from-var?p=1")
	htmlpb.Wb(b1, "base64 from File", "/image/base64-from-file?p=static/pberg1.png")

	htmlpb.Wb(b1, "Namespaces + Task Queues", "")
	htmlpb.Wb(b1, "Increment", "/namespaced-counters/increment")
	htmlpb.Wb(b1, "Read", "/namespaced-counters/read")
	htmlpb.Wb(b1, "Push to task-queue", "/namespaced-counters/queue-push")

	htmlpb.Wb(b1, "URLs with/without ancestors", "")
	htmlpb.Wb(b1, "Backend", "/save-url/backend")

	htmlpb.Wb(b1, "Statistics", "/_ah/stats")

	b1.WriteString("<br>\n")
	b1.WriteString("<hr>\n")
	b1.WriteString("<a target='_gae' href='https://console.developers.google.com/project/347979071940' ><b>global</b> developer console</a><br>\n")
	b1.WriteString(" &nbsp; &nbsp; <a target='_gae' href='http://localhost:8000/mail' >app console local</a><br>\n")
	b1.WriteString(" &nbsp; &nbsp; <a target='_gae' href='https://appengine.google.com/settings?&app_id=s~libertarian-islands' >app console online</a><br>\n")

	b1.WriteString(` &nbsp; &nbsp; <a target='_gae' 
			href='http://go-lint.appspot.com/github.com/pbberlin/tools/dsu' 
			>lint package</a><br>`)

	b1.WriteString("<br>\n")
	b1.WriteString("<a target='_gae'   href='http://localhost:8085/' >app local</a><br>\n")
	b1.WriteString("<a target='_gae_r' href='http://libertarian-islands.appspot.com/' >app online</a><br>\n")

	dir := m["dir"].(string)
	base := m["base"].(string)
	b1.WriteString("<br>\n")
	b1.WriteString("Dir: --" + dir + "-- &nbsp; &nbsp; &nbsp; &nbsp;   Base: --" + base + "-- <br>\n")

	b1.WriteString("<br>\n")
	s := fmt.Sprintf("IntegerSequenes a, b: %v %v %v<br>\n", util_err.MyIntSeq01(), util_err.MyIntSeq01(), util_err.MyIntSeq02())
	b1.WriteString(s)

	// b1.WriteString("<br>\n")
	// b1.WriteString(fmt.Sprintf("Temp dir is %s<br>\n", os.TempDir()))

	b1.WriteString("<br>\n")
	b2 := new(bytes.Buffer)
	b2.WriteString("data:image/png;base64,...")
	b1.WriteString(fmt.Sprintf("Mime from %q is %q<br>\n", b2.String(),
		conv.MimeFromBase64(b2)))

	b1.WriteString("<br>\n")

	io.WriteString(b1, "Date: "+util.TimeMarker()+"  - ")
	b1.WriteString(fmt.Sprintf("Last Month %q - 24 Months ago is %q<br>\n", util.MonthsBack(0),
		util.MonthsBack(24)))

	b1.WriteString("<br>\n")
	x1 := " z" + pbstrings.IncrementString("--z")
	x2 := " Z" + pbstrings.IncrementString("--Z")
	x3 := " 9" + pbstrings.IncrementString("--9")
	x4 := " Peter" + pbstrings.IncrementString("--Peter")
	sEnc := "Theo - wir fahrn nach Łódź <  " + pbstrings.IncrementString("Łódź") + x1 + x2 + x3 + x4
	b1.WriteString(fmt.Sprint("restore string string(  []byte(sEnc) ): ", string([]byte(sEnc)), "<br>"))

	w.Header().Set("Content-Type", "text/html")
	w.Write(b1.Bytes())

}

func init() {
	http.HandleFunc("/backend", util_appengine.Adapter(backend))

}
