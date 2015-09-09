package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/pbberlin/tools/appengine/instance_mgt"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/net/http/tplx"
)

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

		CondenseTrailingDirs: 2,
		DepthTolerance:       1,
		DesiredNumber:        5,
	},
}

/*

[{ 	'Host':           'www.handelsblatt.com',
 	'SearchPrefix':   '/politik/international',
 	'RssXMLURI':      '/contentexport/feed/schlagzeilen',
}]



curl -X POST -d "[{ \"Host\": \"www.handelsblatt.com\",  \"SearchPrefix\":  \"/politik/deutschland\"         }]"  localhost:8085/fetch/command-receive
curl -X POST -d "[{ \"Host\": \"www.welt.de\"         ,  \"SearchPrefix\":  \"/wirtschaft/deutschland\"      }]"  localhost:8085/fetch/command-receive
curl -X POST -d "[{ \"Host\": \"www.economist.com\"   ,  \"SearchPrefix\":  \"/news/business-and-finance\"   }]"  localhost:8085/fetch/command-receive

curl -X POST -d "[{ \"Host\": \"test.economist.com\"  ,  \"SearchPrefix\":  \"/news/business-and-finance\"   }]"  localhost:8085/fetch/command-receive
curl -X POST -d "[{ \"Host\": \"test.economist.com\"  ,  \"SearchPrefix\":  \"/\"                            }]"  localhost:8085/fetch/command-receive

curl -X POST -d "[{ \"Host\": \"www.welt.de\",           \"SearchPrefix\": \"/wirtschaft/deutschland\" ,  \"RssXMLURI\": \"/wirtschaft/?service=Rss\" }]" localhost:8085/fetch/command-receive


curl localhost:8085/fetch/similar?uri-x=www.welt.de/politik/ausland/article146154432/Tuerkische-Bodentruppen-marschieren-im-Nordirak-ein.html

curl --data url-x=a.com  localhost:8085/fetch/similar
curl --data url-x=https://www.welt.de/politik/ausland/article146154432/Tuerkische-Bodentruppen-marschieren-im-Nordirak-ein.html  localhost:8085/fetch/similar
curl --data url-x=http://www.economist.com/news/britain/21663648-hard-times-hard-hats-making-britain-make-things-again-proving-difficult  localhost:8085/fetch/similar
curl --data url-x=http://www.economist.com/news/americas/21661804-gender-equality-good-economic-growth-girl-power  localhost:8085/fetch/similar


*/

// Submit test commands internally, without http request.
func staticFetchDirect(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	FetchHTML(w, r, testCommands)
}

// Submit test commands by http posting them.
func staticFetchViaPosting2Receiver(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {

	lg, lge := loghttp.Logger(w, r)

	wpf(w, tplx.ExecTplHelper(tplx.Head, map[string]string{"HtmlTitle": "JSON Post"}))
	defer wpf(w, tplx.Foot)

	wpf(w, "<pre>")
	defer wpf(w, "</pre>")

	b, err := Post2Receiver(r, testCommands)

	lge(err)
	lg("msg from Post2Receiver:")
	lg(b.String())

}

// Post2Receiver takes commands and http posts them to
// the command receiver
func Post2Receiver(r *http.Request, commands []FetchCommand) (*bytes.Buffer, error) {

	b := new(bytes.Buffer)

	if commands == nil || len(commands) == 0 {
		return b, fmt.Errorf("Slice of commands nil or empty %v", commands)
	}

	ii := instance_mgt.Get(r)
	fullURL := fmt.Sprintf("https://%s%s", ii.PureHostname, uriFetchCommandReceiver)
	wpf(b, "sending to URL:    %v\n", fullURL)

	bcommands, err := json.MarshalIndent(commands, "", "\t")
	if err != nil {
		wpf(b, "marshalling to []byte failed\n")
		return b, err
	}

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(bcommands))
	if err != nil {
		wpf(b, "creation of POST request failed\n")
		return b, err
	}
	req.Header.Set("X-Custom-Header-Counter", "nocounter")
	req.Header.Set("Content-Type", "application/json")

	bts, reqUrl, err := fetch.UrlGetter(r, fetch.Options{Req: req})
	_, _ = bts, reqUrl
	if err != nil {
		wpf(b, "Sending the POST request failed\n")
		return b, err
	}

	wpf(b, "effective req url: %v\n", reqUrl)
	wpf(b, "response body:\n")
	wpf(b, "%s\n", html.EscapeString(string(bts)))

	return b, nil
}
