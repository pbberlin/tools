// package gitkit expends the google identity toolkit;
// wrapping a user inside cookie SESSIONID;
// as opposed to appengine login cookie SACSID.
package gitkit

// Taken from
// https://github.com/googlesamples/identity-toolkit-go/tree/master/favweekday
//
// The complete concept is expained here:
// https://developers.google.com/identity/toolkit/web/federated-login
// https://developers.google.com/identity/choose-auth
//
// https://developers.google.com/identity/toolkit/web/configure-service
// https://developers.google.com/identity/toolkit/web/setup-frontend
//
//

import (
	"encoding/gob"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/adg/xsrftoken"
	"github.com/pbberlin/tools/net/http/htmlfrag" // issues certificates (tokens) for possible http requests, making other requests impossible

	"github.com/google/identity-toolkit-go-client/gitkit"
	gorillaContext "github.com/gorilla/context"
	"github.com/gorilla/sessions"

	"google.golang.org/appengine"
	aelog "google.golang.org/appengine/log"
)

// Action URLs.
// These need to be updated
// in https://console.developers.google.com/project/tec-news/apiui/credential
// in https://console.developers.google.com/project/tec-news/apiui/apiview/identitytoolkit/identity_toolkit
// and on facebook
const (
	homeAndSigninSuccessURL           = "/auth"
	widgetSigninAuthorizedRedirectURL = "/auth/authorized-redirect"
	signOutURL                        = "/auth/signout"
	oobActionURL                      = "/auth/send-email"
	updateURL                         = "/auth/update"
	accountChooserBrandingURL         = "/auth/accountChooserBranding.html"
)

// Identity toolkit configurations.
const (
	serverAPIKey  = "AIzaSyCnFQTG9WlS-y-eDpv3GtCUQhsUy61q8B8"
	browserAPIKey = "AIzaSyAnarmnl8f0nHkGSqyU6CUdZxeN9e_5LhM"

	clientID       = "153437159745-cong6hlqenujf9o8fvl0gvum5gb9np1t.apps.googleusercontent.com"
	serviceAccount = "153437159745-c79ndj0k7csi118tj489v14jkm7iln1f@developer.gserviceaccount.com"
)

// The pseudo absolute path to the pem keyfile
var CodeBaseDirectory = "/not-initialized"
var privateKeyPath = "[CodeBaseDirectory]appaccess-only/tec-news-49bc2267287d.pem"

// Cookie/Form input names.
const (

	// contains jws from appengine/user.CurrentUser() ...;
	// not used here
	aeUserSessName = "SACSID"

	// = cookie name;
	// contains jwt from google/facebook/twitter;
	// remains, even when "signed out"
	// remains, even when logging out of google/twitter
	// cannot be overwritten by "eraser"
	sessionName = "SESSIONID"

	// Created on top of sessionName on "signin"
	// Remains
	gtokenCookieName = "gtoken"

	xsrfTokenName = "xsrftoken"
	favoriteName  = "favorite"

	maxAgeSessionAndToken = 1800
)

var (
	xsrfKey      string
	cookieStore  *sessions.CookieStore
	gitkitClient *gitkit.Client
)

// User information.
type User struct {
	ID            string
	Email         string
	Name          string
	EmailVerified bool
}

// Key used to store the user information in the current session.
type SessionUserKey int

const sessionUserKey SessionUserKey = 0

