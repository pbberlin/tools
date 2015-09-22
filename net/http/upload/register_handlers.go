package upload

import (
	"bytes"
	"net/http"

	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
)

const UrlUploadSend = "/filesys/post-upload-send"
const UrlUploadReceive = "/filesys/post-upload-receive"

func InitHandlers() {
	http.HandleFunc(UrlUploadSend, loghttp.Adapter(sendUpload))
	http.HandleFunc(UrlUploadReceive, loghttp.Adapter(receiveUpload))
	http.HandleFunc("/mnt00/", loghttp.Adapter(ServeDsFsFile))
	http.HandleFunc("/mnt01/", loghttp.Adapter(ServeDsFsFile))
	http.HandleFunc("/mnt02/", loghttp.Adapter(ServeDsFsFile))
}

// userinterface rendered to HTML - not only the strings for title and url
func BackendUIRendered() *bytes.Buffer {
	var b1 = new(bytes.Buffer)
	htmlfrag.Wb(b1, "Upload zip files into dsfs", "")
	htmlfrag.Wb(b1, "send", UrlUploadSend, "via command line or via this form")
	htmlfrag.Wb(b1, "receive", UrlUploadReceive, "receive a plain file or a zip archive")

	htmlfrag.Wb(b1, "serve file mnt00", "/mnt00/test.jpg", "")
	htmlfrag.Wb(b1, "serve file mnt01", "/mnt01/test.jpg", "")
	htmlfrag.Wb(b1, "serve file mnt02", "/mnt02/test.jpg", "")

	return b1
}
