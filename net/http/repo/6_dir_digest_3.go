package repo

import (
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"time"

	"appengine"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/net/http/fetch"
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

		m.lg("\t\t file only %4.2v hours old, skipping", age.Hours())
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
	runtime.Gosched()
	bts, inf, err := fetch.UrlGetter(m.r, fetch.Options{URL: m.SURL, KnownProtocol: m.Protocol, RedirectHandling: 1})
	runtime.Gosched()

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
	var cx appengine.Context
	cx = util_appengine.SafelyExtractGaeContext(m.r)
	if cx == nil {
		m.lg("timed out - returning")
		return bts, inf.Mod, false, fmt.Errorf("req timed out")
	}

	m.lg("retrivd %q; %vkB ", inf.URL.Host+inf.URL.Path, len(bts)/1024)

	//
	//
	m.lg("saved   %q crawled file", fn)
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
