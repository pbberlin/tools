// package gitkit1 omits gitkit functions for updating user data or account-deletion;
// instead all app-specific settings and the redirect-urls are customizable.
package gitkit1

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
	"net/http"
	"time"

	// issues certificates (tokens) for possible http requests, making other requests impossible

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

	WidgetSigninAuthorizedRedirectURL = "/auth/authorized-redirect" // THIS one needs to be registered all over
	signOutURL                        = "/auth/signout"
	accountChooserBrandingURL         = "/auth/accountChooserBranding.html"

	oobActionURL = "/auth/send-email" // not needed, but kept
)

var (
	signinLandingDefaultURL  = "/auth/signin-default-landing"
	signoutLandingDefaultURL = "/auth/signout-default-landing"
	// successLandingURL = "/auth/signin-landing"
	// signoutLandingURL = "/auth/signout-landing"
)

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

	maxTokenAge = 1200 // 20 minutes

	maxSessionIDAge = 1800
)

var (
	xsrfKey      string
	cookieStore  *sessions.CookieStore
	gitkitClient *gitkit.Client
	client       *gitkit.Client
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

// SignedIn is on top of LoggedIn.
// Can only be measured by presence of cookie gtokenCookieName
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
		// client, err := gitkit.NewWithContext(c, gitkitClient)
		// if err != nil {
		// 	aelog.Errorf(c, "Failed to create a gitkit.Client with a context: %s", err)
		// 	return nil
		// }
		// Extract the token string from request.
		ts := client.TokenFromRequest(r)
		if ts == "" {
			return nil
		}
		// Check the token issue time. Only accept token that is no more than 15
		// minutes old even if it's still valid.
		// token, err := client.ValidateToken(ts)
		token, err := client.ValidateToken(appengine.NewContext(r), ts, []string{clientID})
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
