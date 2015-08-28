package fetch_rss

import (
	"log"
	"net/http"
	"path"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/runtimepb"
	"golang.org/x/net/html"
)

// Some websites have url paths like
// 	.../news/some--title/32168.html
// 	.../news/other-title/82316.html
//
// We want to condense these to
// 	.../news/some--title-32168.html
// 	.../news/other-title-82316.html
func condenseTrailingDir(uri string, n int) (ret string) {

	switch n {
	case 0:
		return uri
	case 1:
		return uri
	case 2:
		base1 := path.Base(uri)
		rdir1 := path.Dir(uri) // rightest Dir

		base2 := path.Base(rdir1)

		rdir2 := path.Dir(rdir1)

		ret = path.Join(rdir2, base2+"-"+base1)
	default:
		runtimepb.StackTrace(4)
		log.Fatalf("not implemented n > 2 (%v)", n)
	}

	return

}

// Adding preconfigured settings to a fetch command
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
// Function matchingRSSURI returns the most fitting RSS URL
// for a given  SearchPrefix, or empty string.
func matchingRSSURI(w http.ResponseWriter, r *http.Request, c FetchCommand) (ret string) {

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
			lg("Did not find a RSS URL for %v %q", c.SearchPrefix, ret)
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
