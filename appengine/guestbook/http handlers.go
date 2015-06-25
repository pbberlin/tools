package guestbook

import (
	"net/http"

	"github.com/pbberlin/tools/net/http/loghttp"
)

func init() {
	http.HandleFunc("/guest-entry", loghttp.Adapter(guestEntry))
	http.HandleFunc("/guest-save", loghttp.Adapter(guestSave))
	http.HandleFunc("/guest-view", loghttp.Adapter(guestView))
}
