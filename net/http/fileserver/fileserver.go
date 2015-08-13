// Package fileserver replaces http.Fileserver
package fileserver

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/stringspb"
)

var wpf = fmt.Fprintf

// We cannot use http.FileServer(http.Dir("./css/")
// to dispatch our dsfs files.
// We need the appengine context to initialize dsfs.
// Thus we have to re-implement a serveFile method:
func FsiFileServer(fs fsi.FileSystem, prefix string, w http.ResponseWriter, r *http.Request) {

	b1 := new(bytes.Buffer)

	fclose := func() {
		// Only upon error.
		// If everything is fine, we reset fclose at the end.
		w.Write(b1.Bytes())
	}
	defer fclose()

	wpf(b1, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Half-Static-File-Server"}))
	wpf(b1, "<pre>\n")

	p := r.URL.Path

	if strings.HasPrefix(p, prefix) {
		p = p[len(prefix):]
	} else {
		wpf(b1, "route must start with prefix %v - but is %v\n", prefix, p)
	}

	if strings.HasPrefix(p, "/") {
		p = p[1:]
	}
	wpf(b1, "effective path = %q \n", p)

	// fullP := path.Join(docRootDataStore, p)
	fullP := p

	f, err := fs.Open(fullP)
	if err != nil {
		wpf(b1, "err opening file %v - %v\n", fullP, err)
		return
	}

	inf, err := f.Stat()
	if err != nil {
		wpf(b1, "err opening fileinfo %v - %v\n", fullP, err)
		return
	}

	if inf.IsDir() {

		wpf(b1, "%v is a directory - trying index.html...\n", fullP)

		fullP += "/index.html"

		fIndex, err := fs.Open(fullP)
		if err == nil {
			inf, err = fIndex.Stat()
			if err != nil {
				wpf(b1, "err opening index fileinfo %v - %v\n", fullP, err)
				return
			}

			f = fIndex
		} else {

			wpf(b1, "err opening index file %v - %v\n", fullP, err)
			dirList(w, f)
			b1 = new(bytes.Buffer) // success => reset the message log => dumps an empty buffer
			return
		}

	}

	wpf(b1, "opened file %v - %v -  %v\n", f.Name(), inf.Size(), err)

	bts1, err := ioutil.ReadAll(f)
	if err != nil {
		wpf(b1, "err with ReadAll %v - %v\n", fullP, err)
		return
	}

	tp := mime.TypeByExtension(path.Ext(fullP))

	w.Header().Set("Content-Type", tp)
	w.Write(bts1)

	b1 = new(bytes.Buffer) // success => reset the message log => dumps an empty buffer

}

// inspired by https://golang.org/src/net/http/fs.go

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",

	`"`, "&#34;",

	"'", "&#39;",
)

func dirList(w http.ResponseWriter, f fsi.File) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	for {
		dirs, err := f.Readdir(100)
		if err != nil || len(dirs) == 0 {
			break
		}
		for _, d := range dirs {
			name := d.Name()
			if d.IsDir() {
				name += "/"
			}
			// name may contain '?' or '#', which must be escaped to remain
			// part of the URL path, and not indicate the start of a query
			// string or fragment.
			linktitle := htmlReplacer.Replace(name)
			linktitle = stringspb.Ellipsoider(linktitle, 40)

			url := url.URL{Path: name}
			wpf(w, "<a  style='display:inline-block;min-width:600px;' href=\"%s\">%s</a>", url.String(), linktitle)
			wpf(w, " %v\n", d.ModTime().Format("2006-01-02 15:04:05 MST"))
			wpf(w, "<br>\n")
		}
	}

}
