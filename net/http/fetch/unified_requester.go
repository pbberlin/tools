// Package fetch performs a http request and returns the byte slice,
// also operating on google app engine.
package fetch

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pbberlin/tools/appengine/util_appengine"

	"appengine"
	"appengine/urlfetch"
)

var LogLevel = 0

var ErrNoRedirects = fmt.Errorf("redirect called off")
var sNoRedirects = ErrNoRedirects.Error()

type Options struct {
	URL string

	Req *http.Request

	HttpsOnly        bool
	RedirectHandling int // 1 => call off upon redirects
}

// Response info
type Info struct {
	URL *url.URL
	Mod time.Time
}

// UrlGetter universal http getter for app engine and standalone go programs.
// Previously response was returned. Forgot why. Dropped it.
func UrlGetter(gaeReq *http.Request, options Options) (
	[]byte, Info, error,
) {

	var err error
	var inf Info = Info{}

	if options.Req == nil {
		options.Req, err = http.NewRequest("GET", options.URL, nil)
		if err != nil {
			return nil, inf, err
		}
	}
	req := options.Req

	// Prevent u.Host from "google.com" without scheme is ""
	// Also make no protocol default to https
	surl := req.URL.String()
	if !strings.HasPrefix(surl, "http://") && !strings.HasPrefix(surl, "https://") {
		surl = "https://" + surl
	}
	req.URL, err = url.Parse(surl)
	if err != nil {
		return nil, inf, err
	}

	// Optional: force https
	if options.HttpsOnly {
		req.URL.Scheme = "https"
	}

	// Unifiy appengine plain http.client
	client := &http.Client{}
	if gaeReq == nil {
		client.Timeout = time.Duration(5 * time.Second) // GAE does not allow
	} else {
		c := util_appengine.SafelyExtractGaeContext(gaeReq)
		if c != nil {
			client = urlfetch.Client(c)

			// this does not prevent urlfetch: SSL_CERTIFICATE_ERROR
			// it merely leads to err = "DEADLINE_EXCEEDED"
			tr := urlfetch.Transport{Context: c, AllowInvalidServerCertificate: true}
			// thus
			tr = urlfetch.Transport{Context: c, AllowInvalidServerCertificate: false}
			client.Transport = &tr
		}

		// appengine dev server => always fallback to http
		if c != nil && appengine.IsDevAppServer() {
			req.URL.Scheme = "http"
		}
	}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return ErrNoRedirects
	}

	if LogLevel > 0 {
		log.Printf("host: %v, uri: %v \n", req.URL.Host, req.URL.RequestURI())
	}

	if _, ok := TestData[req.URL.Host+req.URL.Path]; ok {
		return TestData[req.URL.Host+req.URL.Path], Info{URL: req.URL}, nil
	}

	resp, err := client.Do(req)

	// Swallow redirect errors
	if err != nil {
		if options.RedirectHandling == 1 {
			serr := err.Error()
			if strings.Contains(serr, sNoRedirects) {
				bts := []byte(serr)
				return bts, Info{URL: req.URL, Mod: time.Now().Add(-10 * time.Minute)}, nil
			}
		}
	}

	cond := false
	if err != nil {
		cond = strings.Contains(err.Error(), "SSL_CERTIFICATE_ERROR") ||
			strings.Contains(err.Error(), "tls: oversized record received with length")
	}

	// Under narrow conditions => fallback to http
	if err != nil {
		if req.URL.Scheme == "https" && cond && req.Method == "GET" {
			req.URL.Scheme = "http"
			var err2nd error
			resp, err2nd = client.Do(req)
			if err2nd != nil {
				return nil, Info{URL: req.URL}, fmt.Errorf("cannot do https requests on dev server, fallback to http: %v",
					err2nd)
			}
			err = nil // CLEAR error
		}
	}

	if err != nil {
		hint := ""
		if req.URL.Scheme == "https" && cond {
			// We cannot repeat a post request - the r.Body.Reader is consumed
			// options.Req.URL.Scheme = "http"
			// resp, err = client.Do(options.Req)
			return nil, Info{URL: req.URL}, fmt.Errorf("cannot do https requests on dev server: %v", err)
		} else if strings.Contains(err.Error(), "net/http: Client Transport of type init.failingTransport doesn't support CancelRequest; Timeout not supported") {
			hint = "\n\n Did you forget to submit the AE Request?\n"
		}
		return nil, Info{URL: req.URL}, fmt.Errorf("request failed: %v - %v", err, hint)
	}

	//
	if resp.StatusCode != http.StatusOK {
		return nil, Info{URL: req.URL}, fmt.Errorf("bad http resp code: %v - %v", resp.StatusCode, req.URL.String())
	}

	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, Info{URL: req.URL}, fmt.Errorf("cannot read resp body: %v", err)
	}

	// time stamp
	var tlm time.Time // time last modified
	lm := resp.Header.Get("Last-Modified")
	if lm != "" {
		tlm, err = time.Parse(time.RFC1123, lm) // Last-Modified: Sat, 29 Aug 2015 21:15:39 GMT
		if err != nil {
			tlm, err = time.Parse(time.RFC1123Z, lm) // with numeric time zone
			if err != nil {
				var zeroTime time.Time
				tlm = zeroTime
			}
		}
	}
	// log.Printf("    hdr  %v %v\n", lm, tlm.Format(time.ANSIC))

	return bts, Info{URL: req.URL, Mod: tlm}, nil

}