//
// currentUser extracts the user information stored in current session.
//
// If there is no existing session, identity toolkit token is checked.
// If the token is valid, a new session is created.
//
// If any error happens, nil is returned.
func currentUser(r *http.Request) *User {
	c := appengine.NewContext(r)
	sess, _ := cookieStore.Get(r, sessionName)
	if sess.IsNew {
		// Create an identity toolkit client associated with the GAE context.
		client, err := gitkit.NewWithContext(c, gitkitClient)
		if err != nil {
			aelog.Errorf(c, "Failed to create a gitkit.Client with a context: %s", err)
			return nil
		}
		// Extract the token string from request.
		ts := client.TokenFromRequest(r)
		if ts == "" {
			return nil
		}
		// Check the token issue time. Only accept token that is no more than 15
		// minitues old even if it's still valid.
		token, err := client.ValidateToken(ts)
		if err != nil {
			aelog.Errorf(c, "Invalid token %s: %s", ts, err)
			return nil
		}
		if time.Now().Sub(token.IssueAt) > maxAgeSessionAndToken*time.Second {
			aelog.Infof(c, "Token %s is too old. Issused at: %s", ts, token.IssueAt)
			return nil
		}
		// Fetch user info.
		u, err := client.UserByLocalID(token.LocalID)
		if err != nil {
			aelog.Errorf(c, "Failed to fetch user info for %s[%s]: %s", token.Email, token.LocalID, err)
			return nil
		}
		return &User{
			ID:            u.LocalID,
			Email:         u.Email,
			Name:          u.DisplayName,
			EmailVerified: u.EmailVerified,
		}
	} else {
		// Extracts user from current session.
		v, ok := sess.Values[sessionUserKey]
		if !ok {
			aelog.Errorf(c, "no user found in current session")
		}
		return v.(*User)
	}
}

