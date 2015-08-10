package upload

import (
	"bytes"
	"net/http"

	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
)

const UrlUploadReceive = "/blob2/post-receive"

func InitHandlers() {
	http.HandleFunc("/blob2/post-send", loghttp.Adapter(sendUpload))
	http.HandleFunc(UrlUploadReceive, loghttp.Adapter(receiveUpload))
	http.HandleFunc("/mnt00/", loghttp.Adapter(serveFile))
	http.HandleFunc("/mnt01/", loghttp.Adapter(serveFile))
	http.HandleFunc("/mnt02/", loghttp.Adapter(serveFile))

}

// userinterface rendered to HTML - not only the strings for title and url
func BackendUIRendered() *bytes.Buffer {
	var backendFragBlob = new(bytes.Buffer)
	htmlfrag.Wb(backendFragBlob, "Upload zip files into aefs", "")
	htmlfrag.Wb(backendFragBlob, "send", "/blob2/post-send", "via command line or via this form")
	htmlfrag.Wb(backendFragBlob, "receive", UrlUploadReceive, "receive a plain file or a zip archive")

	htmlfrag.Wb(backendFragBlob, "serve file mnt00", "/mnt00/test.jpg", "")
	htmlfrag.Wb(backendFragBlob, "serve file mnt01", "/mnt01/test.jpg", "")
	htmlfrag.Wb(backendFragBlob, "serve file mnt02", "/mnt02/test.jpg", "")

	return backendFragBlob
}
