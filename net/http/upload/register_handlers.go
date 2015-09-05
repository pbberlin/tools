package upload

import (
	"bytes"
	"net/http"

	"appengine"

	"github.com/pbberlin/tools/net/http/fileserver"
	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi/dsfs"
	"github.com/pbberlin/tools/os/fsi/memfs"
)

const UrlUploadReceive = "/blob2/post-receive"

func InitHandlers() {
	http.HandleFunc("/blob2/post-send", loghttp.Adapter(sendUpload))
	http.HandleFunc(UrlUploadReceive, loghttp.Adapter(receiveUpload))
	http.HandleFunc("/mnt00/", loghttp.Adapter(ServeDsFsFile))
	http.HandleFunc("/mnt01/", loghttp.Adapter(ServeDsFsFile))
	// http.HandleFunc("/mnt02/", loghttp.Adapter(ServeDsFsFile))

	var fs1 = memfs.New(
		memfs.Ident("mnt02"),
	)
	dynSrv := func(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
		appID := appengine.AppID(appengine.NewContext(r))
		if appID == "credit-expansion" {

			prefix := "/mnt02"
			// prefix = "/xxx"

			fs2 := dsfs.New(
				dsfs.MountName(prefix[1:]),
				dsfs.AeContext(appengine.NewContext(r)),
			)

			fs1.SetOption(
				memfs.ShadowFS(fs2),
			)

			fileserver.FsiFileServer(fs1, prefix+"/", w, r)
		} else {
			ServeDsFsFile(w, r, m)
			w.Write([]byte("chalamacuca"))
		}
	}
	dmpMemfs := func(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
		htmlfrag.SetNocacheHeaders(w, false)
		w.Write([]byte("<pre>"))
		w.Write(fs1.Dump())
	}
	http.HandleFunc("/mnt02/", loghttp.Adapter(dynSrv))
	http.HandleFunc("/memfsdmp", loghttp.Adapter(dmpMemfs))
	// http.HandleFunc("/", loghttp.Adapter(dynSrv))

}

// userinterface rendered to HTML - not only the strings for title and url
func BackendUIRendered() *bytes.Buffer {
	var b1 = new(bytes.Buffer)
	htmlfrag.Wb(b1, "Upload zip files into dsfs", "")
	htmlfrag.Wb(b1, "send", "/blob2/post-send", "via command line or via this form")
	htmlfrag.Wb(b1, "receive", UrlUploadReceive, "receive a plain file or a zip archive")

	htmlfrag.Wb(b1, "serve file mnt00", "/mnt00/test.jpg", "")
	htmlfrag.Wb(b1, "serve file mnt01", "/mnt01/test.jpg", "")
	htmlfrag.Wb(b1, "serve file mnt02", "/mnt02/test.jpg", "")

	return b1
}
