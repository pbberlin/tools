package weedout

import (
	"fmt"
	"os"
	"path"

	"github.com/pbberlin/tools/net/http/fileserver"
	"github.com/pbberlin/tools/net/http/loghttp"
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

func similar(surl string) ([]string, error) {

	lg, lge := loghttp.Logger(nil, nil)

	if surl == "" {
		surl = "www.welt.de/wirtschaft/article145511716-Deutsche-Rentner-zieht-es-ins-Ausland.html"
	}
	head, base := path.Split(surl)

	dirs1, fils1, msg, err := fileserver.GetDirContents(repoURL, path.Dir(surl))
	if err != nil {
		lge(err)
		lg("%s", msg)
	}

	lg("dirs1")
	for _, v := range dirs1 {
		lg("    %v", v)
	}
	lg("fils1")
	for _, v := range fils1 {
		lg("    %v", v)
	}

	least3Files := []string{}
	if len(fils1) > numTotal-1 {
		for _, v1 := range fils1 {

			if v1 == base {
				continue
			}
			least3Files = append(least3Files, path.Join(head, v1))
			if len(least3Files) == numTotal-1 {
				break
			}
		}

	}

	if len(least3Files) < numTotal {
		return []string{}, fmt.Errorf("not enough files in rss fetcher cache")
	}

	return least3Files, nil

}
