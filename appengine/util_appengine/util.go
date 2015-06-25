package util_appengine

/*
	separated from common utils, because non-app engine projects
	can not use the util package otherwise
*/

import (
	"net/http"
	"regexp"
	"time"

	"appengine"
)
import "os"
import "strings"

// http://regex101.com/
var validRequestPath = regexp.MustCompile(`^([a-zA-Z0-9\.\-\_\/]*)$`)

// IsLocalEnviron tells us, if we are on the
// local development server, or on the google app engine cloud maschine
func IsLocalEnviron() bool {

	return appengine.IsDevAppServer()

	s := os.TempDir()
	s = strings.ToLower(s)
	if s[0:2] == "c:" || s[0:2] == "d:" {
		// we are on windoofs - we are NOT on GAE
		return true
	}
	return false

}

func authenticate(w http.ResponseWriter, r *http.Request) bool {
	return true
}

func requestStats(c appengine.Context, start time.Time) {
	age := time.Now().Sub(start)
	c.Infof("  request took %v", age)
}
