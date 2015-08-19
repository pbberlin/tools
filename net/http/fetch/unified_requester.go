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

type Options struct {
	URL string

	Req *http.Request

	HttpsOnly bool
}

// UrlGetter universal http getter for app engine and standalone go programs.
// Previously response was returned. Forgot why. Dropped it.
func UrlGetter(gaeReq *http.Request, options Options) (
	[]byte, *url.URL, error,
) {

	var err error

	if options.Req == nil {
		options.Req, err = http.NewRequest("GET", options.URL, nil)
		if err != nil {
			return nil, nil, err
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
		return nil, nil, err
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

	if LogLevel > 0 {
		log.Printf("host: %v, uri: %v \n", req.URL.Host, req.URL.RequestURI())
	}

	resp, err := client.Do(req)

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
				return nil, req.URL, fmt.Errorf("cannot do https requests on dev server, fallback to http: %v",
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
			return nil, req.URL, fmt.Errorf("cannot do https requests on dev server: %v", err)
		} else if strings.Contains(err.Error(), "net/http: Client Transport of type init.failingTransport doesn't support CancelRequest; Timeout not supported") {
			hint = "\n\n Did you forget to submit the AE Request?\n"
		}
		return nil, req.URL, fmt.Errorf("request failed: %v - %v", err, hint)
	}

	//
	if resp.StatusCode != http.StatusOK {
		return nil, req.URL, fmt.Errorf("bad http resp code: %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, req.URL, fmt.Errorf("cannot read resp body: %v", err)
	}

	// log.Printf("len %v bytes\n", len(bts))

	return bts, req.URL, nil

}
