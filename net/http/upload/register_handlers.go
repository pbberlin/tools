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
	http.HandleFunc("/blob2/get-file", loghttp.Adapter(displayUpload))
}

// userinterface rendered to HTML - not only the strings for title and url
func BackendUIRendered() *bytes.Buffer {
	var backendFragBlob = new(bytes.Buffer)
	htmlfrag.Wb(backendFragBlob, "Upload zip files into aefs", "")
	htmlfrag.Wb(backendFragBlob, "send", "/blob2/post-send", "via command line or via this form")
	htmlfrag.Wb(backendFragBlob, "receive", UrlUploadReceive, "receive a plain file or a zip archive")
	htmlfrag.Wb(backendFragBlob, "show", "/blob2/get-file", "show an aefs stored file by get param path")
	return backendFragBlob
}
