package parse2

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/pbberlin/tools/pblog"
)

func fetchFileName(sUrl string, idx int) (orig, numbered string) {

	var err error

	u, err := url.Parse(sUrl)
	pblog.LogE(err, "url unparseable: %v")

	uri := u.RequestURI()

	orig = filepath.Base(uri)
	orig = filepath.Join(docRoot, u.Host, orig)

	numbered = filepath.Join(docRoot, u.Host, spf("art%02v.html", idx))

	// pf("orig: %v | numbered: %v\n", orig, numbered)
	return orig, numbered

}

func weedoutFilename(articleId, weedoutStage int) (string, string) {
	fn := fmt.Sprintf("outp_%03v_%v.html", articleId, weedoutStage)
	prefix := fmt.Sprintf("outp_%03v", articleId)
	return fn, prefix
}
