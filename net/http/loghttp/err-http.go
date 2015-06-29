package loghttp

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/runtimepb"
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
func E(w http.ResponseWriter, r *http.Request,
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
		line, file := runtimepb.LineFileXUp(1)
		// we cannot determine, whether html is already sent
		// we cannot determine, whether we are in plaintext or html context
		// thus we need the <br>
		s := fmt.Sprintf("ERR: %v  <br>\n\t /%s:%d \n", err, file, line)
		if len(vs) > 0 {
			s = s + "\t" + fmt.Sprint(vs...) + "\n"
		}

		if continueExecution {
			c, _ := util_appengine.SafeGaeCheck(r)
			if c == nil {
				log.Printf(s)
			} else {
				c.Infof(s)
			}
		} else {
			c, _ := util_appengine.SafeGaeCheck(r)
			if c == nil {
				log.Printf(s)
			} else {
				c.Errorf(s)
			}
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, s, http.StatusInternalServerError)
			panic("abort_handler_processing")
		}
	}

}

func Pf(w http.ResponseWriter, r *http.Request, f string, vs ...interface{}) {

	// Prepare the string
	var s string
	if len(vs) > 0 {
		s = fmt.Sprintf(f, vs...)
	} else {
		s = f
	}

	// Write it to http response
	w.Write([]byte(s))
	w.Write([]byte{'\n'})

	// Write to log/gae-log
	// Adding src code info
	line, file := runtimepb.LineFileXUp(1)
	s = fmt.Sprintf("%v \t\t- %v:%v", s, file, line)

	// Log it
	c, _ := util_appengine.SafeGaeCheck(r)
	if c == nil {
		// log.SetFlags(0)
		log.Printf(s)
		// log.SetFlags(log.Lshortfile)
	} else {
		c.Infof(s)
	}

}
