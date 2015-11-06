package repo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"golang.org/x/net/context"
	"golang.org/x/net/html"
)

// Fetches URL if local file is outdated.
// saves fetched file
//
// link extraction, link addition to treeX now accumulated one level higher
// bool return value: use existing => true
func fetchSave(m *MyWorker) ([]byte, time.Time, bool, error) {

	// w http.ResponseWriter,
	// r *http.Request,

	// Determine FileName
	ourl, err := fetch.URLFromString(m.SURL)
	fc := FetchCommand{}
	fc.Host = ourl.Host
	fc = addDefaults(fc)
	semanticUri := condenseTrailingDir(m.SURL, fc.CondenseTrailingDirs)
	fn := path.Join(docRoot, semanticUri)

	m.lg("crawlin %q", m.SURL)

	// File already exists?
	// Open file for age check
	var bts []byte
	var mod time.Time
	f := func() error {
		file1, err := m.fs1.Open(fn)
		// m.lg(err) // file may simply not exist
		if err != nil {
			return err // file may simply not exist
		}
		defer file1.Close() // file close *fast* at the end of *this* anonymous func

		fi, err := file1.Stat()
		m.lg(err)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			m.lg("\t\t file is a directory, skipping - %v", fn)
			return fmt.Errorf("is directory: %v", fn)
		}

		mod = fi.ModTime()
		age := time.Now().Sub(mod)
		if age.Hours() > 10 {
			m.lg("\t\t file %4.2v hours old, refetch ", age.Hours())
			return fmt.Errorf("too old: %v", fn)
		}

		m.lg("\t\t file only %4.2v hours old, take %4.2vkB from datastore", age.Hours(), fi.Size()/1024)
		bts, err = ioutil.ReadAll(file1)
		if err != nil {
			return err
		}
		return nil
	}

	err = f()
	if err == nil {
		return bts, mod, true, err
	}

	//
	// Fetch
	bts, inf, err := fetch.UrlGetter(m.r, fetch.Options{URL: m.SURL, KnownProtocol: m.Protocol, RedirectHandling: 1})

	m.lg(err)
	if err != nil {
		m.lg("tried to fetch %v, %v", m.SURL, inf.URL)
		m.lg("msg %v", inf.Msg)
		return []byte{}, inf.Mod, false, err
	}
	if inf.Mod.IsZero() {
		inf.Mod = time.Now().Add(-75 * time.Minute)
	}

	//
	//
	// main request still exists?
	if false {
		var cx context.Context
		cx = util_appengine.SafelyExtractGaeContext(m.r)
		if cx == nil {
			m.lg("timed out - returning")
			return bts, inf.Mod, false, fmt.Errorf("req timed out")
		}
	}

	m.lg("retrivd+saved %q; %vkB ", inf.URL.Host+inf.URL.Path, len(bts)/1024)

	if len(bts) > 1024*1024-1 {
		bts = removeScriptsAndComments(m.lg, bts)
		m.lg("size reduzed to %vkB ", len(bts)/1024)

		if len(bts) > 1024*1024-1 {
			bts = zipIt(m.lg, bts)
			m.lg("size reduzed to %vkB ", len(bts)/1024)
		}

	}

	//
	//
	dir := path.Dir(fn)
	err = m.fs1.MkdirAll(dir, 0755)
	m.lg(err)
	err = m.fs1.Chtimes(dir, time.Now(), time.Now())
	m.lg(err)
	err = m.fs1.WriteFile(fn, bts, 0644)
	m.lg(err)
	err = m.fs1.Chtimes(fn, inf.Mod, inf.Mod)
	m.lg(err)

	return bts, inf.Mod, false, nil

}

func removeScriptsAndComments(lg loghttp.FuncBufUniv, bts []byte) []byte {
	doc, err := html.Parse(bytes.NewReader(bts))
	lg(err)
	if err != nil {
		return []byte{}
	}
	var fr func(*html.Node) // function recursive
	fr = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fr(c)
		}
		removeUnwanted(n)

	}
	fr(doc)
	var b bytes.Buffer
	err = html.Render(&b, doc)
	return b.Bytes()
}

func removeUnwanted(n *html.Node) {
	cc := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cc = append(cc, c)
	}
	for _, c := range cc {
		if n.Type == html.ElementNode && n.Data == "script" || n.Type == html.CommentNode {
			n.RemoveChild(c)
		}
	}
}

func zipIt(lg loghttp.FuncBufUniv, bts []byte) []byte {
	var b bytes.Buffer
	return b.Bytes()
}
