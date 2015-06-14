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

func UrlGetter(sUrl string, gaeReq *http.Request, httpExclusively bool) (*http.Response, []byte, error) {

	client := &http.Client{}
	if gaeReq != nil {
		c := appengine.NewContext(gaeReq)
		if c != nil {
			client = urlfetch.Client(c)
		}
	}

	u, err := url.Parse(sUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("url unparseable: %v", err)
	}
	if !httpExclusively {
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
		return nil, nil, fmt.Errorf("get request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("bad http resp code: %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read resp body: %v", err)
	}

	return resp, byteContent, nil

}