// saveCurrentUser stores the user information in current session.
func saveCurrentUser(r *http.Request, w http.ResponseWriter, u *User) {
	if u == nil {
		return
	}
	sess, _ := cookieStore.Get(r, sessionName)
	sess.Values[sessionUserKey] = *u
	err := sess.Save(r, w)
	if err != nil {
		aelog.Errorf(appengine.NewContext(r), "Cannot save session: %s", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {

	u := currentUser(r)
	var d time.Weekday
	if u != nil {
		d = weekdayForUser(r, u)
	}
	saveCurrentUser(r, w, u)
	var xf string
	if u != nil {
		xf = xsrftoken.Generate(xsrfKey, u.ID, updateURL)
	}

	homeTemplate := getHomeTpl(w, r)
	homeTemplate.Execute(w, map[string]interface{}{
		"CookieDump":             template.HTML(htmlfrag.CookieDump(r)),
		"WidgetURL":              widgetSigninAuthorizedRedirectURL,
		"SignOutURL":             signOutURL,
		"User":                   u,
		"WeekdayIndex":           d,
		"Weekdays":               weekdays,
		"UpdateWeekdayURL":       updateURL,
		"UpdateWeekdayXSRFToken": xf,
	})
}

func handleWidget(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Extract the POST body if any.
	b, _ := ioutil.ReadAll(r.Body)
	body, _ := url.QueryUnescape(string(b))

	gitkitTemplate := getWidgetTpl(w, r)

	// gitkitTemplate.Execute(
	// 	w, struct {
	// 		BrowserAPIKey    string
	// 		SignInSuccessUrl string
	// 		SignOutURL       string
	// 		OOBActionURL     string
	// 		POSTBody         string
	// 	}{browserAPIKey, homeAndSigninSuccessURL, signOutURL, oobActionURL,
	// 		body,
	// 	})
	gitkitTemplate.Execute(w, map[string]interface{}{
		"BrowserAPIKey":    browserAPIKey,
		"SignInSuccessUrl": homeAndSigninSuccessURL,
		"SignOutURL":       signOutURL,
		"OOBActionURL":     oobActionURL,
		"POSTBody":         body,
	})

}

func handleSignOut(w http.ResponseWriter, r *http.Request) {
	sess, _ := cookieStore.Get(r, sessionName)
	sess.Options = &sessions.Options{
		MaxAge: -1, // MaxAge<0 means delete session cookie.
	}
	err := sess.Save(r, w)
	if err != nil {
		aelog.Errorf(appengine.NewContext(r), "Cannot save session: %s", err)
	}

	if false {
		// The above deletion does not remove SESSIONID cookie.
		// This also does not remove SESSIONID.
		eraser := &http.Cookie{Name: sessionName, MaxAge: -1}
		eraser.Value = "erased"
		http.SetCookie(w, eraser)
	}

	// Also clear identity toolkit token.
	http.SetCookie(w, &http.Cookie{Name: gtokenCookieName, MaxAge: -1})
	// Redirect to home page for sign in again.
	http.Redirect(w, r, homeAndSigninSuccessURL, http.StatusFound)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var (
		d   int
		day time.Weekday
		err error
	)
	// Check if there is a signed in user.
	u := currentUser(r)
	if u == nil {
		aelog.Errorf(c, "No signed in user for updating")
		goto out
	}
	// Validate XSRF token first.
	if !xsrftoken.Valid(r.PostFormValue(xsrfTokenName), xsrfKey, u.ID, updateURL) {
		aelog.Errorf(c, "XSRF token validation failed")
		goto out
	}
	// Extract the new favorite weekday.
	d, err = strconv.Atoi(r.PostFormValue(favoriteName))
	if err != nil {
		aelog.Errorf(c, "Failed to extract new favoriate weekday: %s", err)
		goto out
	}
	day = time.Weekday(d)
	if day < time.Sunday || day > time.Saturday {
		aelog.Errorf(c, "Got wrong value for favorite weekday: %d", d)
	}
	// Update the favorite weekday.
	updateWeekdayForUser(r, u, day)
out:
	// Redirect to home page to show the update result.
	http.Redirect(w, r, homeAndSigninSuccessURL, http.StatusFound)
}

// Is called by AccountChooser to retrieve some layout.
// Dynamic execution required because of Access-Control header ...
func accountChooserBranding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	str := `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  </head>
  <body>
    <div style="width:256px;margin:auto">
      <img src="/img/house-of-cards-mousepointer-03-04.gif" 
      	style="display:block;height:120px;margin:auto">
      <p style="font-size:14px;opacity:.54;margin-top:20px;text-align:center">
        Welcome to tec-news insights.
      </p>
    </div>
  </body>
</html>`

	w.Write([]byte(str))

}

func initCodeBaseDir() {
	var err error
	CodeBaseDirectory, err = os.Getwd()
	if err != nil {
		panic("could not call the code base directory: " + err.Error() + "<br>\n")
	}
	// Make the path working
	CodeBaseDirectory = path.Clean(CodeBaseDirectory) // remove trailing slash
	if !strings.HasSuffix(CodeBaseDirectory, "/") {
		CodeBaseDirectory += "/"
	}
	privateKeyPath = strings.Replace(privateKeyPath, "[CodeBaseDirectory]", CodeBaseDirectory, -1)

}

func init() {

	initCodeBaseDir()

	// Register datatypes such that it can be saved in the session.
	gob.Register(SessionUserKey(0))
	gob.Register(&User{})

	// Initialize XSRF token key.
	xsrfKey = "My very secure XSRF token key"

	// Create a session cookie store.
	cookieStore = sessions.NewCookieStore(
		[]byte("My very secure authentication key for cookie store or generate one using securecookies.GenerateRamdonKey()")[:64],
		[]byte("My very secure encryption key for cookie store or generate one using securecookies.GenerateRamdonKey()")[:32])

	cookieStore.Options = &sessions.Options{
		MaxAge:   maxAgeSessionAndToken, // Session valid for two hours.
		HttpOnly: true,
	}

	// Create identity toolkit client.
	c := &gitkit.Config{
		ServerAPIKey: serverAPIKey,
		ClientID:     clientID,
		WidgetURL:    widgetSigninAuthorizedRedirectURL,
	}
	// Service account and private key are not required in GAE Prod.
	// GAE App Identity API is used to identify the app.
	if appengine.IsDevAppServer() {
		c.ServiceAccount = serviceAccount
		c.PEMKeyPath = privateKeyPath
	}
	var err error
	gitkitClient, err = gitkit.New(c)
	if err != nil {
		log.Fatal(err)
	}

	// The gorilla sessions use gorilla request context
	ClearHandler := func(fc http.HandlerFunc) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer gorillaContext.Clear(r)
			fc(w, r)
		})
	}

	// http.Handle(homeAndSigninSuccessURL, r)
	http.Handle(homeAndSigninSuccessURL, ClearHandler(handleHome))
	http.Handle(widgetSigninAuthorizedRedirectURL, ClearHandler(handleWidget))
	http.Handle(signOutURL, ClearHandler(handleSignOut))
	http.Handle(oobActionURL, ClearHandler(handleOOBAction))
	http.Handle(updateURL, ClearHandler(handleUpdate))
	http.HandleFunc(accountChooserBrandingURL, accountChooserBranding)
}
