package oauthpb

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
/*

	JavaScript init:

{
  "widgetUrl": "https://tec-news.appspot.com/auth/authorized-redirect",
  "signInSuccessUrl": "https://tec-news.appspot.com/auth/signin-success",
  "signOutUrl": "https://tec-news.appspot.com/auth/signout",
  "oobActionUrl": "https://tec-news.appspot.com/auth/send-email",
  "apiKey": "AIzaSyAnarmnl8f0nHkGSqyU6CUdZxeN9e_5LhM",
  "siteName": "this site",
  "signInOptions": ["password","google","facebook"]
}


Server configuration

*/
var CodeBaseDirectory = "/not-initialized"

var gitkit_server_config_json = `{
  "clientId": "153437159745-cong6hlqenujf9o8fvl0gvum5gb9np1t.apps.googleusercontent.com",
  "serviceAccountEmail": "153437159745-c79ndj0k7csi118tj489v14jkm7iln1f@developer.gserviceaccount.com",
  "serviceAccountPrivateKeyFile": "[CodeBaseDirectory]/static/app-accessible/test.txt",
  "widgetUrl": "https://tec-news.appspot.com/auth/authorized-redirect",
  "cookieName": "gtoken"
}`

func init() {
	var err error
	CodeBaseDirectory, err = os.Getwd()
	if err != nil {
		panic("could not call the code base directory: " + err.Error() + "<br>\n")
	}
	gitkit_server_config_json = strings.Replace(gitkit_server_config_json, "[CodeBaseDirectory]", CodeBaseDirectory, -1)
}

func AuthIndex(w http.ResponseWriter, r *http.Request) {
	str := `
		<a href="/auth/authorized-redirect" >AuthorizedRedirect</a><br>
		<a href="/auth/signin-success" >SigninSuccess</a><br>
		<a href="/auth/signout" >Signout</a><br>
		<a href="/auth/send-email" >Send Email</a><br>
	`
	w.Write([]byte(str))

	bts, err := ioutil.ReadFile(CodeBaseDirectory + "/static/app-accessible/test.txt")
	if err != nil {
		w.Write([]byte(err.Error() + "<br>\n"))
	}
	w.Write(bts)

	w.Write([]byte("<pre>\n"))
	w.Write([]byte(gitkit_server_config_json))
	w.Write([]byte("</pre>\n"))
}

func AuthorizedRedirect(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("AuthorizedRedirect"))
}
func SigninSuccess(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("SigninSuccess"))
}
func Signout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Signout"))
}
func SendEmail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("SendEmail"))
}
