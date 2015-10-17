// package login wraps the appengine user inside cookie name SACSID,
// as opposed to gitkit user wrapped inside SESSIONID.
package login

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/pbberlin/tools/net/http/htmlfrag"

	"appengine"
	"appengine/user"
)

func UserFromReq(r *http.Request) (*user.User, string) {

	c := appengine.NewContext(r)
	uType := ""

	//
	//
	u := user.Current(c)
	if u == nil {
		uType += "OAuth1 fail "
	} else {
		uType += "OAuth1 succ "
	}
	uType += "\n"

	//
	//
	u2, err := user.CurrentOAuth(c, "")
	if err != nil {
		uType += fmt.Sprintf("OAuth2 fail: %v", err)
	}
	if u2 != nil {
		uType += "OAuth2 succ "
	}
	uType += "\n"

	if appengine.IsDevAppServer() {
		if u2.Email == "example@example.com" {
			uType += fmt.Sprintf("OAuth2 reset %q. ", u2.Email)
			u2 = nil
		}
		uType += "CurrentOAuth() always exists on DEV system."
	}

	// Replace
	if u == nil {
		u = u2
	}

	return u, uType
}

func CheckForNormalUser(r *http.Request) (bool, *user.User, string) {

	if appengine.IsDevAppServer() {
		return true, &user.User{Email: "dev@server.com", Admin: true, ID: "32168"}, "DevServer login granted"
	}

	u, msg := UserFromReq(r)

	if u == nil {
		msg = "google appengine oauth required - normal rights - no login found\n" + msg
		return false, nil, msg
	}

	return true, u, msg

}

func CheckForAdminUser(r *http.Request) (bool, *user.User, string) {

	if appengine.IsDevAppServer() {
		return true, &user.User{Email: "dev@server.com", Admin: true, ID: "32168"}, "DevServer login granted"
	}

	u, msg := UserFromReq(r)

	if u == nil {
		msg = "google appengine oauth required - admin rights - no login found\n" + msg
		return false, nil, msg
	}
	if u != nil && !u.Admin {
		msg = "google appengine oauth required - admin rights - login found without admin\n" + msg
		return false, nil, msg
	}

	// if u.ID != "108853175242330402880" && u.ID != "S-1-5-21-2175189548-897864986-1736798499-1000" {
	// }

	return true, u, msg

}

// Show status and show login/logut url
func login(w http.ResponseWriter, r *http.Request) {

	r.Header.Set("X-Custom-Header-Counter", "nocounter")
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusFound)

	c := appengine.NewContext(r)
	u, uType := UserFromReq(r)

	if u == nil {

		fmt.Fprintf(w, "%v<br>\n", uType)
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// forward to login
		// w.Header().Set("Location", url)

		fmt.Fprintf(w, "<a href='%v'>login</a><br>", url)

	} else {

		// this gets never executed on dev server
		fmt.Fprintf(w, "Hello, %v, %v, %v, %v!<br>\n", u, u.ID, u.Email, u.FederatedIdentity)
		fmt.Fprintf(w, "Login type <pre>%v</pre><br>\n", uType)
		url2, err := user.LogoutURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "<a href='%v'>logout</a><br>", url2)
		urlLogoutDocumented := "/_ah/login?action=logout"
		fmt.Fprintf(w, "<a href='%v'>%v</a><br>", urlLogoutDocumented, urlLogoutDocumented)
	}

	fmt.Fprintf(w, htmlfrag.CookieDump(r))

}

func InitHandlers() {
	http.HandleFunc("/appengine/login", login)
}

// userinterface rendered to HTML - not only the strings for title and url
func BackendUIRendered() *bytes.Buffer {
	var b1 = new(bytes.Buffer)
	htmlfrag.Wb(b1, "Login AppEngine", "/appengine/login", "opposite of gitkit login")
	return b1
}
