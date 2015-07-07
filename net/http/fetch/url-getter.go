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

	"appengine/urlfetch"
)

var LogLevel = 0

// UrlGetter universal http getter for app engine and standalone go programs.
// Previously response was returned. Forgot why. Dropped it.
func UrlGetter(sUrl string, gaeReq *http.Request, httpsOnly bool) ([]byte, error) {

	client := &http.Client{}
	if gaeReq == nil {
		client.Timeout = time.Duration(5 * time.Second) // GAE does not allow
	} else {
		c, _ := util_appengine.SafeGaeCheck(gaeReq)
		if c != nil {
			client = urlfetch.Client(c)

			// this does not prevent urlfetch: SSL_CERTIFICATE_ERROR
			// it merely leads to err = "DEADLINE_EXCEEDED"
			tr := urlfetch.Transport{Context: c, AllowInvalidServerCertificate: true}
			// thus
			tr = urlfetch.Transport{Context: c, AllowInvalidServerCertificate: false}
			client.Transport = &tr
		}
	}

	if !strings.HasPrefix(sUrl, "http://") && !strings.HasPrefix(sUrl, "https://") {
		sUrl = "https://" + sUrl
	}

	u, err := url.Parse(sUrl)
	if err != nil {
		return nil, fmt.Errorf("url unparseable: %v", err)
	}

	if httpsOnly {
		u.Scheme = "https"
	}

	if LogLevel > 0 {
		log.Printf("host: %v, uri: %v \n", u.Host, u.RequestURI())
	}

	resp, err := client.Get(u.String())
	if err != nil {
		if strings.Contains(err.Error(), "SSL_CERTIFICATE_ERROR") && u.Scheme == "https" {
			u.Scheme = "http"
			resp, err = client.Get(u.String())
		}
	}

	if err != nil {
		return nil, fmt.Errorf("get request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad http resp code: %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read resp body: %v", err)
	}

	// log.Printf("len %v bytes\n", len(bts))

	return bts, nil

}
