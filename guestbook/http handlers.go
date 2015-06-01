package guestbook

import (
	"net/http"

	"github.com/pbberlin/tools/util_appengine"
)

func init() {

	http.HandleFunc("/guest-entry", util_appengine.Adapter(guestEntry))
	http.HandleFunc("/guest-save", util_appengine.Adapter(guestSave))
	http.HandleFunc("/guest-view", util_appengine.Adapter(guestView))

}
