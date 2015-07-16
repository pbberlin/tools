package fetch

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/pbberlin/tools/logif"
)

func HostFromReq(r *http.Request) string { return inner(r.Host) }

func HostFromUrl(u *url.URL) string { return inner(u.Host) }

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
