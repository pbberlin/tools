package pbfetch

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"

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

	u, err := url.Parse(sUrl)
	if err != nil {
		return nil, fmt.Errorf("url unparseable: %v", err)
	}
	if httpsOnly || u.Scheme == "" {
		u.Scheme = "https"
	}

	host, port, err = net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}
	log.Println("get url:", u.String())
	log.Println("host and port: ", host, port, "standalone:", u.Host)

	resp, err := client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("get request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad http resp code: %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read resp body: %v", err)
	}

	return byteContent, nil

}
