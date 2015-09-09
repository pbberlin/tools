package repo

import (
	"net/http"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/httpfs"
	"github.com/pbberlin/tools/os/fsi/memfs"
)

var docRoot = ""                                           // no relative path, 'cause working dir too flippant
var whichType = 0                                          // which type of filesystem, default is memfs
var memMapFileSys = memfs.New(memfs.DirSort("byDateDesc")) // package variable required as "persistence"

const uriSetType = "/fetch/set-fs-type"
const mountName = "mntftch"
const UriMountNameY = "/" + mountName + "/serve-file/"

const cTestHostDev = "localhost:8085"

var repoURL = cTestHostDev + UriMountNameY

const uriFetchCommandReceiver = "/fetch/command-receive"
const uriFetchCommandSender = "/fetch/command-send"
const UriFetchSimilar = "/fetch/similar"

//
var httpFSys = &httpfs.HttpFs{SourceFs: fsi.FileSystem(memMapFileSys)} // memMap is always ready
var fileserver1 = http.FileServer(httpFSys.Dir(docRoot))

var msg = []byte(`<p>This is an embedded static http server.</p>
<p>
It serves previously downloaded pages<br>
 i.e. from handelsblatt or economist.
</p>`)

// ConfigDefaults are default values for FetchCommands
var ConfigDefaults = map[string]FetchCommand{
	"unspecified": FetchCommand{
		RssXMLURI:            map[string]string{},
		CondenseTrailingDirs: 0,
		DepthTolerance:       1,
		DesiredNumber:        5,
	},
	"www.handelsblatt.com": FetchCommand{
		RssXMLURI: map[string]string{
			"/":                      "/contentexport/feed/schlagzeilen",
			"/politik":               "/contentexport/feed/schlagzeilen",
			"/politik/international": "/contentexport/feed/schlagzeilen",
			"/politik/deutschland":   "/contentexport/feed/schlagzeilen",
		},
		CondenseTrailingDirs: 2,
		DepthTolerance:       1,
		DesiredNumber:        5,
	},
	"www.economist.com": FetchCommand{
		RssXMLURI: map[string]string{
			"/news/europe":               "/sections/europe/rss.xml",
			"/news/business-and-finance": "/sections/business-finance/rss.xml",
		},
		CondenseTrailingDirs: 0,
		DepthTolerance:       2,
		DesiredNumber:        5,
	},
	"test.economist.com": FetchCommand{
		RssXMLURI: map[string]string{
			"/news/business-and-finance": "/sections/business-finance/rss.xml",
		},
		CondenseTrailingDirs: 0,
		DepthTolerance:       2,
		DesiredNumber:        5,
	},
	"www.welt.de": FetchCommand{
		RssXMLURI: map[string]string{
			"/wirtschaft/deutschland":   "/wirtschaft/?service=Rss",
			"/wirtschaft/international": "/wirtschaft/?service=Rss",
		},

		// CondenseTrailingDirs: 2,
		CondenseTrailingDirs: 0,
		DepthTolerance:       1,
		DesiredNumber:        5,
	},
}
