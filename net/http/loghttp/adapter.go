// Package loghttp helps logging and or printing to the http response, branching for GAE.
package loghttp

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/dsu/distributed_unancestored"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	aelog "google.golang.org/appengine/log"
)

// added @ for /_ah/mail
var validRequestPath = regexp.MustCompile(`^([a-zA-Z0-9\.\-\_\/@]*)$`)

// works like an interface - functions just have to fit in the signature
type ExtendedHandler func(http.ResponseWriter, *http.Request, map[string]interface{})

// mjibson.appengine handler
type AppengineHandler func(context.Context, http.ResponseWriter, *http.Request)

/*

	This an adapter for adding a map to each handlerFunc

	http://golang.org/doc/articles/wiki/

		1.)  requi(a1)
		2.)  given(a1,a2)
	=> 3.)  requi(a1) = adapter( given(a1,a2) )

	func adapter(	 given func( t1, t2)	){
		return func( a1 ) {						  // signature of requi
			a2 := something							// set second argument
			given( a1, a2)
		}
	}

	No chance for closure context variables.
	They can not flow into given(),
	   because given() is not anonymous.


	Adding a map to each handlerFunc
		Precomputing keys to the directory and the base of the URL
			/aaa/bbbb/ccc.html => aaa/bbb ; ccc.html
		Can be enhanced with more precomputed stuff

	Validation: Checking path for legal chars

	Introducing a "global defer" func
		Catching panics and writing shortened stacktrace

*/

var C context.Context

// Adapter() checks the path, takes the time, precomputes values into a map
// provides a global panic catcher
// The typed signature is cleaner than the long version:
//   func Adapter(given func(http.ResponseWriter, *http.Request, map[string]interface{})) http.HandlerFunc {
func Adapter(given ExtendedHandler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		c, _ := util_appengine.SafelyExtractGaeCtxError(r)
		lgi := log.Printf
		lge := log.Fatalf
		if c != nil {
			defer logServerTime(c, start)
			// lgi = c.Infof
			lgi = func(format string, v ...interface{}) {
				aelog.Infof(c, format, v...)
			}

			// lge = c.Errorf
			lge = func(format string, v ...interface{}) {
				aelog.Errorf(c, format, v...)
			}

			C = c
		}

		if !authenticate(w, r) {
			return
		}

		//check_against := r.URL.String()
		check_against := r.URL.Path
		matches := validRequestPath.FindStringSubmatch(check_against)
		if matches == nil {
			s := "illegal chars in path: " + check_against
			lgi(s)
			http.Error(w, s, http.StatusInternalServerError)
			return
		}

		s, err := url.Parse(r.URL.String())
		if err != nil {
			panic("Could not url.Parse current url")
		}
		mp := map[string]interface{}{
			"dir":  path.Dir(s.Path),
			"base": path.Base(s.Path),
		}

		defer func() {
			// note: Println works even in panic

			panicSignal := recover()
			if panicSignal != nil {
				miniStacktrace := ""
				for i := 1; i < 11; i++ {
					_, file, line, _ := runtime.Caller(i)
					if strings.Index(file, `/src/pkg/runtime/`) > -1 {
						miniStacktrace += fmt.Sprintf("<br>\n\t\t %s:%d ", file[len(file)-20:], line)
					} else {
						dir := filepath.Dir(file)
						dirLast := filepath.Base(dir)
						file = filepath.Join(dirLast, filepath.Base(file))

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
					lge(s)
					w.Write([]byte(s))
				} else if panicSignal != nil {
					s := fmt.Sprintf("\tPANIC caught by util_err.Adapter: %v %s\n", panicSignal, miniStacktrace)
					lge(s)
					w.Write([]byte(s))
				}

			}
		}()

		r.Header.Set("adapter_01", "a string set by adapter")

		if c == nil {
			given(w, r, mp)
		} else {
			var given1 AppengineHandler
			given1 = func(c context.Context, w http.ResponseWriter, r *http.Request) {

				given(w, r, mp)

				// automatically set on appengine live, but not on appengine dev
				if r.Header.Get("Content-Type") == "" {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
				}

				if r.Header.Get("X-Custom-Header-Counter") != "nocounter" {
					cntr := 0
					if true {
						// This seems to cause problems with new applications
						// possible because of missing indize
						distributed_unancestored.Increment(c, mp["dir"].(string)+mp["base"].(string))
						cntr, _ = distributed_unancestored.Count(c, mp["dir"].(string)+mp["base"].(string))
					}
					fmt.Fprintf(w, "<br>\n%v Views<br>\n", cntr)
				}
			}

			if true || appengine.IsDevAppServer() {
				given1(c, w, r)
			} else {
				// wrapped := appstats.NewHandler(given1) // mjibson
				// wrapped.ServeHTTP(w, r)
			}

		}

	}
}

func authenticate(w http.ResponseWriter, r *http.Request) bool {
	return true
}

func logServerTime(c context.Context, start time.Time) {
	age := time.Now().Sub(start)
	if age.Seconds() < 0.01 {
		aelog.Infof(c, "  request took %v nano secs", age.Nanoseconds())
	} else {
		aelog.Infof(c, "  request took %2.2v secs", age.Seconds())
	}
}
