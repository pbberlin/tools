package repo

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

	"github.com/golang/snappy"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/osutilpb"
	"github.com/pbberlin/tools/stringspb"
	"golang.org/x/net/html"
)

func dirTreeStrRec(buf *bytes.Buffer, d *DirTree, lvl int) {
	ind2 := strings.Repeat("    ", lvl+1)
	keys := make([]string, 0, len(d.Dirs))
	for k, _ := range d.Dirs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		buf.WriteString(ind2)
		indir := d.Dirs[key]
		buf.WriteString(stringspb.ToLen(indir.Name, 44-len(ind2)))
		if indir.EndPoint {
			buf.WriteString(fmt.Sprintf(" EP"))
		}
		buf.WriteByte(10)
		dirTreeStrRec(buf, &indir, lvl+1)
	}
}

func (d DirTree) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(d.Name)
	// buf.WriteString(fmt.Sprintf(" %v ", len(d.Dirs)))
	if d.Dirs == nil {
		buf.WriteString(" (nil)")
	}
	buf.WriteByte(10)
	dirTreeStrRec(buf, &d, 0)
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

func path2DirTree(lg loghttp.FuncBufUniv, treeX *DirTree, articles []FullArticle, domain string, IsRSS bool) {

	if treeX == nil {
		return
	}
	var trLp *DirTree
	trLp = treeX

	pfx1 := "http://" + domain
	pfx2 := "https://" + domain

	for _, art := range articles {
		href := art.Url
		if art.Mod.IsZero() {
			art.Mod = time.Now()
		}
		href = strings.TrimPrefix(href, pfx1)
		href = strings.TrimPrefix(href, pfx2)
		if strings.HasPrefix(href, "/") { // ignore other domains
			parsed, err := url.Parse(href)
			lg(err)
			href = parsed.Path
			// lg("%v", href)
			trLp = treeX
			// lg("trLp is %v", trLp.String())
			dir, remainder, remDirs := "", href, []string{}
			lvl := 0
			for {

				dir, remainder, remDirs = osutilpb.PathDirReverse(remainder)

				if dir == "/" && remainder == "" {
					// skip root
					break
				}

				if lvl > 0 {
					trLp.Name = dir // lvl==0 => root
				}
				trLp.LastFound = art.Mod.Truncate(time.Minute)

				// lg("   %v, %v", dir, remainder)

				// New creation
				if _, ok := trLp.Dirs[dir]; !ok {
					if IsRSS {
						trLp.Dirs[dir] = DirTree{Name: dir, Dirs: map[string]DirTree{}, SrcRSS: true}
					} else {
						trLp.Dirs[dir] = DirTree{Name: dir, Dirs: map[string]DirTree{}}
					}
				}

				// We "cannot assign" to map struct directly:
				// trLp.Dirs[dir].LastFound = art.Mod   // fails with "cannot assign"
				addressable := trLp.Dirs[dir]
				addressable.LastFound = art.Mod.Truncate(time.Minute)

				// We can rely that the *last* dir or html is an endpoint.
				// We cannot tell about higher paths, unless explicitly linked somewhere
				// Previous distinction between RSS URLs and crawl URLs dropped
				if len(remDirs) < 1 {
					addressable.EndPoint = true
				}

				if dir == "/2015" || dir == "/08" || dir == "/09" {
					addressable.EndPoint = true
				}

				trLp.Dirs[dir] = addressable
				trLp = &addressable

				if remainder == "" {
					// lg("break\n")
					break
				}

				lvl++
			}

		}
	}

}

// Append of all links of a DOM to an in-memory dirtree
func addAnchors(lg loghttp.FuncBufUniv, host string, bts []byte, dirTree *DirTree) {

	doc, err := html.Parse(bytes.NewReader(bts))
	lg(err)
	if err != nil {
		return
	}
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
	path2DirTree(lg, dirTree, anchors, host, false)
	lg("\t\tadded %v anchors", len(anchors))
	dirTree.LastFound = time.Now() // Marker for later accumulated saving

}

func loadDigest(w http.ResponseWriter, r *http.Request, lg loghttp.FuncBufUniv, fs fsi.FileSystem, fnDigest string, treeX *DirTree) {

	fnDigestSnappied := strings.Replace(fnDigest, ".json", ".json.snappy", -1)
	bts, err := fs.ReadFile(fnDigestSnappied)
	if err == nil {
		btsDec := []byte{}
		lg("encoded digest loaded, size %vkB", len(bts)/1024)
		btsDec, err := snappy.Decode(nil, bts)
		if err != nil {
			lg(err)
			return
		}
		lg("digest decoded from %vkB to %vkB", len(bts)/1024, len(btsDec)/1024)
		bts = btsDec
	} else {
		bts, err = fs.ReadFile(fnDigest)
		lg(err)
	}

	if err == nil {
		err = json.Unmarshal(bts, &treeX)
		lg(err)
	}

	lg("DirTree   %5.2vkB loaded for %v", len(bts)/1024, fnDigest)

}

// requesting via http; not from filesystem
// unused
func fetchDigest(hostWithPrefix, domain string) (*DirTree, error) {

	lg, lge := loghttp.Logger(nil, nil)
	_ = lg

	surl := path.Join(hostWithPrefix, domain, "digest2.json")
	bts, _, err := fetch.UrlGetter(nil, fetch.Options{URL: surl})
	lge(err)
	if err != nil {
		return nil, err
	}

	// lg("%s", bts)
	dirTree := &DirTree{Name: "/", Dirs: map[string]DirTree{}, EndPoint: true}

	if err == nil {
		err = json.Unmarshal(bts, dirTree)
		lge(err)
		if err != nil {
			return nil, err
		}
	}

	lg("DirTree   %5.2vkB loaded for %v", len(bts)/1024, surl)

	age := time.Now().Sub(dirTree.LastFound)
	lg("DirTree is %5.2v hours old (%v)", age.Hours(), dirTree.LastFound.Format(time.ANSIC))

	return dirTree, nil

}

func saveDigest(lg loghttp.FuncBufUniv, fs fsi.FileSystem, fnDigest string, treeX *DirTree) {

	treeX.LastFound = time.Now()

	b, err := json.MarshalIndent(treeX, "", "\t")
	lg(err)

	if len(b) > 1024*1024-1 || true {
		b1 := snappy.Encode(nil, b)
		lg("digest encoded from %vkB to %vkB ", len(b)/1024, len(b1)/1024)
		b = b1
		fnDigest = strings.Replace(fnDigest, ".json", ".json.snappy", -1)
	}

	err = fs.MkdirAll(path.Dir(fnDigest), 0755)
	lg(err)

	err = fs.WriteFile(fnDigest, b, 0755)
	lg(err)

}
