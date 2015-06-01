package util_appengine

/*
	separated from common utils, because non-app engine projects
	can not use the util package otherwise
*/

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"runtime"
	"time"

	"appengine"
)
import "os"
import "strings"

import "github.com/pbberlin/tools/util_err"

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

// Adapter() checks the path, takes the time, precomputes values into a map
// provides a global panic catcher
// The typed signature is cleaner than the long version:
//   func Adapter(given func(http.ResponseWriter, *http.Request, map[string]interface{})) http.HandlerFunc {
func Adapter(given util_err.ExtendedHandler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		c := appengine.NewContext(r)

		start := time.Now()
		defer requestStats(c, start)

		if !authenticate(w, r) {
			return
		}
		//check_against := r.URL.String()
		check_against := r.URL.Path
		matches := validRequestPath.FindStringSubmatch(check_against)
		if matches == nil {
			s := "illegal chars in path: " + check_against
			c.Infof(s)
			http.Error(w, s, http.StatusInternalServerError)
			return
		}

		s, err := url.Parse(r.URL.String())
		if err != nil {
			panic("Could not url.Parse current url")
		}
		dir := path.Dir(s.Path)
		base := path.Base(s.Path)

		map1 := map[string]interface{}{
			"dir":  dir,
			"base": base,
		}

		defer func() {
			// note: Println works even in panic
			//fmt.Println("--apapter(): going to catch panic on higher level")

			panicSignal := recover()
			if panicSignal != nil {
				miniStacktrace := ""
				for i := 1; i < 6; i++ {
					_, file, line, _ := runtime.Caller(i)
					if strings.Index(file, `/src/pkg/runtime/`) > -1 {
						miniStacktrace += fmt.Sprintf("<br>\n\t\t %s:%d ", file[len(file)-20:], line)
					} else {
						file = strings.TrimPrefix(file, util_err.AppDir)
						// we cannot determine, whether html is already sent
						// we cannot determine, whether we are in plaintext or html context
						// thus we need the <br>
						miniStacktrace += fmt.Sprintf("<br>\n\t\t /%s:%d ", file, line)
					}
				}
				// headers := w.Header()
				// for k, v := range headers {
				// 	miniStacktrace += fmt.Sprintf("%#v %#v<br>\n", k, v)
				// }
				if panicSignal == "abort_handler_processing" {
					s := fmt.Sprint("\thttp processing aborted\n", miniStacktrace)
					c.Errorf(s)
					w.Write([]byte(s))
				} else if panicSignal != nil {
					s := fmt.Sprintf("\tPANIC caught by util_err.Adapter: %v %s\n", panicSignal, miniStacktrace)
					c.Errorf(s)
					w.Write([]byte(s))
				}

			}
		}()

		given(w, r, map1)

	}
}

func requestStats(c appengine.Context, start time.Time) {
	age := time.Now().Sub(start)
	c.Infof("  request took %v", age)
}

/*

 A generic error function

 Utility functions pass   errors up to the caller


 Higher level "request functions" handle errors directly
    often we want to abort further request processing
 	and issue an message into the http response AND into the logs

 Sometimes we only want to write the error into the logs and
    continue operation => continueExecution true

  In addition to the generic error messages
  we may add specific error explanations or values
  via parameter vs - for display and logging
  We also show the source file+location.

  A "global panic catcher" in util_err.Adapter() ...defer(){}
  cooperates - suppressing the stacktrace and
  healing the panic

*/
func Err_http(w http.ResponseWriter, r *http.Request,
	bool_or_err interface{},
	continueExecution bool,
	vs ...interface{}) {

	var err error

	switch bool_or_err.(type) {
	default:
		type_unknown := fmt.Sprintf("%T", bool_or_err)
		err = errors.New("only bool or error - instead: -" + type_unknown + "-")
		panic(err)
	case nil:
		return
	case bool:
		if bool_or_err.(bool) {
			return
		}
		err = errors.New("Not OK (type conv?)")
	case error:
		err = bool_or_err.(error)
	}

	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		file = strings.TrimPrefix(file, util_err.AppDir)

		// we cannot determine, whether html is already sent
		// we cannot determine, whether we are in plaintext or html context
		// thus we need the <br>
		s := fmt.Sprintf("ERR: %v  <br>\n\t /%s:%d \n", err, file, line)
		if len(vs) > 0 {
			s = s + "\t" + fmt.Sprint(vs...) + "\n"
		}
		c := appengine.NewContext(r)

		if continueExecution {
			c.Infof(s)
		} else {
			c.Errorf(s)
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, s, http.StatusInternalServerError)
			panic("abort_handler_processing")
		}
	}

}
