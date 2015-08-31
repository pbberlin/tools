// +build similar
// go test -tags=similar

package repo

import (
	"testing"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/stringspb"
)

func Test2(t *testing.T) {

	lg, lge := loghttp.Logger(nil, nil)

	dirTree, err := FetchDigest(repoURL, "www.economist.com")
	_ = dirTree
	lge(err)
	if err != nil {
		return
	}
	lg("dirtree 400 chars is %v \nend of dirtree", stringspb.ToLen(dirTree.String(), 400))

	path2 := "/news"
	path2 = "/news/europe"
	path2 = "/news"
	path2 = "/blogs"
	path2 = "/blogs/graphicdetail"
	subtree, streePath := DiveToDeepestMatch(dirTree, path2)

	lg("subtreePath is %v", streePath)

	if subtree == nil {
		lg("subtree is nil")
	} else {

		opt := LevelWiseDeeperOptions{}
		opt.Rump = streePath
		opt.ExcludeDir = "/news/americas"
		opt.ExcludeDir = "/blogs/buttonwood"
		opt.MinDepthDiff = 3
		opt.MaxDepthDiff = 5
		opt.CondenseTrailingDirs = 0
		opt.MaxNumber = 16

		articles := LevelWiseDeeper(nil, nil, subtree, opt)
		for _, art := range articles {
			lg("found %v", art.Url)
		}
	}

	// lg("subtr 400 chars is %v \nend of subtr", stringspb.ToLen(subtree.String(), 400))

}
