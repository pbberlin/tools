// package gitkit expends the google identity toolkit;
// wrapping a user inside cookie SESSIONID;
// as opposed to appengine login cookie SACSID.
package gitkit

// Code taken from
// https://github.com/googlesamples/identity-toolkit-go/tree/master/favweekday
//
// The complete concept is expained here:
// https://developers.google.com/identity/choose-auth
// https://developers.google.com/identity/toolkit/web/federated-login
//
// https://developers.google.com/identity/toolkit/web/configure-service
// https://developers.google.com/identity/toolkit/web/setup-frontend
//
//
// Remove apps:
// https://security.google.com/settings/security/permissions
// https://www.facebook.com/settings?tab=applications

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/adg/xsrftoken"
	"github.com/pbberlin/tools/net/http/tplx" // issues certificates (tokens) for possible http requests, making other requests impossible

	"github.com/google/identity-toolkit-go-client/gitkit"
	"github.com/gorilla/sessions"

	"google.golang.org/appengine"
	aelog "google.golang.org/appengine/log"
)

// Action URLs.
// These need to be updated
// https://console.developers.google.com/project/tec-news/apiui/credential
// https://console.developers.google.com/project/tec-news/apiui/apiview/identitytoolkit/identity_toolkit
// https://developers.facebook.com/apps/942324259171809/settings/advanced/

const (
	homeURL = "/auth"

	widgetSigninAuthorizedRedirectURL = "/auth/authorized-redirect" // THIS one needs to be registered all over
	successLandingURL_X               = "SHOULD BE DYNAMIC"
	signOutURL                        = "/auth/signout"
	updateURL                         = "/auth/update"
	accountChooserBrandingURL         = "/auth/accountChooserBranding.html"
)

