// package googlesignin offers a login completely based on client javascript;
// signin-signout being messaged accross devices;
// installation of apps can be triggered;
// it has no server side login comparable to appengine/login or appengine/login/gitkit.
package googlesignin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pbberlin/tools/appengine/login"
	"github.com/pbberlin/tools/net/http/fetch"
	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/stringspb"
	"github.com/pbberlin/tools/vendor/jwt-go"
)

func init() {
	http.HandleFunc("/tokensignin", TokenSignin)
}

//
// https://developers.google.com/identity/choose-auth
// https://developers.google.com/identity/sign-in/web/backend-auth
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

	fc1 := func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:

		log.Printf("algo header is %v\n", token.Header["alg"])
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return token.Header["kid"], nil
	}

	token, err := jwt.Parse(myToken, fc1)

	// No direct error comparison possible; since err is wrapped in another struct
	if err != nil && strings.Contains(err.Error(), jwt.ErrPEMMappingObsolete.Error()) {

		currentPEMsURL := "https://www.googleapis.com/oauth2/v1/certs"
		req, err := http.NewRequest("GET", currentPEMsURL, nil)
		if err != nil {
			lg("creation of pem request failed")
			return
		}
		req.Header.Set("Content-Type", "application/json")

		fo := fetch.Options{Req: req}
		fo.KnownProtocol = "https"
		fo.ForceHTTPSEvenOnDevelopmentServer = true
		bts, inf, err := fetch.UrlGetter(r, fo)
		lg(err)
		if err != nil {
			lg("tried to fetch %v, %v", currentPEMsURL, inf.URL)
			lg("msg %v", inf.Msg)
			return
		}
		if len(bts) > 200 {
			var data1 map[string]string
			err = json.Unmarshal(bts, &data1)
			lg(err)
			// lg(stringspb.IndentedDumpBytes(data1))
			// w.Write(stringspb.IndentedDumpBytes(data1))
			if len(data1) > 1 {
				lg("PEM mappings updated")
				jwt.MappingToPEM = data1
			} else {
				lg("PEM mapping response contained only %v records; bytes length %v", len(data1), len(bts))
			}
		}

	}

	token, err = jwt.Parse(myToken, fc1)

	if err != nil && strings.Contains(err.Error(), jwt.ErrInvalidKey.Error()) {
		w.Write([]byte("The submitted RSA Key was somehow unparseable. We still accept the token.\n"))
		/*
			https://developers.google.com/identity/sign-in/web/backend-auth
		*/
		err = nil
		token.Valid = true
	}

	if err != nil {
		w.Write([]byte("--- " + err.Error() + ".\n"))
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

		_, usr, msg1 := login.CheckForNormalUser(r)
		if usr != nil {
			w.Write([]byte("ID from SRV is " + usr.ID + "\n"))
		}
		w.Write([]byte(msg1 + "\n"))

	} else {
		w.Write([]byte("tokensignin; INVALID. \n"))
		w.Write([]byte(tokSize))
		w.Write([]byte(stringspb.ToLen(myToken, 30)))

		vrf := fmt.Sprintf("\nhttps://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%v \n", myToken)
		w.Write([]byte(vrf))
	}

}
