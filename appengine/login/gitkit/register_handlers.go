package gitkit

import (
	"bytes"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/google/identity-toolkit-go-client/gitkit"
	gorillaContext "github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/pbberlin/tools/net/http/htmlfrag"

	"appengine"
)

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

func InitHandlers() {

	initCodeBaseDir()

	// Register datatypes such that it can be saved in the session.
	gob.Register(SessionUserKey(0))
	gob.Register(&User{})

	// Initialize XSRF token key.
	xsrfKey = "My personal very secure XSRF token key"

	sessKey := []byte("secure-key-234002395432-wsasjasfsfsfsaa-234002395432-wsasjasfsfsfsaa-234002395432-wsasjasfsfsfsaa")

	// Create a session cookie store.
	cookieStore = sessions.NewCookieStore(
		sessKey[:64],
		sessKey[:32],
	)

	cookieStore.Options = &sessions.Options{
		MaxAge:   maxSessionIDAge, // Session valid for 30 Minutes.
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

	http.Handle(homeURL, ClearHandler(handleHome))

	http.Handle(widgetSigninAuthorizedRedirectURL, ClearHandler(handleWidget))
	http.Handle(successLandingURL, ClearHandler(handleSuccess))

	http.Handle(signOutURL, ClearHandler(handleSignOut))
	http.Handle(signoutLandingURL, ClearHandler(handleSignoutLanding))

	http.Handle(updateURL, ClearHandler(handleUpdate))

	http.HandleFunc(accountChooserBrandingURL, accountChooserBranding)
}

// userinterface rendered to HTML - not only the strings for title and url
func BackendUIRendered() *bytes.Buffer {
	var b1 = new(bytes.Buffer)

	htmlfrag.Wb(b1, "Login GitKit", homeURL, "opposite of appengine login")
	htmlfrag.Wb(b1, "Signin", widgetSigninAuthorizedRedirectURL+"?mode=select", "")
	htmlfrag.Wb(b1, "Success Landing", successLandingURL, "")
	htmlfrag.Wb(b1, "Signout", signOutURL, "")
	htmlfrag.Wb(b1, "Signout Landing", signoutLandingURL, "")
	return b1
}
