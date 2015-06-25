// Package loghttp helps logging and or printing to the http response, branching for GAE.
package loghttp

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/pbberlin/tools/appengine/util_appengine"

	"appengine"
)

var validRequestPath = regexp.MustCompile(`^([a-zA-Z0-9\.\-\_\/]*)$`)

// works like an interface - functions just have to fit in the signature
type ExtendedHandler func(http.ResponseWriter, *http.Request, map[string]interface{})

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

// Adapter() checks the path, takes the time, precomputes values into a map
// provides a global panic catcher
// The typed signature is cleaner than the long version:
//   func Adapter(given func(http.ResponseWriter, *http.Request, map[string]interface{})) http.HandlerFunc {
func Adapter(given ExtendedHandler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		c, _ := util_appengine.SafeGaeCheck(r)
		if c != nil {
			defer logServerTime(c, start)
		}

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

func authenticate(w http.ResponseWriter, r *http.Request) bool {
	return true
}

func logServerTime(c appengine.Context, start time.Time) {
	age := time.Now().Sub(start)
	c.Infof("  request took %v", age)
}
