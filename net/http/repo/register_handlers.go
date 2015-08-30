package repo

import (
	"bytes"
	"net/http"

	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
)

const uriSetType = "/fetch/set-fs-type"
const mountName = "mntftch"
const UriMountNameY = "/" + mountName + "/serve-file/"

const uriFetchCommandReceiver = "/fetch/command-receive"
const uriFetchCommandSender = "/fetch/command-send"

// InitHandlers is called from outside,
// and makes the EndPoints available.
func InitHandlers() {
	http.HandleFunc(uriSetType, loghttp.Adapter(setFSType))

	http.HandleFunc("/fetch/request-static", loghttp.Adapter(staticFetchDirect))

	http.HandleFunc(uriFetchCommandSender, loghttp.Adapter(staticFetchViaPosting2Receiver))
	http.HandleFunc(uriFetchCommandReceiver, loghttp.Adapter(fetchCommandReceiver))

	http.HandleFunc(UriMountNameY, loghttp.Adapter(serveFile))

	// working only for memfs
	http.Handle("/fetch/reservoire/static/", http.StripPrefix("/fetch/reservoire/static/", fileserver1))

}

// BackendUIRendered returns a userinterface rendered to HTML
func BackendUIRendered() *bytes.Buffer {
	var b1 = new(bytes.Buffer)
	htmlfrag.Wb(b1, "Fetcher", "")
	htmlfrag.Wb(b1, "type", uriSetType, "storage type")

	htmlfrag.Wb(b1, "static command", "/fetch/request-static", "send direct")

	htmlfrag.Wb(b1, "send command", uriFetchCommandSender, "dynamic")
	htmlfrag.Wb(b1, "recv", uriFetchCommandReceiver, "receive fetch command, takes commands by curl")

	htmlfrag.Wb(b1, "reservoire BOTH", UriMountNameY, "browse ANY fsi.FileSystem - human readable with ?fmt=http ")

	htmlfrag.Wb(b1, "reservoire static", "/fetch/reservoire/static/", "browse - memfs only")

	return b1
}
