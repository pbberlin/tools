package upload

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"appengine"

	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi/aefs"
	"github.com/pbberlin/tools/os/fsi/fsc"
)

var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var pfRestore func(format string, a ...interface{}) (int, error) = fmt.Printf
var spf func(format string, a ...interface{}) string = fmt.Sprintf
var wpf func(w io.Writer, format string, a ...interface{}) (int, error) = fmt.Fprintf

func InitHandlers() {
	http.HandleFunc("/blob2/zipupload", loghttp.Adapter(receiveUpload))
	http.HandleFunc("/blob2/zipdisplay", loghttp.Adapter(displayUpload))
}

func GetBackend() *bytes.Buffer {
	var backendFragBlob = new(bytes.Buffer)
	htmlfrag.Wb(backendFragBlob, "upload-receive", "/blob2/zipupload")
	htmlfrag.Wb(backendFragBlob, "upload-show", "/blob2/zipdisplay")
	return backendFragBlob
}

func receiveUpload(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	r.ParseMultipartForm(1024 * 1024 * 2)

	fields := []string{"title", "author", "description"}
	for _, v := range fields {
		s := spf("%12v => %q", v, r.FormValue(v))
		fmt.Fprintln(w, s)
		c.Infof(s)
	}

	ff := "filefield"

	file, handler, err := r.FormFile(ff)
	if err != nil {
		log.Println(err)
	}
	if handler == nil {
		fmt.Fprintf(w, "no multipart file %q\n", ff)
		c.Infof("no multipart file %q\n", ff)
	} else {
		fmt.Fprintf(w, "extracted file %v\n", handler.Filename)

		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, "extracted file content;  %v bytes\n", len(data))
		c.Infof("extracted file %q ;  %v bytes\n", handler.Filename, len(data))

		newFilename := "c:\\TEMP\\xxx" + handler.Filename

		err = ioutil.WriteFile(newFilename, data, 0777)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, "saved file content to %v - %v\n", newFilename, err)

		fs1 := aefs.New(
			aefs.MountName(aefs.MountPointLast()),
			aefs.AeContext(c),
		)

		dir, bname := fsc.PathInternalize(newFilename, fs1.RootDir(), fs1.RootName())

		err = fs1.MkdirAll(dir, 0777)
		fmt.Fprintf(w, "mkdir %v - %v\n", dir, err)

		err = fs1.WriteFile(path.Join(dir, bname), data, 0777)
		fmt.Fprintf(w, "saved file content to %v - %v\n", newFilename, err)

	}

}

func displayUpload(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
}
