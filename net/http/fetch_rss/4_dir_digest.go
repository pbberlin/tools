package fetch_rss

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/osutilpb"
	"golang.org/x/net/html"
)

type DirsWithFiles struct {
	Name      string // Name == key of Parent.Dirs
	LastFound time.Time
	Dirs      map[string]DirsWithFiles
	// Fils []string
}

func dirsWithFilesStr(buf *bytes.Buffer, d *DirsWithFiles, lvl int) {
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
		dirsWithFilesStr(buf, &indir, lvl+1)
	}
}

func (d DirsWithFiles) String() string {
	buf := new(bytes.Buffer)
	dirsWithFilesStr(buf, &d, 0)
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

func crawl(w http.ResponseWriter, r *http.Request, fs fsi.FileSystem, c FetchCommand) (DirsWithFiles, error) {

	lg, lge := loghttp.Logger(w, r)

	crawl1URL := path.Join(c.Host, path.Dir(c.SearchPrefix))
	lg("crawl %v", crawl1URL)

	dwf := DirsWithFiles{Name: "root1", Dirs: map[string]DirsWithFiles{}, LastFound: time.Now()}
	dwfLoop := dwf

	fnDigest := path.Join(docRoot, c.Host, "digest2.json")

	if c.Host == "test.economist.com" {
		switchTData(w, r)
	}

	bdwf, err := fs.ReadFile(fnDigest)
	lge(err)
	if err == nil {
		err = json.Unmarshal(bdwf, &dwf)
		lge(err)
	}

	var crawl2URL *url.URL
	bts, crawl2URL, err := fetch.UrlGetter(r, fetch.Options{URL: crawl1URL})
	lge(err)
	if err != nil {
		return dwf, err
	}

	lg("retrieved %v; %vkB ", crawl2URL.String(), len(bts)/1024)

	doc, err := html.Parse(bytes.NewReader(bts))
	lge(err)

	anchors := []string{}
	var fr func(*html.Node)
	fr = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			href := attrX(n.Attr, "href")
			anchors = append(anchors, href)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fr(c)
		}
	}
	fr(doc)

	pfx1 := "http://" + crawl2URL.Host
	pfx2 := "https://" + crawl2URL.Host
	for _, href := range anchors {
		href = strings.TrimPrefix(href, pfx1)
		href = strings.TrimPrefix(href, pfx2)
		if strings.HasPrefix(href, "/") { // ignore other domains
			parsed, err := url.Parse(href)
			lge(err)
			href = parsed.Path
			// lg("%v", href)
			dwfLoop = dwf
			dir, remainder := "", href
			for {

				dir, remainder = osutilpb.PathDirReverse(remainder)
				dwfLoop.Name = dir
				dwfLoop.LastFound = time.Now()

				if _, ok := dwfLoop.Dirs[dir]; !ok {
					dwfLoop.Dirs[dir] = DirsWithFiles{Name: dir, Dirs: map[string]DirsWithFiles{}}
				}

				dwfLoop = dwfLoop.Dirs[dir]

				// Since we "cannot assign" to map struct directly:
				// dwfLoop.Dirs[dir].LastFound = time.Now()   // fails
				dwfLoop.LastFound = time.Now()

				if remainder == "" {
					break
				}
			}

		}
	}

	b, err := json.MarshalIndent(dwf, "", "\t")
	lge(err)

	dir := path.Join(docRoot, crawl2URL.Host)
	err = fs.MkdirAll(dir, 0755)
	lge(err)

	err = fs.WriteFile(fnDigest, b, 0755)
	lge(err)

	return dwf, nil

}
