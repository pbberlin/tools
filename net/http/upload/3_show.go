package upload

import (
	"bytes"
	"io/ioutil"
	"mime"
	"net/http"
	"path"

	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/os/fsi/aefs"

	"appengine"
)

func displayUpload(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	c := appengine.NewContext(r)

	b1 := new(bytes.Buffer)

	fclose := func() {
		// Only upon error.
		// If everything is fine, we reset fclose at the end.
		w.Write(b1.Bytes())
	}
	defer fclose()

	b1.WriteString(tplx.Head)
	wpf(b1, "<pre>\n")

	err := r.ParseForm()
	if err != nil {
		wpf(b1, "ParseFormErr %v", err)
		return
	}

	p := r.FormValue("path")
	wpf(b1, "path = %q\n", p)

	if len(p) > 0 {

		fs1 := aefs.New(
			aefs.MountName(aefs.MountPointLast()),
			aefs.AeContext(c),
		)

		fullP := path.Join(baseDirY, p)

		f, err := fs1.Open(fullP)
		if err != nil {
			wpf(b1, "err opening file %v - %v\n", fullP, err)
			return
		}

		inf, err := f.Stat()

		if err != nil {
			wpf(b1, "err opening fileinfo %v - %v\n", fullP, err)
			return
		}

		wpf(b1, "opened file %v - %v -  %v\n", f.Name(), inf.Size(), err)

		bts1, err := ioutil.ReadAll(f)

		tp := mime.TypeByExtension(path.Ext(fullP))

		w.Header().Set("Content-Type", tp)
		w.Write(bts1)

		b1 = new(bytes.Buffer) // reset the log

	}

}
