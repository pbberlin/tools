package upload

import (
	"bytes"
	"net/http"

	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
)

func InitHandlers() {
	http.HandleFunc("/blob2/zipupload", loghttp.Adapter(receiveUpload))
	http.HandleFunc("/blob2/zipdisplay", loghttp.Adapter(displayUpload))
}

// userinterface rendered to HTML - not only the strings for title and url
func BackendUIRendered() *bytes.Buffer {
	var backendFragBlob = new(bytes.Buffer)
	htmlfrag.Wb(backendFragBlob, "Upload zip files into aefs", "")
	htmlfrag.Wb(backendFragBlob, "receive", "/blob2/zipupload", "receive a plain file or a zip archive")
	htmlfrag.Wb(backendFragBlob, "show", "/blob2/zipdisplay", "show an aefs stored file by get param path")
	return backendFragBlob
}
