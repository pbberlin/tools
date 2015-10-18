package gitkit

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"appengine/mail"

	"github.com/adg/xsrftoken"
	"github.com/pbberlin/tools/net/http/htmlfrag" // issues certificates (tokens) for possible http requests, making other requests impossible

	"github.com/google/identity-toolkit-go-client/gitkit"

	gorillaContext "github.com/gorilla/context"

	"google.golang.org/appengine"
	aelog "google.golang.org/appengine/log"

	aeOrig "appengine"
)

const home3 = `{{if .User}}  
  <p>Tired of FavWeekday?</p>
  <form method="POST" action="{{.DeleteAccountURL}}">
    <input type="hidden" name="xsrftoken" value="{{.DeleteAccountXSRFToken}}">
    <button type="submit">delete account</button>
  </form>
{{end}}`

const (
	deleteAccountURL = "/auth/deleteAccount"
	oobActionURL     = "/auth/send-email"
)

func UNUSEDinit() {

	// The gorilla sessions use gorilla request context
	ClearHandler := func(fc http.HandlerFunc) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer gorillaContext.Clear(r)
			fc(w, r)
		})
	}

	http.Handle(deleteAccountURL, ClearHandler(handleDeleteAccount))
	http.Handle(oobActionURL, ClearHandler(handleOOBAction))
}

func UNUSEDhandleHome(w http.ResponseWriter, r *http.Request) {

	u := currentUser(r)
	var d time.Weekday
	if u != nil {
		d = weekdayForUser(r, u)
	}
	saveCurrentUser(r, w, u)
	var xf, xd string
	if u != nil {
		xf = xsrftoken.Generate(xsrfKey, u.ID, updateURL)
		xd = xsrftoken.Generate(xsrfKey, u.ID, deleteAccountURL)
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
		"DeleteAccountURL":       deleteAccountURL,
		"DeleteAccountXSRFToken": xd,
	})
}

func handleOOBAction(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// Create an identity toolkit client associated with the GAE context.
	client, err := gitkit.NewWithContext(c, gitkitClient)
	if err != nil {
		aelog.Errorf(c, "Failed to create a gitkit.Client with a context: %s", err)
		w.Write([]byte(gitkit.ErrorResponse(err)))
		return
	}
	resp, err := client.GenerateOOBCode(r)
	if err != nil {
		aelog.Errorf(c, "Failed to get an OOB code: %s", err)
		w.Write([]byte(gitkit.ErrorResponse(err)))
		return
	}
	msg := &mail.Message{
		Sender: "FavWeekday Support <support@favweekday.appspot.com>",
	}
	switch resp.Action {
	case gitkit.OOBActionResetPassword:
		msg.Subject = "Reset your FavWeekday account password"
		msg.HTMLBody = fmt.Sprintf(emailTemplateResetPassword, resp.Email, resp.OOBCodeURL.String())
		msg.To = []string{resp.Email}
	case gitkit.OOBActionChangeEmail:
		msg.Subject = "FavWeekday account email address change confirmation"
		msg.HTMLBody = fmt.Sprintf(emailTemplateChangeEmail, resp.Email, resp.NewEmail, resp.OOBCodeURL.String())
		msg.To = []string{resp.NewEmail}
	case gitkit.OOBActionVerifyEmail:
		msg.Subject = "FavWeekday account registration confirmation"
		msg.HTMLBody = fmt.Sprintf(emailTemplateVerifyEmail, resp.OOBCodeURL.String())
		msg.To = []string{resp.Email}
	}
	c2 := aeOrig.NewContext(r)
	if err := mail.Send(c2, msg); err != nil {
		aelog.Errorf(c, "Failed to send %s message to user %s: %s", resp.Action, resp.Email, err)
		w.Write([]byte(gitkit.ErrorResponse(err)))
		return
	}
	w.Write([]byte(gitkit.SuccessResponse()))
}

/*

Failed to delete user {ID:14423325142879445183 Email:peter.buchmann.68@gmail.com
Name:Peter Buchmann EmailVerified:true}:
googleapi: Error 400: INVALID_LOCAL_ID, invalid

Failed to delete 00880189686365773816


Failed to delete user {ID: }: googleapi: Error 400: INVALID_LOCAL_ID, invalid
*/
func handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var (
		client *gitkit.Client
		err    error
	)
	// Check if there is a signed in user.
	u := currentUser(r)
	if u == nil {
		aelog.Errorf(c, "No signed in user for updating")
		goto out
	}
	// Validate XSRF token first.
	if !xsrftoken.Valid(r.PostFormValue(xsrfTokenName), xsrfKey, u.ID, deleteAccountURL) {
		aelog.Errorf(c, "XSRF token validation failed")
		goto out
	}
	// Create an identity toolkit client associated with the GAE context.
	client, err = gitkit.NewWithContext(c, gitkitClient)
	if err != nil {
		aelog.Errorf(c, "Failed to create a gitkit.Client with a context: %s", err)
		goto out
	}
	// Delete account.
	err = client.DeleteUser(&gitkit.User{LocalID: u.ID})
	if err != nil {
		aelog.Errorf(c, "Failed to delete user %v %v: %s", u.ID, u.Email, err)
		goto out
	}
	// Account deletion succeeded.
	// Call sign out to clear session and identity toolkit token.
	aelog.Infof(c, "Account deletion succeeded")

	handleSignOut(w, r)
	return
out:
	http.Redirect(w, r, signinSuccessAndHomeURL, http.StatusFound)
}
