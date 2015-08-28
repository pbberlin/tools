package fetch_rss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi"
	"golang.org/x/net/html"
)

type DirTree struct {
	Name      string // Name == key of Parent.Dirs
	LastFound time.Time
	Dirs      map[string]DirTree
	// Fils []string
}

func DirTreeStr(buf *bytes.Buffer, d *DirTree, lvl int) {
	ind2 := strings.Repeat("    ", lvl+1)
	keys := make([]string, 0, len(d.Dirs))
	for k, _ := range d.Dirs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		buf.WriteString(ind2)
		indir := d.Dirs[key]
		buf.WriteString(indir.Name)
		buf.WriteByte(10)
		DirTreeStr(buf, &indir, lvl+1)
	}
}

func (d DirTree) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(d.Name)
	buf.WriteString(fmt.Sprintf(" %v ", len(d.Dirs)))
	if d.Dirs == nil {
		buf.WriteString(" (nil)")
	}
	buf.WriteByte(10)
	DirTreeStr(buf, &d, 0)
	return buf.String()
}

func switchTData(w http.ResponseWriter, r *http.Request) {

	lg, lge := loghttp.Logger(w, r)
	_ = lge

	b := fetch.TestData["test.economist.com"]
	sub1 := []byte(`<li><a href="/sections/newcontinent">xxx</a></li>`)

	sub2 := []byte(`<li><a href="/sections/asia">Asia</a></li>`)
	sub3 := []byte(`<li><a href="/sections/asia">Asia</a></li>
		<li><a href="/sections/newcontinent">xxx</a></li>`)

	if bytes.Contains(b, sub1) {
		b = bytes.Replace(b, sub1, []byte{}, -1)
	} else {
		b = bytes.Replace(b, sub2, sub3, -1)
	}

	if bytes.Contains(b, sub1) {
		lg("now contains %s", sub1)
	} else {
		lg("NOT contains %s", sub1)
	}

	fetch.TestData["test.economist.com"] = b

}

func loadDigest(w http.ResponseWriter, r *http.Request, fs fsi.FileSystem, fnDigest string, treeX *DirTree) {

	lg, lge := loghttp.Logger(w, r)
	_ = lg

	bts, err := fs.ReadFile(fnDigest)
	lge(err)
	if err == nil {
		err = json.Unmarshal(bts, &treeX)
		lge(err)
	}

}

func saveDigest(w http.ResponseWriter, r *http.Request, fs fsi.FileSystem, fnDigest string, treeX *DirTree) {

	lg, lge := loghttp.Logger(w, r)
	_ = lg

	b, err := json.MarshalIndent(treeX, "", "\t")
	lge(err)

	err = fs.MkdirAll(path.Dir(fnDigest), 0755)
	lge(err)

	err = fs.WriteFile(fnDigest, b, 0755)
	lge(err)

}

func crawl(w http.ResponseWriter, r *http.Request, treeX *DirTree, fs fsi.FileSystem, c FetchCommand) error {

	lg, lge := loghttp.Logger(w, r)

	if treeX == nil {
		treeX = &DirTree{Name: "root1", Dirs: map[string]DirTree{}, LastFound: time.Now()}
	}

	crawl1URL := path.Join(c.Host, path.Dir(c.SearchPrefix))
	lg("crawl %v", crawl1URL)

	var crawl2URL *url.URL
	bts, crawl2URL, err := fetch.UrlGetter(r, fetch.Options{URL: crawl1URL})
	lge(err)
	if err != nil {
		return err
	}

	lg("retrieved %v; %vkB ", crawl2URL.String(), len(bts)/1024)

	doc, err := html.Parse(bytes.NewReader(bts))
	lge(err)

	anchors := []FullArticle{}
	var fr func(*html.Node)
	fr = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			art := FullArticle{}
			art.Url = attrX(n.Attr, "href")
			art.Mod = time.Now()
			anchors = append(anchors, art)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fr(c)
		}
	}
	fr(doc)

	path2DirTree(w, r, treeX, anchors, crawl2URL.Host)

	return nil

}
