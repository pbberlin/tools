package repo

import (
	"net/http"
	"path"
	"strings"

	"appengine"

	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/routes"
	"github.com/pbberlin/tools/stringspb"
)

// FetchSimilar is an extended version of Fetch
func FetchSimilar(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, lge := loghttp.Logger(w, r)

	fs1 := GetFS(appengine.NewContext(r))

	err := r.ParseForm()
	lge(err)

	surl := r.FormValue(routes.URLParamKey)
	ourl, err := fetch.URLFromString(surl)
	lge(err)
	if err != nil {
		return
	}
	if ourl.Host == "" {
		lg("host is empty")
		return
	}

	srcDepth := strings.Count(ourl.Path, "/")

	cmd := FetchCommand{}
	cmd.Host = ourl.Host
	cmd.SearchPrefix = ourl.Path
	cmd = addDefaults(w, r, cmd)

	dirTree := &DirTree{Name: "/", Dirs: map[string]DirTree{}, EndPoint: true}
	fnDigest := path.Join(docRoot, cmd.Host, "digest2.json")
	loadDigest(w, r, fs1, fnDigest, dirTree) // previous
	lg("dirtree 400 chars is %v \nend of dirtree", stringspb.ToLen(dirTree.String(), 400))

	err = fetchCrawlSave(w, r, dirTree, fs1, path.Join(cmd.Host, ourl.Path))
	lge(err)
	if err != nil {
		return
	}

	var treePath string
	treePath = "/news"
	treePath = "/blogs"
	treePath = "/blogs/freeexchange"
	treePath = "/news/europe"
	treePath = path.Dir(ourl.Path)

	opt := LevelWiseDeeperOptions{}
	opt.Rump = treePath
	opt.ExcludeDir = "/news/americas"
	opt.ExcludeDir = "/blogs/buttonwood"
	opt.ExcludeDir = "/something-impossible"
	opt.MinDepthDiff = 1
	opt.MaxDepthDiff = 1
	opt.CondenseTrailingDirs = cmd.CondenseTrailingDirs
	opt.MaxNumber = 2266 // cmd.DesiredNumber

	var subtree *DirTree
	articles := []FullArticle{}

Mark1:
	for j := 0; j < srcDepth; j++ {
		treePath = path.Dir(ourl.Path)
	Mark2:
		for i := 1; i < 22; i++ {

			subtree, treePath = DiveToDeepestMatch(dirTree, treePath)

			lg("\nLooking from height %v to level %v  - %v", srcDepth-i, srcDepth-j, treePath)

			err = fetchCrawlSave(w, r, dirTree, fs1, path.Join(cmd.Host, treePath))
			lge(err)
			if err != nil {
				return
			}

			if subtree == nil {
				lg("\n#%v treePath %q ; subtree is nil", i, treePath)
			} else {
				// lg("\n#%v treePath %q ; subtree exists", i, treePath)

				opt.Rump = treePath
				opt.MinDepthDiff = i - j
				opt.MaxDepthDiff = i - j
				lpArticles := LevelWiseDeeper(nil, nil, subtree, opt)
				articles = append(articles, lpArticles...)
				for _, art := range lpArticles {
					lg("#%v     fnd %v", i, stringspb.ToLen(art.Url, 100))
				}

				if len(articles) >= opt.MaxNumber {
					lg("found enough")
					break Mark1
				}

				pathPrev := treePath
				treePath = path.Dir(treePath)
				// lg("#%v  bef %v - aft %v", i, pathPrev, treePath)

				if pathPrev == "." && treePath == "." ||
					pathPrev == "/" && treePath == "/" ||
					pathPrev == "" && treePath == "." {
					lg("break to innner")
					break Mark2
				}
			}

		}
	}

	// lg("subtr 400 chars is %v \nend of subtr", stringspb.ToLen(subtree.String(), 400))

}
