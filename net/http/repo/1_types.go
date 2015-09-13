package repo

import "time"

// FullArticle is the main struct passed
// between the pipeline stages
type FullArticle struct {
	Url  string
	Mod  time.Time
	Body []byte
}

// FetchCommand contains a RSS location
// and details which items we want to fetch from it.
type FetchCommand struct {
	Host         string // www.handelsblatt.com,
	SearchPrefix string // /politik/international/aa/bb,

	RssXMLURI            map[string]string // SearchPrefix => RSS-URLs
	DesiredNumber        int
	CondenseTrailingDirs int // The last one or two directories might be article titles or ids
	DepthTolerance       int
}

type DirTree struct {
	Name      string // Name == key of Parent.Dirs
	LastFound time.Time

	SrcRSS   bool // dir came first from RSS
	EndPoint bool // dir came first from RSS - and is at the end of the path

	Dirs map[string]DirTree
	// Fils []string
}

type LevelWiseDeeperOptions struct {
	Rump       string // for instance /blogs/buttonwood
	ExcludeDir string // for instance /blogs/buttonwood/2014

	MinDepthDiff int // Relative to the depth of rump!; 0 equals Rump-Depth; thus 2 is the first restricting setting; set to 4 to exclude /blogs/buttonwood/2015/08/article1
	MaxDepthDiff int // To include /blogs/buttonwood/2015/08/article1 from /blogs/buttonwood => set to 3

	CondenseTrailingDirs int // See FetchCommand - equal member
	MaxNumber            int //
}
