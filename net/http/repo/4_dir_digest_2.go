package repo

import (
	"net/http"
	"sort"
	"strings"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/os/osutilpb"
)

func DiveToDeepestMatch(dirTree *DirTree, uriPrefixIncl string) (*DirTree, string) {

	var subtree *DirTree
	subtree = dirTree
	head, dir, remainder := "", "", uriPrefixIncl

	if uriPrefixIncl == "/" || uriPrefixIncl == "" || uriPrefixIncl == "." {
		// exception for root
		head = "" // not "/"
	} else {
		// recur deeper
		for {
			dir, remainder, _ = osutilpb.PathDirReverse(remainder)
			head += dir
			// lg("    %-10q %-10q %-10q - head - dir - remainder", head, dir, remainder)
			if newSubtr, ok := subtree.Dirs[dir]; ok {
				subtree = &newSubtr
				// lg("    recursion found  %-10v %-10v for %v (%v)", dir, subtree.Name, uriPrefixIncl, subtree.Name)
			} else {
				// lg("    recursion failed %-10v %-10v for %v (%v)", dir, subtree.Name, uriPrefixIncl, subtree.Name)

				// Calling off searching on this level
				// StuffStage() itsself can step up if it wants to
				//
				// *not* setting subtree = nil
				// would keep us one level higher than this level.
				subtree = nil

				break
			}

			if remainder == "" {
				break
			}
		}
	}

	return subtree, head
}

type LevelWiseDeeperOptions struct {
	Rump       string // for instance /blogs/buttonwood
	ExcludeDir string // for instance /blogs/buttonwood/2014

	MinDepthDiff int // Relative to the depth of rump!; 0 equals Rump-Depth; thus 2 is the first restricting setting; set to 4 to exclude /blogs/buttonwood/2015/08/article1
	MaxDepthDiff int // To include /blogs/buttonwood/2015/08/article1 from /blogs/buttonwood => set to 3

	CondenseTrailingDirs int // See FetchCommand - equal member
	MaxNumber            int //
}

func LevelWiseDeeper(w http.ResponseWriter, r *http.Request, dtree *DirTree, opt LevelWiseDeeperOptions) []FullArticle {

	lg, lge := loghttp.Logger(w, r)
	_ = lge

	depthRump := strings.Count(opt.Rump, "/")

	arts := []FullArticle{}

	var fc func(string, *DirTree, int)

	fc = func(rmp1 string, dr1 *DirTree, lvl int) {

		// lg("      lvl %2v %v", lvl, dr1.Name)
		keys := make([]string, 0, len(dr1.Dirs))
		for k, _ := range dr1.Dirs {
			keys = append(keys, k)
		}
		// We could sort by LastFound
		// but we rather seek most current
		// files *later*
		sort.Strings(keys) // for debugging clarity
		for _, key := range keys {
			dr2 := dr1.Dirs[key]
			rmp2 := rmp1 + dr2.Name

			// lg("        %v", rmp2)

			//
			// rmp2 a candidate?
			if len(arts) > opt.MaxNumber-1 {
				return
			}

			if !dr2.EndPoint {
				continue
			}

			semanticUri := condenseTrailingDir(rmp2, opt.CondenseTrailingDirs)
			depthUri := strings.Count(semanticUri, "/")
			if depthUri-depthRump <= opt.MaxDepthDiff &&
				depthUri-depthRump >= opt.MinDepthDiff {
			} else {
				continue // we could also "break"
			}

			if opt.ExcludeDir == rmp2 {
				lg("    exclude dir %v", opt.ExcludeDir)
				continue
			}

			// lg("        including %v", semanticUri)
			art := FullArticle{Url: rmp2}
			if dr2.SrcRSS {
				art.Mod = dr2.LastFound
			}
			arts = append(arts, art)

		}

		//
		// recurse horizontally
		for _, key := range keys {
			dr2 := dr1.Dirs[key]
			rmp2 := rmp1 + dr2.Name
			if len(dr2.Dirs) == 0 {
				// lg("      LevelWiseDeeper - no children")
				continue
			}
			fc(rmp2, &dr2, lvl+1)
		}

	}

	fc(opt.Rump, dtree, 0)

	return arts

}
