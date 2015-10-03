// Package fileserver replaces http.Fileserver
package fileserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/pbberlin/tools/net/http/htmlfrag"
	"github.com/pbberlin/tools/net/http/tplx"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/common"
	"github.com/pbberlin/tools/stringspb"
)

var wpf = func(w io.Writer, format string, a ...interface{}) (int, error) {
	fmt.Fprintf(w, format, a...)
	fmt.Fprintf(w, "\n")
	return 0, nil
}

var spf = fmt.Sprintf

// We cannot use http.FileServer(http.Dir("./css/")
// to dispatch our dsfs files.
// We need the appengine context to initialize dsfs.
// Thus we have to re-implement a serveFile method:
func FsiFileServer(fs fsi.FileSystem, prefix string, w http.ResponseWriter, r *http.Request) {

	r.Header.Set("X-Custom-Header-Counter", "nocounter")

	b1 := new(bytes.Buffer)

	fclose := func() {
		// Only upon error.
		// If everything is fine, we reset fclose at the end.
		w.Write(b1.Bytes())
	}
	defer fclose()

	wpf(b1, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "Half-Static-File-Server"}))
	wpf(b1, "<pre>")

	err := r.ParseForm()
	if err != nil {
		wpf(b1, "err parsing request (ParseForm)%v", err)
	}

	p := r.URL.Path

	if strings.HasPrefix(p, prefix) {
		// p = p[len(prefix):]
		p = strings.TrimPrefix(p, prefix)
	} else {
		wpf(b1, "route must start with prefix %v - but is %v", prefix, p)
	}

	if strings.HasPrefix(p, "/") {
		p = p[1:]
	}
	wpf(b1, "effective path = %q", p)

	// fullP := path.Join(docRootDataStore, p)
	fullP := p

	f, err := fs.Open(fullP)
	if err != nil {
		wpf(b1, "err opening file %v - %v", fullP, err)
		return
	}
	defer f.Close()

	inf, err := f.Stat()
	if err != nil {
		wpf(b1, "err opening fileinfo %v - %v", fullP, err)
		return
	}

	if inf.IsDir() {

		wpf(b1, "%v is a directory - trying index.html...", fullP)

		fullP += "/index.html"

		fIndex, err := fs.Open(fullP)
		if err == nil {
			defer fIndex.Close()
			inf, err = fIndex.Stat()
			if err != nil {
				wpf(b1, "err opening index fileinfo %v - %v", fullP, err)
				return
			}

			f = fIndex
		} else {

			wpf(b1, "err opening index file %v - %v", fullP, err)

			if r.FormValue("fmt") == "html" {
				dirListHtml(w, r, f)
			} else {
				dirListJson(w, r, f)
			}

			b1 = new(bytes.Buffer) // success => reset the message log => dumps an empty buffer
			return
		}

	}

	wpf(b1, "opened file %v - %v -  %v", f.Name(), inf.Size(), err)

	bts1, err := ioutil.ReadAll(f)
	if err != nil {
		wpf(b1, "err with ReadAll %v - %v", fullP, err)
		return
	}

	ext := path.Ext(fullP)
	ext = strings.ToLower(ext)
	tp := mime.TypeByExtension(ext)

	w.Header().Set("Content-Type", tp)

	//
	// caching
	// either explicitly discourage
	// or     explicitly  encourage
	if false ||
		ext == ".css" || ext == ".js" ||
		ext == "css" || ext == "js" ||
		ext == ".jpg" || ext == ".gif" ||
		ext == "jpg" || ext == "gif" ||
		false {

		if strings.Contains(fullP, "tamper-monkey") {
			htmlfrag.SetNocacheHeaders(w)
		} else {
			htmlfrag.CacheHeaders(w)
		}
	} else {
		htmlfrag.SetNocacheHeaders(w)
	}

	w.Write(bts1)

	b1 = new(bytes.Buffer) // success => reset the message log => dumps an empty buffer

}

// inspired by https://golang.org/src/net/http/fs.go
//
// name may contain '?' or '#', which must be escaped to remain
// part of the URL path, and not indicate the start of a query
// string or fragment.
var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",

	`"`, "&#34;",

	"'", "&#39;",
)

func dirListJson(w http.ResponseWriter, r *http.Request, f fsi.File) {

	r.Header.Set("Content-Type", "application/json")

	mp := []map[string]string{}

	for {
		dirs, err := f.Readdir(100)
		if err != nil || len(dirs) == 0 {
			break
		}
		for _, d := range dirs {
			name := d.Name()
			if d.IsDir() {
				name = common.Directorify(name)
			}
			name = htmlReplacer.Replace(name)

			url := url.URL{Path: name}

			mpl := map[string]string{
				"path": url.String(),
				"mod":  d.ModTime().Format(time.RFC1123Z),
			}

			mp = append(mp, mpl)
		}
	}

	bdirListHtml, err := json.MarshalIndent(mp, "", "\t")
	if err != nil {
		wpf(w, "marshalling to []byte failed - mp was %v", mp)
		return
	}
	w.Write(bdirListHtml)

}

func dirListHtml(w http.ResponseWriter, r *http.Request, f fsi.File) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	for {
		dirs, err := f.Readdir(100)
		if err != nil || len(dirs) == 0 {
			break
		}
		for _, d := range dirs {
			name := d.Name()

			suffix := ""
			if d.IsDir() {
				suffix = "/"
			}

			linktitle := htmlReplacer.Replace(name)
			linktitle = stringspb.Ellipsoider(linktitle, 40)
			if d.IsDir() {
				linktitle = common.Directorify(linktitle)
			}

			surl := path.Join(r.URL.Path, name) + suffix + "?fmt=html"

			oneLine := spf("<a  style='display:inline-block;min-width:600px;' href=\"%s\">%s</a>", surl, linktitle)
			// wpf(w, " %v", d.ModTime().Format("2006-01-02 15:04:05 MST"))
			oneLine += spf(" %v<br>", d.ModTime().Format(time.RFC1123Z))
			wpf(w, oneLine)
		}
	}

}
