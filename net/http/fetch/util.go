package fetch

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/pbberlin/tools/net/http/loghttp"
)

func HostFromReq(r *http.Request) string {
	return splitPort(r.Host)
}

// A better url.Parse
func URLFromString(surl string) (*url.URL, error) {

	// Prevent "google.com" => u.Host == "" when scheme == ""
	if !strings.HasPrefix(surl, "http://") && !strings.HasPrefix(surl, "https://") {
		surl = "https://" + surl
	}
	ourl, err := url.Parse(surl)

	return ourl, err

}

func HostFromStringUrl(surl string) string {

	// Prevent "google.com" => u.Host == "" when scheme == ""
	if !strings.HasPrefix(surl, "http://") && !strings.HasPrefix(surl, "https://") {
		surl = "https://" + surl
	}
	url2, _ := url.Parse(surl)

	return splitPort(url2.Host)
}

func PathFromStringUrl(surl string) string {

	// Prevent "google.com" => u.Host == "" when scheme == ""
	if !strings.HasPrefix(surl, "http://") && !strings.HasPrefix(surl, "https://") {
		surl = "https://" + surl
	}
	url2, _ := url.Parse(surl)

	return splitPort(url2.Path)
}

func HostFromUrl(u *url.URL) string {

	// Prevent "google.com" => u.Host == "" when scheme == ""
	surl := u.String()
	return HostFromStringUrl(surl)
}

func splitPort(hp string) string {

	lg, b := loghttp.BuffLoggerUniversal(nil, nil)
	_ = b

	host, port, err := net.SplitHostPort(hp)
	_ = port

	if err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			// normal
		} else {
			lg(err)
		}
		host = hp
	}

	return host
}
