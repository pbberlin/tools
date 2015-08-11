package fetch_rss

import (
	"net/url"
	"path"

	"github.com/pbberlin/tools/logif"
)

func fetchFileName(sUrl string, idx int) (orig, numbered string) {

	var err error

	u, err := url.Parse(sUrl)
	logif.E(err, "url unparseable: %v")

	uri := u.RequestURI()

	orig = path.Base(uri)
	orig = path.Join(docRoot, u.Host, orig)

	numbered = path.Join(docRoot, u.Host, spf("art%02v.html", idx))

	// pf("orig: %v | numbered: %v\n", orig, numbered)
	return orig, numbered

}
