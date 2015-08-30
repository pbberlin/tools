// +build weed2
// go test -tags=weed2

package weedout

import (
	"testing"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/repo"
	"github.com/pbberlin/tools/stringspb"
)

func Test2(t *testing.T) {

	lg, lge := loghttp.Logger(nil, nil)

	logdir := prepareLogDir()
	_ = logdir

	dirTree, err := fetchDirTree(repoURL, "www.economist.com")
	_ = dirTree
	lge(err)
	if err != nil {
		return
	}
	lg("dirtree 400 chars is %v \nend of dirtree", stringspb.ToLen(dirTree.String(), 400))

	path2 := "/news"
	path2 = "/blogs"
	path2 = "/news/europe"
	path2 = "/news"
	subtree, streePath := repo.DiveToDeepestMatch(dirTree, path2)

	lg("subtreePath is %v", streePath)

	if subtree == nil {
		lg("subtree is nil")
	} else {
		articles := repo.LevelWiseDeeper(nil, nil, streePath, subtree, "/americas", 2, 0, 2)
		for _, art := range articles {
			lg("found %v", art.Url)
		}
	}

	// lg("subtr 400 chars is %v \nend of subtr", stringspb.ToLen(subtree.String(), 400))

}
