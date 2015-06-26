package json

import (
	"encoding/json"
	"fmt"
	"net/http"

	"appengine"
)

const client_secret string = `{
    "installed": {
        "client_id":"347979071940-kuu2a99r5i8h9q334i1k9r0a182pdr1f.apps.googleusercontent.com",
        "client_secret":"AwR-DDoRIK3iG9ai-4KP7rJm",
        "redirect_uris":["urn:ietf:wg:oauth:2.0:oob","oob"],
        "auth_uri":"https://accounts.google.com/o/oauth2/auth",
        "token_uri":"https://accounts.google.com/o/oauth2/token",
        "client_email":"",
        "client_x509_cert_url":"",
        "auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs"
    }
}`

type Servable_As_HTTP_JSON struct {
	Title, Body string
}

func (h Servable_As_HTTP_JSON) ServeHTTP(
	w http.ResponseWriter, r *http.Request,
) {

	c := appengine.NewContext(r)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(h); nil != err {
		c.Errorf(`{"error":"%s"}`, err)
	} else {
		c.Infof("encoding succeeded\n")
	}
	c.Infof("log: %v %T\n", h, h)
}

func jsonDecode(w http.ResponseWriter, r *http.Request) {

	var mapOuter map[string]interface{}
	byt := []byte(client_secret)

	if err := json.Unmarshal(byt, &mapOuter); err != nil {
		panic(err)
	}

	mapInner := mapOuter["installed"].(map[string]interface{})
	for key, val := range mapInner {
		s := fmt.Sprintf("%v \t\t %v\n", key, val)
		fmt.Fprint(w, s)
	}

}

func init() {
	http.HandleFunc("/json-decode", (jsonDecode))
	http.Handle("/json-encode", (Servable_As_HTTP_JSON{Title: "unused", Body: client_secret}))
}
