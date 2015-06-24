package pbfetch

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"appengine"
	"appengine/urlfetch"
)

// UrlGetter universal http getter for app engine and standalone go programs.
// Previously response was returned. Forgot why. Dropped it.
func UrlGetter(sUrl string, gaeReq *http.Request, httpsOnly bool) ([]byte, error) {

	client := &http.Client{}
	if gaeReq != nil {
		c := appengine.NewContext(gaeReq)
		if c != nil {
			client = urlfetch.Client(c)
		}
	}
	client.Timeout = time.Duration(5 * time.Second)
	// client.Timeout = time.Duration(500 * time.Millisecond)

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

	// log.Println("host:   ", u.Host)
	// log.Println("get url:", u.RequestURI())

	resp, err := client.Get(u.String())
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

	return bts, nil

}
