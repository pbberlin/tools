package oauthpb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"appengine"

	"github.com/google/identity-toolkit-go-client/gitkit"
	"golang.org/x/net/context"
)

func init() {
	http.HandleFunc("/auth", AuthIndex)
	http.HandleFunc("/auth/authorized-redirect", AuthorizedRedirect)
	http.HandleFunc("/auth/signin-success", SigninSuccess)
	http.HandleFunc("/auth/signout", Signout)
	http.HandleFunc("/auth/send-email", SendEmail)
}

//
// https://developers.google.com/identity/choose-auth
// https://developers.google.com/identity/toolkit/web/configure-service
// https://developers.facebook.com/apps/942324259171809/dashboard

//
// Server configuration
var CodeBaseDirectory = "/not-initialized"

var client *gitkit.Client

// Provide configuration. gitkit.LoadConfig() can also be used to load
// the configuration from a JSON file.
var gitkit_server_config_json = `{
  "clientId": "153437159745-cong6hlqenujf9o8fvl0gvum5gb9np1t.apps.googleusercontent.com",
  "widgetUrl": "https://tec-news.appspot.com/auth/authorized-redirect",
  "serviceAccountPrivateKeyFile": "[CodeBaseDirectory]appaccess-only/tec-news-49bc2267287d.p12",
  "cookieName": "gtoken"
  "serviceAccountEmail": "153437159745-c79ndj0k7csi118tj489v14jkm7iln1f@developer.gserviceaccount.com",
}`

var config = &gitkit.Config{
	ClientID:     "153437159745-cong6hlqenujf9o8fvl0gvum5gb9np1t.apps.googleusercontent.com",
	WidgetURL:    "[protoc_host]/auth/authorized-redirect",
	ServerAPIKey: "AIzaSyCnFQTG9WlS-y-eDpv3GtCUQhsUy61q8B8",
}

func init() {
	var err error
	CodeBaseDirectory, err = os.Getwd()
	if err != nil {
		panic("could not call the code base directory: " + err.Error() + "<br>\n")
	}

	if appengine.IsDevAppServer() {
		config.WidgetURL = strings.Replace(config.WidgetURL, "[protoc_host]", "https://tec-news.appspot.com", -1)
	} else {
		config.WidgetURL = strings.Replace(config.WidgetURL, "[protoc_host]", "http://localhost:8087", -1)
	}

	// Service account and private key are not required in Google App Engine
	// Prod environment. GAE App Identity API is used to identify the app.
	if appengine.IsDevAppServer() {
		config.ServiceAccount = "153437159745-c79ndj0k7csi118tj489v14jkm7iln1f@developer.gserviceaccount.com"
		config.PEMKeyPath = "[CodeBaseDirectory]appaccess-only/tec-news-49bc2267287d.pem"

		CodeBaseDirectory = path.Clean(CodeBaseDirectory) // remove trailing slash
		if !strings.HasSuffix(CodeBaseDirectory, "/") {
			CodeBaseDirectory += "/"
		}
		config.PEMKeyPath = strings.Replace(config.PEMKeyPath, "[CodeBaseDirectory]", CodeBaseDirectory, -1)
	}

	client, err = gitkit.New(config)
	if err != nil {
		panic("could not instantiate gitkit client: " + err.Error() + "<br>\n")
	}

}

func AuthIndex(w http.ResponseWriter, r *http.Request) {
	str := `
		<a href="/auth/authorized-redirect" >AuthorizedRedirect</a><br>
		<a href="/auth/signin-success" >SigninSuccess</a><br>
		<a href="/auth/signout" >Signout</a><br>
		<a href="/auth/send-email" >Send Email</a><br>
	`
	w.Write([]byte(str))

	bts, err := ioutil.ReadFile(CodeBaseDirectory + "appaccess-only/test.txt")
	if err != nil {
		w.Write([]byte(err.Error() + "<br>\n"))
	}
	w.Write(bts)

	w.Write([]byte("<pre>\n"))
	w.Write([]byte(gitkit_server_config_json))
	w.Write([]byte("</pre>\n"))
	w.Write([]byte("<br>\n"))

	{
		bts, err := ioutil.ReadFile(CodeBaseDirectory + "appaccess-only/tec-news-49bc2267287d.p12")
		if err != nil {
			w.Write([]byte(err.Error() + "<br>\n"))
		}
		w.Write(bts)
	}

}

func AuthorizedRedirect(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("AuthorizedRedirect"))

	// JavaScript init:
	str := `{
  "widgetUrl": "https://tec-news.appspot.com/auth/authorized-redirect",
  "signInSuccessUrl": "https://tec-news.appspot.com/auth/signin-success",
  "signOutUrl": "https://tec-news.appspot.com/auth/signout",
  "oobActionUrl": "https://tec-news.appspot.com/auth/send-email",
  "apiKey": "AIzaSyAnarmnl8f0nHkGSqyU6CUdZxeN9e_5LhM",
  "siteName": "this site",
  "signInOptions": ["password","google","facebook"]
}`

	w.Write([]byte(str))

}
func SigninSuccess(w http.ResponseWriter, r *http.Request) {

	// If there is no valid session, check identity tookit ID token.

	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	// ctx2 := appengine.NewContext(r)
	clientWAEC, err := gitkit.NewWithContext(ctx, client) // client with AE Context

	clientWAEC = client

	if err != nil {
		w.Write([]byte("Can not get gitkit client with AE context: " + err.Error() + "<br>\n"))
		return
	}

	// If there is no valid session, check identity tookit ID token.
	ts := clientWAEC.TokenFromRequest(r)
	token, err := clientWAEC.ValidateToken(ts)
	if err != nil {
		w.Write([]byte("Not a valid token: " + err.Error() + "<br>\n"))
		return
	}

	// Token is validate and it contains the user account information
	// including user ID, email address, etc.
	// Issue your own session cookie to finish the sign in.

	tss := fmt.Sprintf("SigninSuccess: \n<br>%#v", token)

	w.Write([]byte(tss))

}
func Signout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Signout"))
}
func SendEmail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("SendEmail"))
}
