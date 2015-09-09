package repo

import (
	"fmt"
	"time"
)

var pf = fmt.Printf
var pfRestore = fmt.Printf

var spf = fmt.Sprintf
var wpf = fmt.Fprintf

// parallel fetchers routines
const numWorkers = 3

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
