package repo

import (
	"net/http"

	"github.com/pbberlin/tools/net/http/routes"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/httpfs"
	"github.com/pbberlin/tools/os/fsi/memfs"
)

// parallel fetchers routines
const numWorkers = 3

var docRoot = ""  // no relative path, 'cause working dir too flippant
var whichType = 0 // which type of filesystem, default is dsfs

var memMapFileSys = memfs.New(memfs.DirSort("byDateDesc"))             // package variable required as "persistence"
var httpFSys = &httpfs.HttpFs{SourceFs: fsi.FileSystem(memMapFileSys)} // memMap is always ready
var fileserver1 = http.FileServer(httpFSys.Dir(docRoot))

const mountName = "mntftch"

const uriSetType = "/fetch/set-fs-type"
const UriMountNameY = "/" + mountName + "/serve-file/"

const uriFetchCommandReceiver = "/fetch/command-receive"
const uriFetchCommandSender = "/fetch/command-send"

var RepoURL = routes.AppHost01 + UriMountNameY

var msg = []byte(`<p>This is an embedded static http server.</p>
<p>
It serves previously downloaded pages<br>
 i.e. from handelsblatt or economist.
</p>`)

var testCommands = []FetchCommand{
	FetchCommand{
		Host:         "www.handelsblatt.com",
		SearchPrefix: "/politik/deutschland/aa/bb",
	},
	FetchCommand{
		Host:         "www.handelsblatt.com",
		SearchPrefix: "/politik/international/aa/bb",
	},
	FetchCommand{
		Host:         "www.economist.com",
		SearchPrefix: "/news/europe/aa",
	},
}

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
