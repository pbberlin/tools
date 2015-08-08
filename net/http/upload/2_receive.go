package upload

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/os/fsi/aefs"

	"archive/zip"

	"appengine"
)

func receiveUpload(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, _ := loghttp.Logger(w, r)
	c := appengine.NewContext(r)

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Receive an Upload"}))
	defer wpf(w, tplx.Foot)
	wpf(w, "<pre>")
	defer wpf(w, "</pre>")

	err := r.ParseMultipartForm(1024 * 1024 * 2)
	if err != nil {
		lg("Multipart parsing failed: %v", err)
		return
	}

	fields := []string{"getparam1", "mountname", "description"}
	for _, v := range fields {
		lg("%12v => %q", v, r.FormValue(v))
	}

	mountPoint := aefs.MountPointLast()
	if len(r.FormValue("mountname")) > 0 {
		mountPoint = r.FormValue("mountname")
	}
	lg("mount point is %v", mountPoint)

	fs1 := aefs.New(
		aefs.MountName(mountPoint),
		aefs.AeContext(c),
	)

	// As closure, since we cannot define aefs.aeFileSys as parameter
	funcSave := func(argName string, data []byte) (error, *bytes.Buffer) {

		b1 := new(bytes.Buffer)

		fs1 := aefs.New(
			aefs.MountName(mountPoint),
			aefs.AeContext(c),
		)

		dir, bname := fs1.SplitX(argName)

		err := fs1.MkdirAll(dir, 0777)
		wpf(b1, "mkdir %v - %v\n", dir, err)
		if err != nil {
			return err, b1
		}

		err = fs1.WriteFile(path.Join(dir, bname), data, 0777)
		wpf(b1, "saved file content to %v - %v\n", argName, err)

		return err, b1
	}

	ff := "filefield"

	file, handler, err := r.FormFile(ff)
	if err != nil {
		lg("error calling FormFile from %q  => %v", ff, err)
		return
	}

	if handler == nil {
		lg("no multipart file %q", ff)
	} else {
		lg("extracted file %v", handler.Filename)

		data, err := ioutil.ReadAll(file)
		if err != nil {
			lg("ReadAll on uploaded file failed: %v", err)
			return
		}
		defer file.Close()
		lg("extracted file content;  %v bytes", len(data))

		newFilename := docRootDataStore + handler.Filename
		ext := path.Ext(newFilename)

		if ext == ".zip" {

			lg("found zip - treat as dir-tree %q", newFilename)

			r, err := zip.NewReader(file, int64(len(data)))
			if err != nil {
				lg("open as zip failed: %v", err)
				return
			}

			for _, f := range r.File {
				newFilename = docRootDataStore + f.Name

				dir, bname := fs1.SplitX(newFilename)

				if f.FileInfo().IsDir() {

					lg("\t dir %s", newFilename)

					err := fs1.MkdirAll(path.Join(dir, bname), 0777)
					if err != nil {
						lg("MkdirAll %v failed: %v", newFilename, err)
						return
					}

				} else {

					lg("\t file %s", newFilename)

					rc, err := f.Open()
					if err != nil {
						return
					}
					defer func(rc io.ReadCloser) {
						if err := rc.Close(); err != nil {
							panic(err)
						}
					}(rc)

					bts := new(bytes.Buffer)
					size, err := io.Copy(bts, rc)
					if err != nil {
						lg("Could not copy from zipped file %v: %v", newFilename, err)
						return
					}

					err = fs1.WriteFile(path.Join(dir, bname), bts.Bytes(), 0777)
					if err != nil {
						lg("WriteFile of zipped file %v failed: %v", newFilename, err)
						return
					}
					lg("\t  saved %v - %v Bytes", newFilename, size)

				}

			}

		} else {

			err, b2 := funcSave(newFilename, data)
			lg("%s", b2)
			if err != nil {
				return
			}

		}
		lg("--------------------\n")

	}

}
