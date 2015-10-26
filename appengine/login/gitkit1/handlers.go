package gitkit1

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/pbberlin/tools/net/http/tplx"

	"google.golang.org/appengine"
	aelog "google.golang.org/appengine/log"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	// HandleHomeVari(w, r, successLandingURL, signoutLandingURL)
	HandleHomeVari(w, r, signinLandingDefaultURL, signoutLandingDefaultURL)
}

func HandleHomeVari(w http.ResponseWriter, r *http.Request, successLandingURL, signoutLandingURL string) {

	format := `
		<a href='%v?mode=select'>Signin with Redirect (Widget)</a><br><br> 

		<a href='%v'>Signin Success Landing</a><br><br> 
		<a href='%v'>Signout </a><br><br>

		<a href='%v'>Signout Landing</a><br> 
		<a href='%v'>Branding for Account Chooser</a><br> 
	`

	str := fmt.Sprintf(format,
		WidgetSigninAuthorizedRedirectURL,
		successLandingURL,
		signOutURL,

		signoutLandingURL,
		accountChooserBrandingURL,
	)

	bstpl := tplx.TemplateFromHugoPage(w, r) // the jQuery irritates
	fmt.Fprintf(w, tplx.ExecTplHelper(bstpl, map[string]interface{}{
		"HtmlTitle":       "Google Identity Toolkit Overview",
		"HtmlDescription": "", // reminder
		"HtmlContent":     template.HTML(str),
	}))

}

func handleWidget(w http.ResponseWriter, r *http.Request) {

	// param "red"  for redirect
	// must correspond with the widget param queryParameterForSignInSuccessUrl
	// It comes back with the encoded answer params - "source" and "state"
	// But it's only
	red := signinLandingDefaultURL

	if r.FormValue("red") != "" {
		red = r.FormValue("red") // This is not even necessary
		red = signinLandingDefaultURL
	}

	// log.Printf("\n-----------------------------------")
	// for key, v := range r.Form {
	// 	log.Printf("%10v is %#v", key, v)
	// }

	HandleVariWidget(w, r, red)
}

func HandleVariWidget(w http.ResponseWriter, r *http.Request, successLandingURL string) {

	defer r.Body.Close()
	// Extract the POST body if any.
	b, _ := ioutil.ReadAll(r.Body)
	body, _ := url.QueryUnescape(string(b))

	gitkitTemplate := GetWidgetTpl(w, r, siteName+" Identity Toolkit")

	gitkitTemplate.Execute(w, map[string]interface{}{
		"BrandingURL":         getConfig(siteName, "protocDomain") + accountChooserBrandingURL,
		"FaviconURL":          getConfig(siteName, "protocDomain") + "/favicon.ico",
		"BrowserAPIKey":       getConfig(siteName, "browserAPIKey"),
		"SignInSuccessURLTRY": template.URL(`/harcoded/ use queryParameterForSignInSuccessUrl=red insteads`),
		"SignInSuccessURL":    template.URL(successLandingURL), // widget
		"SignOutURL":          signOutURL,
		"OOBActionURL":        oobActionURL, // unnecessary, since we don't offer "home account", but kept
		"SiteName":            siteName,
		"POSTBody":            body,
	})

}

func handleSigninSuccessLanding(w http.ResponseWriter, r *http.Request) {
	HandleVariSuccess(w, r,
		siteName+" member home",
		UserInfoHTML+"<br><br>"+IDCardHTML+"<br><br>",
	)
}

func HandleVariSuccess(w http.ResponseWriter, r *http.Request, title, body string) {

	u := CurrentUser(r)

	if ok := IsSignedIn(r); !ok {
		u = nil
	}

	if u == nil {
		http.Redirect(w, r, WidgetSigninAuthorizedRedirectURL+"?mode=select&user=wasNil", http.StatusFound)
	}

	saveCurrentUser(r, w, u)

	//
	homeTemplate := GetHomeTpl(w, r, title, body)

	homeTemplate.Execute(w, map[string]interface{}{
		"WidgetURL":  WidgetSigninAuthorizedRedirectURL,
		"SignOutURL": signOutURL,
		"User":       u,
		// "CookieDump": template.HTML(htmlfrag.CookieDump(r)),
	})
}

func handleSignOut(w http.ResponseWriter, r *http.Request) {
	HandleVariSignOut(w, r, signoutLandingDefaultURL)
}

func HandleVariSignOut(w http.ResponseWriter, r *http.Request, signoutLandingURL string) {

	sess, _ := cookieStore.Get(r, sessionName)
	sess.Options = &sessions.Options{MaxAge: -1} // MaxAge<0 means delete session cookie.
	err := sess.Save(r, w)
	if err != nil {
		aelog.Errorf(appengine.NewContext(r), "Cannot save session: %s", err)
	}
	// Impossible to delete SESSIONID cookie

	// Also clear identity toolkit token.
	http.SetCookie(w, &http.Cookie{Name: gtokenCookieName, MaxAge: -1})

	// Redirect to home page for sign in again.
	http.Redirect(w, r, signoutLandingURL+"?logout=true", http.StatusFound)
	// w.Write([]byte("<a href='" + signoutLandingURL + "'>Home<a>"))

}

func handleSignOutLanding(w http.ResponseWriter, r *http.Request) {

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
      <img src="%v" 
      	style="display:block;height:120px;margin:auto">
      <p style="font-size:14px;opacity:.54;margin-top:20px;text-align:center">
        %v.
      </p>
    </div>
  </body>
</html>`

	str = fmt.Sprintf(str, getConfig(siteName, "accountChooserImg"), getConfig(siteName, "accountChooserHeadline"))

	w.Write([]byte(str))

}
