package oauthpb

import (
	"fmt"
	"net/http"

	"appengine"
	"appengine/user"
)

func login(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Hello, %v, %v, %v, %v!", u, u.ID, u.Email, u.FederatedIdentity)
}

func Auth(r *http.Request) (bool, *user.User, string) {

	msg := ""
	c := appengine.NewContext(r)

	u := user.Current(c)
	if appengine.IsDevAppServer() {
		return true, u, "Logon always shines on DEV system."
	}
	// var err error
	// u, err = user.Current()
	// if err != nil {
	// 	msg += "user.Current() returned error :" + err.Error()
	// 	return

	if u == nil {
		msg += "google oauth required"
		return false, nil, msg
	}
	if u.ID != "108853175242330402880" && u.ID != "S-1-5-21-2175189548-897864986-1736798499-1000" {
		msg += "you need to be me; not " + u.ID
		return false, u, msg
	}

	return true, u, msg

}

func init() {
	http.HandleFunc("/login", login)
}
