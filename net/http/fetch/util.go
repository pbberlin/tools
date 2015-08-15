package fetch

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/pbberlin/tools/logif"
)

func HostFromReq(r *http.Request) string {
	return inner(r.Host)
}

func HostFromUrl(u *url.URL) string {

	// Prevent u.Host from "google.com" without scheme is ""
	surl := u.String()
	if !strings.HasPrefix(surl, "http://") && strings.HasPrefix(surl, "https://") {
		surl += "http://" + surl
	}
	url2, _ := url.Parse(surl)

	return inner(url2.Host)
}

func inner(hp string) string {

	host, port, err := net.SplitHostPort(hp)
	_ = port

	if err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			// normal
		} else {
			logif.E(err)
		}
		host = hp
	}

	return host
}