var (
	successLandingURL = "/auth/signin-landing"
	signoutLandingURL = "/auth/signout-landing"
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

	xsrfTokenName       = "xsrftoken"
	fieldNameFavWeekDay = "favorite"

	maxTokenAge = 1200 // 20 minutes

	maxSessionIDAge = 1800
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

func IsSignedIn(r *http.Request) bool {

	signedIn := false
	cks := r.Cookies()
	for _, ck := range cks {
		if ck.Name == gtokenCookieName {
			signedIn = true
			break
		}
	}

	return signedIn

}

//
// CurrentUser extracts the user information stored in current session.
//
// If there is no existing session, identity toolkit token is checked.
// If the token is valid, a new session is created.
//
// If any error happens, nil is returned.
func CurrentUser(r *http.Request) *User {
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
		if time.Now().Sub(token.IssueAt) > maxTokenAge*time.Second {
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

func f_HANDLERS() {

}

func handleHome(w http.ResponseWriter, r *http.Request) {

	format := `
		<a href='%v?mode=select'>Signin with Redirect (Widget)</a><br> 
		<a href='%v'>Signin Success Landing</a><br> 
		<a href='%v'>Signout </a><br> 
		<a href='%v'>Signout Landing</a><br> 
		<a href='%v'>Update</a><br> 
		<a href='%v'>Branding for Account Chooser</a><br> 
	`

	str := fmt.Sprintf(format,
		widgetSigninAuthorizedRedirectURL,
		successLandingURL,
		signOutURL,
		signoutLandingURL,
		updateURL,
		accountChooserBrandingURL,
	)

	bstpl := tplx.TemplateFromHugoPage(w, r) // the jQuery irritates
	fmt.Fprintf(w, tplx.ExecTplHelper(bstpl, map[string]interface{}{
		"HtmlTitle":       "Google Identity Toolkit Overview",
		"HtmlDescription": "", // reminder
		"HtmlContent":     template.HTML(str),
	}))

}

func HandleWidget(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	// Extract the POST body if any.
	b, _ := ioutil.ReadAll(r.Body)
	body, _ := url.QueryUnescape(string(b))

	gitkitTemplate := getWidgetTpl(w, r)

	gitkitTemplate.Execute(w, map[string]interface{}{
		"BrowserAPIKey":    browserAPIKey,
		"SignInSuccessUrl": successLandingURL,
		"SignOutURL":       signOutURL,
		"OOBActionURL":     oobActionURL, // unnecessary, since we don't offer "home account", but kept
		"POSTBody":         body,
	})

}

func HandleSuccess(w http.ResponseWriter, r *http.Request) {

	u := CurrentUser(r)

	if ok := IsSignedIn(r); !ok {
		u = nil
	}

	if u == nil {
		http.Redirect(w, r, widgetSigninAuthorizedRedirectURL+"?mode=select&user=wasNil", http.StatusFound)
	}

	saveCurrentUser(r, w, u)
	var xf string
	if u != nil {
		xf = xsrftoken.Generate(xsrfKey, u.ID, updateURL)
	}

	//
	var d time.Weekday
	if u != nil {
		d = weekdayForUser(r, u)
	}

	homeTemplate := getHomeTpl(w, r)
	homeTemplate.Execute(w, map[string]interface{}{
		"WidgetURL":              widgetSigninAuthorizedRedirectURL,
		"SignOutURL":             signOutURL,
		"User":                   u,
		"WeekdayIndex":           d,
		"Weekdays":               weekdays,
		"UpdateWeekdayURL":       updateURL,
		"UpdateWeekdayXSRFToken": xf,
		// "CookieDump": template.HTML(htmlfrag.CookieDump(r)),
	})
}

func handleSignOut(w http.ResponseWriter, r *http.Request) {

	sess, _ := cookieStore.Get(r, sessionName)
	sess.Options = &sessions.Options{MaxAge: -1} // MaxAge<0 means delete session cookie.
	err := sess.Save(r, w)
	if err != nil {
		aelog.Errorf(appengine.NewContext(r), "Cannot save session: %s", err)
	}

	// NONE of this has any effect
	if false {

		w.Header().Del("Set-Cookie")

		// The above deletion does not remove SESSIONID cookie.
		// This also does not remove SESSIONID.
		eraser := &http.Cookie{Name: sessionName, MaxAge: -1, Value: "erased",
			Expires: time.Now().Add(-240 * time.Hour), HttpOnly: true}
		http.SetCookie(w, eraser)
		eraser.Name = "SESSIONID"
		http.SetCookie(w, eraser)

		//
		w.Header().Del("Set-Cookie")
		ck := `set-cookie: SESSIONID=doesnthelp; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; Max-Age=1800; HttpOnly`
		ck = `set-cookie: SESSIONID=doesnthelp; expires=Thu, 01 Jan 1970 00:00:00 GMT`
		w.Header().Add("Set-Cookie", ck)

	}

	// Also clear identity toolkit token.
	http.SetCookie(w, &http.Cookie{Name: gtokenCookieName, MaxAge: -1})

	// Redirect to home page for sign in again.
	http.Redirect(w, r, signoutLandingURL+"?logout=true", http.StatusFound)
	// w.Write([]byte("<a href='" + signoutLandingURL + "'>Home<a>"))

}

func handleSignoutLanding(w http.ResponseWriter, r *http.Request) {

	format := `
		Signed out<br>
		<a href='%v'>Home</a><br> 
	`

	str := fmt.Sprintf(format, homeURL)

	bstpl := tplx.TemplateFromHugoPage(w, r) // the jQuery irritates
	fmt.Fprintf(w, tplx.ExecTplHelper(bstpl, map[string]interface{}{
		"HtmlTitle":       "Google Identity Toolkit Overview",
		"HtmlDescription": "", // reminder
		"HtmlContent":     template.HTML(str),
	}))

}

func handleUpdate(w http.ResponseWriter, r *http.Request) {

	operationResult := "failure"
	outFunc := func() {
		http.Redirect(w, r, successLandingURL+"?update="+operationResult, http.StatusFound)
	}

	var (
		d   int
		day time.Weekday
		err error
	)

	// Generic
	c := appengine.NewContext(r)
	// Check if there is a signed in user.
	u := CurrentUser(r)
	if u == nil {
		aelog.Errorf(c, "No signed in user for updating")
		outFunc()
		goto out
	}
	// Validate XSRF token first.
	if !xsrftoken.Valid(r.PostFormValue(xsrfTokenName), xsrfKey, u.ID, updateURL) {
		aelog.Errorf(c, "XSRF token validation failed")
		outFunc()
		goto out
	}

	//
	// Specific
	// Extract the new favorite weekday.
	d, err = strconv.Atoi(r.PostFormValue(fieldNameFavWeekDay))
	if err != nil {
		aelog.Errorf(c, "Failed to extract new favoriate weekday: %s", err)
		outFunc()
		goto out
	}
	day = time.Weekday(d)
	if day < time.Sunday || day > time.Saturday {
		aelog.Errorf(c, "Got wrong value for favorite weekday: %d", d)
		outFunc()
		goto out
	}
	// Update the favorite weekday.
	updateWeekdayForUser(r, u, day)
	operationResult = "success"

out:
	outFunc()
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
