// +build similar
// go test -tags=similar

package repo

import (
	"path"
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

	treePath := "/news"
	treePath = "/news/europe"
	treePath = "/news"
	treePath = "/blogs"
	treePath = "/blogs/graphicdetail"

	opt := LevelWiseDeeperOptions{}
	opt.Rump = treePath
	opt.ExcludeDir = "/news/americas"
	opt.ExcludeDir = "/blogs/buttonwood"
	opt.MinDepthDiff = 3
	opt.MaxDepthDiff = 5
	opt.CondenseTrailingDirs = 0
	opt.MaxNumber = 16

	var subtree *DirTree

	for i := 0; i < 3; i++ {

		subtree, treePath = DiveToDeepestMatch(dirTree, treePath)

		if subtree == nil {
			lg("#%v treePath %v ; subtree is nil", i, treePath)
		} else {
			lg("#%v treePath %v ; subtree exists", i, treePath)

			articles := LevelWiseDeeper(nil, nil, subtree, opt)
			for _, art := range articles {
				lg("#%v fnd %v", i, art.Url)
			}

			if len(articles) >= opt.MaxNumber {
				break
			}

			pathPrev := treePath
			treePath = path.Dir(treePath)
			lg("#%v  bef %v - aft %v", i, pathPrev, treePath)

			if pathPrev == "." && treePath == "." ||
				pathPrev == "/" && treePath == "/" {
				break
			}
		}

	}

	// lg("subtr 400 chars is %v \nend of subtr", stringspb.ToLen(subtree.String(), 400))

}
