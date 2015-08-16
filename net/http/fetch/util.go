package fetch

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/pbberlin/tools/logif"
)

func HostFromReq(r *http.Request) string {
	return splitPort(r.Host)
}

func HostFromUrl(u *url.URL) string {

	// Prevent "google.com" => u.Host == "" when scheme == ""
	surl := u.String()
	if !strings.HasPrefix(surl, "http://") && !strings.HasPrefix(surl, "https://") {
		surl = "https://" + surl
	}
	url2, _ := url.Parse(surl)

	return splitPort(url2.Host)
}

func splitPort(hp string) string {

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
