package oauthpb

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pbberlin/tools/vendor/jwt-go"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/stringspb"

	"appengine"
	"appengine/user"
)

func login(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	uType := ""

	u := user.Current(c)
	if u != nil {
		uType += "Normal "
	}

	u2, err := user.CurrentOAuth(c, "")
	if err != nil {
		uType += fmt.Sprintf("OAuth failed %v", err)
	}
	if u2 != nil {
		uType += "OAuth2 "
	}

	// Replace
	if u == nil {
		u = u2
	}

	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	} else {
		// this gets never executed on dev server
		fmt.Fprintf(w, "Hello, %v, %v, %v, %v!<br>\n", u, u.ID, u.Email, u.FederatedIdentity)
		fmt.Fprintf(w, "Login type %v<br>\n", uType)
		url2, err := user.LogoutURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "<a href='%v'>logout</a>", url2)
	}

}

func Auth(r *http.Request) (bool, *user.User, string) {

	msg := ""
	c := appengine.NewContext(r)

	u := user.Current(c)
	u2, err := user.CurrentOAuth(c, "")
	if err != nil {
		msg += fmt.Sprintf("oauth user err %v", err)
	}
	if u == nil {
		u = u2
	}

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

	log.Printf("is admin: %v", u.Admin)

	// if u.ID != "108853175242330402880" && u.ID != "S-1-5-21-2175189548-897864986-1736798499-1000" {
	// 	msg += "you need to be me; not " + u.ID
	// 	return false, u, msg
	// }

	return true, u, msg

}

func myLookupKey(k interface{}) (interface{}, error) {
	return k, nil
}

//
//	https://developers.google.com/identity/choose-auth
//  https://developers.google.com/identity/sign-in/web/backend-auth
func TokenSignin(w http.ResponseWriter, r *http.Request) {

	lg, _ := loghttp.BuffLoggerUniversal(w, r)

	// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:1313")

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8087")

	w.Header().Del("Access-Control-Allow-Origin")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// err := r.ParseMultipartForm(1024 * 1024 * 2)
	err := r.ParseForm()
	lg(err)

	myToken := r.Form.Get("idtoken")
	tokSize := fmt.Sprintf("Len of Tok was %v. \n", len(myToken))

	// we can also verify the token from here:
	// https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=XYZ123

	fc1 := func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:

		log.Printf("algo header is %v\n", token.Header["alg"])
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return myLookupKey(token.Header["kid"])
	}

	token, err := jwt.Parse(myToken, fc1)

	if err != nil && strings.Contains(err.Error(), jwt.ErrInvalidKey.Error()) {
		w.Write([]byte("The submitted RSA Key was somehow unparseable. We still accept the token.\n"))
		/*
			https://developers.google.com/identity/sign-in/web/backend-auth
		*/
		err = nil
		token.Valid = true
	}

	if err != nil {
		w.Write([]byte(err.Error() + ".\n"))
	}

	if err == nil && token.Valid {

		tk := ""
		tk += fmt.Sprintf("     Algor:     %v\n", token.Method)
		tk += fmt.Sprintf("     Header:    %v\n", token.Header)
		for k, v := range token.Claims {
			tk += fmt.Sprintf("\t  %-8v %v\n", k, v)
		}
		lg(tk)

		w.Write([]byte("tokensignin; valid.   \n"))
		w.Write([]byte(tokSize))
		sb := "header-sub-not-present"
		if _, ok := token.Claims["sub"]; ok {
			sb = token.Claims["sub"].(string)
		}
		w.Write([]byte("ID from PWT is " + sb + "\n"))

		_, usr, msg1 := Auth(r)
		if usr != nil {
			w.Write([]byte("ID from SRV is " + usr.ID + "\n"))
		}
		w.Write([]byte(msg1 + "\n"))

	} else {
		w.Write([]byte("tokensignin; INVALID. \n"))
		w.Write([]byte(tokSize))
		w.Write([]byte(stringspb.ToLen(myToken, 30)))
	}

}

func init() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/tokensignin", TokenSignin)
}
