package fetch_rss

import (
	"bytes"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/osutilpb"
	"golang.org/x/net/html"
)

func crawl(w http.ResponseWriter, r *http.Request, c FetchCommand) {

	lg, lge := loghttp.Logger(w, r)

	crawlURL := path.Join(c.Host, path.Dir(c.SearchPrefix))
	lg("crawl %v", crawlURL)

	// bts, rssUrlObj, err := fetch.UrlGetter(r, fetch.Options{URL: crawlURL})

	rssUrlObj := url.URL{}
	rssUrlObj.Host = "www.economist.com"
	rssUrlObj.Scheme = "http"
	bts := economistHomepage
	var err error

	lge(err)
	if err == nil {

		lg("retrieved %v; %vkB ", rssUrlObj.String(), len(bts)/1024)

		doc, err := html.Parse(bytes.NewReader(bts))
		lge(err)
		_ = doc

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

		type DirsWithFiles struct {
			Name string
			Dirs map[string]DirsWithFiles
			// Fils []string
		}

		dwf := DirsWithFiles{Name: "root"}
		dwfLoop := dwf
		// dwf.Dirs = map[string]DirsWithFiles{}

		pfx1 := "http://" + rssUrlObj.Host
		pfx2 := "https://" + rssUrlObj.Host
		for _, href := range anchors {
			href = strings.TrimPrefix(href, pfx1)
			href = strings.TrimPrefix(href, pfx2)
			if strings.HasPrefix(href, "/") { // ignore other domains
				lg("%v", href)

				dwfLoop = dwf
				dir, remainder := "", href
				for {
					dir, remainder = osutilpb.PathDirReverse(remainder)
					lg("     fnd %q %q", dir, remainder)

					dwfLoop.Name = dir
					if dwfLoop.Dirs == nil {
						dwfLoop.Dirs = map[string]DirsWithFiles{}
					}

					if _, ok := dwfLoop.Dirs[dir]; !ok {
						lg("      added %v", dir)
						dwfLoop.Dirs[dir] = DirsWithFiles{Name: dir}
					}

					if remainder == "" {
						lg("     breaking")
						lg("     %v", dwfLoop.Dirs)
						break
					}
					// dwfLoop = dwfLoop.Dirs[dir]

				}

				lg("dump %v", dwf.Name)
				for k, v := range dwf.Dirs {
					lg("%v %v [%v]", k, v, dwf.Name)
				}
				lg(" ")
			}
		}

		// for k, v := range dwf.Dirs {
		// 	lg("%v %v", k, v)
		// }

	}

	return

}

func addDefaults(w http.ResponseWriter, r *http.Request, in FetchCommand) FetchCommand {

	lg, lge := loghttp.Logger(w, r)
	_, _ = lg, lge

	var preset FetchCommand

	h := in.Host
	if exactPreset, ok := ConfigDefaults[h]; ok {
		preset = exactPreset
	} else {
		preset = ConfigDefaults["unspecified"]
	}

	in.DepthTolerance = preset.DepthTolerance
	in.CondenseTrailingDirs = preset.CondenseTrailingDirs
	if in.DesiredNumber == 0 {
		in.DesiredNumber = preset.DesiredNumber
	}

	if in.RssXMLURI == nil || len(in.RssXMLURI) == 0 {
		in.RssXMLURI = preset.RssXMLURI
	}

	return in
}

// Each domain might have *several* RSS URLs.
// Function selectRSS returns the most fitting RSS URL
// for a given  SearchPrefix, or empty string.
func selectRSS(w http.ResponseWriter, r *http.Request, c FetchCommand) (ret string) {

	lg, lge := loghttp.Logger(w, r)
	_, _ = lg, lge

	cntr := 0
	sp := c.SearchPrefix

MarkX:
	for {

		// lg("search pref %v", sp)

		if rss, ok := c.RssXMLURI[sp]; ok {
			ret = rss
			lg("found rss url %v for %v", ret, sp)
			break MarkX

		}

		spPrev := sp
		sp = path.Dir(sp)
		if sp == "/" && spPrev == "/" ||
			sp == "." && spPrev == "." {
			lg("Did not find a RSS URL for %v, %q", c.SearchPrefix, ret)
			break
		}

		cntr++
		if cntr > 20 {
			lg("Select RSS Loop did not terminate. %v", c.SearchPrefix)
			break
		}
	}

	return
}

func attrX(attributes []html.Attribute, key string) (s string) {
	for _, a := range attributes {
		if key == a.Key {
			s = a.Val
			break
		}
	}
	return
}
