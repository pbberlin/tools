package pboauth

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

func init() {
	http.HandleFunc("/login", login)
}
