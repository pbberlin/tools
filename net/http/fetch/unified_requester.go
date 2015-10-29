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
	"github.com/pbberlin/tools/stringspb"
	"google.golang.org/appengine"

	oldAE "appengine"
	oldFetch "appengine/urlfetch"
)

var MsgNoRdirects = "redirect cancelled"
var ErrCancelRedirects = fmt.Errorf(MsgNoRdirects)
var ErrNoContext = fmt.Errorf("gaeReq did not yield a context; deadline exceeded?")

type Options struct {
	Req *http.Request

	URL string

	RedirectHandling int // 1 => call off upon redirects

	LogLevel int

	KnownProtocol                     string
	ForceHTTPSEvenOnDevelopmentServer bool
}

// Response info
type Info struct {
	URL *url.URL
	Mod time.Time
	Msg string
}

// UrlGetter universal http getter for app engine and standalone go programs.
// Previously response was returned. Forgot why. Dropped it.
func UrlGetter(gaeReq *http.Request, options Options) (
	[]byte, Info, error,
) {

	options.LogLevel = 2

	var err error
	var inf Info = Info{}

	if options.LogLevel > 0 {
		if options.Req != nil {
			inf.Msg += fmt.Sprintf("orig req url: %#v\n", options.Req.URL.String())
		} else {
			inf.Msg += fmt.Sprintf("orig str url: %#v\n", options.URL)
		}
	}

	//
	// Either take provided request
	// Or build one from options.URL
	if options.Req == nil {
		ourl, err := URLFromString(options.URL) // Normalize
		if err != nil {
			return nil, inf, err
		}
		options.URL = ourl.String()
		options.Req, err = http.NewRequest("GET", options.URL, nil)
		if err != nil {
			return nil, inf, err
		}
	} else {
		if options.Req.URL.Scheme == "" {
			options.Req.URL.Scheme = "https"
		}
	}
	r := options.Req

	if len(options.KnownProtocol) > 1 {
		if strings.HasSuffix(options.KnownProtocol, ":") {
			options.KnownProtocol = strings.TrimSuffix(options.KnownProtocol, ":")
		}
		if options.KnownProtocol == "http" || options.KnownProtocol == "https" {
			r.URL.Scheme = options.KnownProtocol
			inf.Msg += fmt.Sprintf("Using known protocol %q\n", options.KnownProtocol)
		}
	}

	//
	// Unifiy appengine plain http.client
	client := &http.Client{}
	if gaeReq == nil {
		client.Timeout = time.Duration(5 * time.Second) // GAE does not allow
	} else {
		c := util_appengine.SafelyExtractGaeContext(gaeReq)
		if c != nil {

			ctxOld := oldAE.NewContext(gaeReq)
			client = oldFetch.Client(ctxOld)

			// this does not prevent urlfetch: SSL_CERTIFICATE_ERROR
			// it merely leads to err = "DEADLINE_EXCEEDED"
			tr := oldFetch.Transport{Context: ctxOld, AllowInvalidServerCertificate: true}
			// thus
			tr = oldFetch.Transport{Context: ctxOld, AllowInvalidServerCertificate: false}

			tr.Deadline = 20 * time.Second // only possible on aeOld

			client.Transport = &tr
			// client.Timeout = 20 * time.Second // also not in google.golang.org/appengine/urlfetch

		} else {
			return nil, inf, ErrNoContext
		}

		// appengine dev server => always fallback to http
		if c != nil && appengine.IsDevAppServer() && !options.ForceHTTPSEvenOnDevelopmentServer {
			r.URL.Scheme = "http"
		}
	}

	inf.URL = r.URL

	if options.RedirectHandling == 1 {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {

			if len(via) == 1 && req.URL.Path == via[0].URL.Path+"/" {
				// allow redirect from /gesundheit to /gesundheit/
				return nil
			}

			spath := "\n"
			for _, v := range via {
				spath += v.URL.Path + "\n"
			}
			spath += req.URL.Path + "\n"
			return fmt.Errorf("%v %v", MsgNoRdirects, spath)
		}
	}

	if options.LogLevel > 0 {
		inf.Msg += fmt.Sprintf("url standardized to %q  %q %q \n", r.URL.Scheme, r.URL.Host, r.URL.RequestURI())
	}

	//
	//
	// Respond to test.economist.com directly from memory
	if _, ok := TestData[r.URL.Host+r.URL.Path]; ok {
		return TestData[r.URL.Host+r.URL.Path], inf, nil
	}

	// The actual call
	// =============================

	resp, err := client.Do(r)

	// Swallow redirect errors
	if err != nil {
		if options.RedirectHandling == 1 {
			serr := err.Error()
			if strings.Contains(serr, MsgNoRdirects) {
				bts := []byte(serr)
				inf.Mod = time.Now().Add(-10 * time.Minute)
				return bts, inf, nil
			}
		}
	}

	isHTTPSProblem := false
	if err != nil {
		isHTTPSProblem = strings.Contains(err.Error(), "SSL_CERTIFICATE_ERROR") ||
			strings.Contains(err.Error(), "tls: oversized record received with length")
	}

	// Under narrow conditions => fallback to http
	if err != nil {
		if isHTTPSProblem && r.URL.Scheme == "https" && r.Method == "GET" {
			r.URL.Scheme = "http"
			var err2nd error
			resp, err2nd = client.Do(r)
			// while protocol http may go through
			// next obstacle might be - again - a redirect error:
			if err2nd != nil {
				if options.RedirectHandling == 1 {
					serr := err2nd.Error()
					if strings.Contains(serr, MsgNoRdirects) {
						bts := []byte(serr)
						inf.Mod = time.Now().Add(-10 * time.Minute)
						addFallBackSuccessInfo(options, &inf, r, err)
						return bts, inf, nil
					}
				}

				return nil, inf, fmt.Errorf("GET fallback to http failed with %v", err2nd)
			}
			addFallBackSuccessInfo(options, &inf, r, err)
			err = nil // CLEAR error
		}
	}

	//
	// Final error handler
	//
	if err != nil {
		hintAE := ""
		if isHTTPSProblem && r.URL.Scheme == "https" {
			// Not GET but POST:
			// We cannot do a fallback for a post request - the r.Body.Reader is consumed
			// options.r.URL.Scheme = "http"
			// resp, err = client.Do(options.Req)
			return nil, inf, fmt.Errorf("Cannot do https requests. Possible reason: Dev server: %v", err)
		} else if strings.Contains(
			err.Error(),
			"net/http: Client Transport of type init.failingTransport doesn't support CancelRequest; Timeout not supported",
		) {
			hintAE = "\nDid you forget to submit the AE Request?\n"
		}
		return nil, inf, fmt.Errorf("request failed: %v - %v", err, hintAE)
	}

	//
	// We got response, but
	// explicit bad response from server
	if resp.StatusCode != http.StatusOK {

		if resp.StatusCode == http.StatusBadRequest || // 400
			resp.StatusCode == http.StatusNotFound || // 404
			false {
			dmp := ""
			for k, v := range resp.Header {
				dmp += fmt.Sprintf("key: %v - val %v\n", k, v)
			}
			dmp = ""
			dmp += stringspb.IndentedDump(r.URL)

			bts, errRd := ioutil.ReadAll(resp.Body)
			if errRd != nil {
				return nil, inf, fmt.Errorf("cannot read resp body: %v", errRd)
			}
			defer resp.Body.Close()

			err2 := fmt.Errorf("resp %v: %v \n%v \n<pre>%s</pre>", resp.StatusCode, r.URL.String(), dmp, bts)

			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
			var err2nd error
			resp, err2nd = client.Do(r)
			if err2nd != nil || resp.StatusCode != http.StatusOK {
				return nil, inf, fmt.Errorf("failed again %v %v \n%v", err2nd, resp.StatusCode, err2)
			}
			log.Printf("successful retry with '/' to %v after %v\n", r.URL.String(), err)
			err = nil // CLEAR error

			// return nil, inf, err2

		} else {
			return nil, inf, fmt.Errorf("bad http resp code: %v - %v", resp.StatusCode, r.URL.String())
		}
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, inf, fmt.Errorf("cannot read resp body: %v", err)
	}
	defer resp.Body.Close()

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
	inf.Mod = tlm
	// log.Printf("    hdr  %v %v\n", lm, tlm.Format(time.ANSIC))

	return bts, inf, nil

}

func addFallBackSuccessInfo(options Options, inf *Info, r *http.Request, err error) {
	if options.LogLevel > 0 {
		inf.Msg += fmt.Sprintf("\tsuccessful fallback to http %v", r.URL.String())
		inf.Msg += fmt.Sprintf("\tafter %v\n", err)
	}

}
