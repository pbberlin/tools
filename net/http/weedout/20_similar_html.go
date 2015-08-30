package weedout

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/repo"
)

func prepareLogDir() string {

	lg, lge := loghttp.Logger(nil, nil)

	logdir := "outp"
	lg("logdir is %v ", logdir)

	// sweep previous
	rmPath := spf("./%v/", logdir)
	err := os.RemoveAll(rmPath)
	if err != nil {
		lge(err)
		os.Exit(1)
	}
	lg("removed %q", rmPath)

	// create anew
	err = os.Mkdir(logdir, 0755)
	if err != nil && !os.IsExist(err) {
		lge(err)
		os.Exit(1)
	}

	return logdir

}

func fetchDirTree(hostWithPrefix, domain string) (*repo.DirTree, error) {

	lg, lge := loghttp.Logger(nil, nil)
	_ = lg

	surl := path.Join(hostWithPrefix, domain, "digest2.json")
	bts, _, err := fetch.UrlGetter(nil, fetch.Options{URL: surl})
	lge(err)
	if err != nil {
		return nil, err
	}

	// lg("%s", bts)
	dirTree := &repo.DirTree{Name: "/", Dirs: map[string]repo.DirTree{}, EndPoint: true}

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

// func similar(surl string) ([]string, error) {

// 	lg, lge := loghttp.Logger(nil, nil)

// 	if surl == "" {
// 		surl = "www.welt.de/wirtschaft/article145511716-Deutsche-Rentner-zieht-es-ins-Ausland.html"
// 	}
// 	head, base := path.Split(surl)

// 	dirs1, fils1, msg, err := fileserver.GetDirContents(repoURL, path.Dir(surl))
// 	if err != nil {
// 		lge(err)
// 		lg("%s", msg)
// 	}

// 	lg("dirs1")
// 	for _, v := range dirs1 {
// 		lg("    %v", v)
// 	}
// 	lg("fils1")
// 	for _, v := range fils1 {
// 		lg("    %v", v)
// 	}

// 	least3Files := []string{}
// 	if len(fils1) > numTotal-1 {
// 		for _, v1 := range fils1 {

// 			if v1 == base {
// 				continue
// 			}
// 			least3Files = append(least3Files, path.Join(head, v1))
// 			if len(least3Files) == numTotal-1 {
// 				break
// 			}
// 		}

// 	}

// 	if len(least3Files) < numTotal {
// 		return []string{}, fmt.Errorf("not enough files in rss fetcher cache")
// 	}

// 	return least3Files, nil

// }
