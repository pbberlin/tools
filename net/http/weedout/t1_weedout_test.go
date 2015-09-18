// +build weed1
// go test -tags=weed1

package weedout

import (
	"path"
	"testing"

	"appengine/aetest"

	"github.com/pbberlin/tools/net/http/domclean2"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/fileserver"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/repo"
	"github.com/pbberlin/tools/stringspb"
)

func Test1(t *testing.T) {

	lg, b := loghttp.BuffLoggerUniversal(nil, nil)
	_ = b

	c, err := aetest.NewContext(nil)
	lg(err)
	if err != nil {
		return
	}
	defer c.Close()
	fs := GetFS(c, 2)

	remoteHostname := "www.welt.de"
	remoteHostname = "www.welt.de/politik/ausland"

	dirs1, _, msg, err := fileserver.GetDirContents(repo.RepoURL, remoteHostname)
	if err != nil {
		lg(err)
		lg("%s", msg)
	}

	lg("dirs1")
	for _, v := range dirs1 {
		lg("    %v", v)
	}

	least3URLs := []string{}
	for _, v1 := range dirs1 {

		p := path.Join(remoteHostname, v1)
		dirs2, fils2, msg, err := fileserver.GetDirContents(repo.RepoURL, p)
		_ = dirs2
		if err != nil {
			lg(err)
			lg("%s", msg)
		}
		// lg("  dirs2 %v", stringspb.IndentedDump(dirs2))
		// lg("  fils2 %v", stringspb.IndentedDump(fils2))

		for _, v2 := range fils2 {
			least3URLs = append(least3URLs, path.Join(remoteHostname, v1, v2))
		}
	}

	if len(least3URLs) < numTotal {
		lg("not enough files in rss fetcher cache")
		return
	} else {
		least3URLs = least3URLs[:numTotal+1]
	}

	lg("fils2")
	for _, v := range least3URLs {
		lg("    %v", v)
	}

	// domclean

	least3Files := make([]repo.FullArticle, 0, len(least3URLs))
	for i := 0; i < len(least3URLs); i++ {

		surl := spf("%v/%v", repo.RepoURL, least3URLs[i])

		fNamer := domclean2.FileNamer(logDir, i)
		fNamer() // first call yields key

		resBytes, inf, err := fetch.UrlGetter(nil, fetch.Options{URL: surl})
		if err != nil {
			lg(err)
			return
		}
		lg("fetched %4.1fkB from %v", float64(len(resBytes))/1024, stringspb.ToLenR(inf.URL.String(), 60))

		fa := repo.FullArticle{}
		fa.Url = inf.URL.String()
		fa.Mod = inf.Mod
		fa.Body = resBytes
		least3Files = append(least3Files, fa)

	}

	doc := WeedOut(least3Files, lg, fs)

	fNamer := domclean2.FileNamer(logDir, 0)
	fNamer() // first call yields key
	fsPerm := GetFS(c, 2)
	fileDump(lg, fsPerm, doc, fNamer, "_fin.html")

	pf("MapSimiliarCompares: %v SimpleCompares: %v LevenstheinComp: %v\n", breakMapsTooDistinct, appliedLevenshtein, appliedCompare)
	pf("correct finish\n")

}

// func weedoutFilename(articleId, weedoutStage int) (string, string) {
// 	stagedFn := fmt.Sprintf("outp_%03v_%v.html", articleId, weedoutStage)
// 	prefix := fmt.Sprintf("outp_%03v", articleId)
// 	return stagedFn, prefix
// }
