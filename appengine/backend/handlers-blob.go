package backend

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
)

var backendFragBlob = new(bytes.Buffer)

func init() {

	//
	// handler registration
	http.HandleFunc("/blob2/zipupload", loghttp.Adapter(receiveUpload))
	htmlfrag.Wb(backendFragBlob, "upload-receive", "/blob2/zipupload")

	http.HandleFunc("/blob2/zipdisplay", loghttp.Adapter(displayUpload))
	htmlfrag.Wb(backendFragBlob, "upload-show", "/blob2/zipdisplay")

}

func receiveUpload(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	// Here the parameter is the size of the form data that should
	// be loaded in memory, the remaining being put in temporary
	// files
	r.ParseMultipartForm(1024 * 1024 * 2)

	fields := []string{"title", "author", "description"}
	for _, v := range fields {
		s := spf("%12v => %q", v, r.FormValue(v))
		fmt.Fprintln(w, s)
	}

	ff := "filefield"

	file, handler, err := r.FormFile(ff)
	if err != nil {
		log.Println(err)
	}
	if handler == nil {
		fmt.Fprintf(w, "no multipart file %q\n", ff)
	} else {
		fmt.Fprintf(w, "extracted file %v\n", handler.Filename)

		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, "extracted file content;  %v bytes\n", len(data))

		newFilename := "c:\\temp\\xx" + handler.Filename
		err = ioutil.WriteFile(newFilename, data, 0777)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, "saved file content to %v\n", newFilename)
	}

}

func displayUpload(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
}
