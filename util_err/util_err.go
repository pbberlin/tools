package util_err

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

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
		file = strings.TrimPrefix(file, AppDir)

		// we cannot determine, whether html is already sent
		// we cannot determine, whether we are in plaintext or html context
		// thus we need the <br>
		s := fmt.Sprintf("ERR: %v  <br>\n\t /%s:%d \n", err, file, line)
		if len(vs) > 0 {
			s = s + "\t" + fmt.Sprint(vs...) + "\n"
		}

		if continueExecution {
			log.Printf(s)
		} else {
			log.Printf(s)
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, s, http.StatusInternalServerError)
			panic("abort_handler_processing")
		}
	}

}

func LogAndShow(w http.ResponseWriter, r *http.Request, f string, vs ...interface{}) {

	log.Printf(f, vs...)

	var s string
	if len(vs) > 0 {
		s = fmt.Sprintf(f, vs...) + "\n"
	}
	w.Write([]byte(s))
}

// log only
func Err_log(e error) {
	if e != nil {
		_, file, line, _ := runtime.Caller(1)
		s := fmt.Sprintf("ERR: %v  \n\tSRC: %s:%d \n", e, file, line)
		log.Printf(s)
	}
}

// log and panic
func Err_panic(e error) {
	if e != nil {
		_, file, line, _ := runtime.Caller(1)
		s := fmt.Sprintf("ERR: %v  \n\tSRC: %s:%d \n", e, file, line)
		log.Printf(s)
		panic(s)
	}
}

func StackTrace(max int) {
	lg := log.New(os.Stdout, "str", 0)

	for i := 1; i <= max; i++ {
		_, file, line, _ := runtime.Caller(i)
		lg.Printf("        %s:%d ", file, line)
	}
}

// SuppressPanicUponDoubleRegistration registers
// a request hanlder for a route.
//
//
// Because of asynchronicity we need to
// catch the ensuing panic for repeated registration
// of the same handler
func SuppressPanicUponDoubleRegistration(w http.ResponseWriter, r *http.Request,
	urlPattern string, handler func(http.ResponseWriter, *http.Request)) string {
	defer func() {
		panicSignal := recover()
		if panicSignal != nil {
			w.Write([]byte(fmt.Sprintf("panic caught:\n\n %s", panicSignal)))
		}
	}()

	http.HandleFunc(urlPattern, handler)
	return urlPattern

}
