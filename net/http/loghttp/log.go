package loghttp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pbberlin/tools/appengine/util_appengine"
	"github.com/pbberlin/tools/runtimepb"
	"github.com/pbberlin/tools/stringspb"
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
			c, _ := util_appengine.SafelyExtractGaeCtxError(r)
			if c == nil {
				log.Printf(s)
			} else {
				c.Infof(s)
			}
		} else {
			c, _ := util_appengine.SafelyExtractGaeCtxError(r)
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

func Pf(w io.Writer, r *http.Request, f string, vs ...interface{}) {

	// Prepare the string
	var s string
	if len(vs) > 0 {
		s = fmt.Sprintf(f, vs...)
	} else {
		s = f
	}

	if s == "" {
		return
	}

	// Write it to http response - unless prefix 'lo ' - log only
	if w != nil && !strings.HasPrefix(s, "lo ") {
		w.Write([]byte(s))
		w.Write([]byte{'\n'})
	}

	// Write to log/gae-log
	// Adding src code info

	line, file := runtimepb.LineFileXUp(1)
	// if strings.HasSuffix(file, "log.go")
	if strings.HasSuffix(file, runtimepb.ThisFile()) { // change
		line, file = runtimepb.LineFileXUp(2)
	}

	if len(s) < 60 {
		s = stringspb.ToLen(s, 60)
	}
	s = fmt.Sprintf("%v - %v:%v", s, file, line)

	// Log it
	c, _ := util_appengine.SafelyExtractGaeCtxError(r)
	if c == nil {
		lnp.Printf(s)
	} else {
		c.Infof(s)
	}

}

type tLogFunc func(format string, is ...interface{})
type tErrFunc func(error)

var lnp = log.New(os.Stdout, "", 0) // logger no prefix

const maxPref = 32

func Logger(w http.ResponseWriter, r *http.Request) (tLogFunc, tErrFunc) {

	fLog := func(format string, is ...interface{}) {
		Pf(w, r, format, is...)
	}
	fErr := func(err error) {
		if err != nil {
			Pf(w, r, "Err %v", err)
		}
	}
	return fLog, fErr

}

// universal logger func, for log(err) and log("format", ...)
type FuncBufUniv func(...interface{})

func BuffLoggerUniversal(w http.ResponseWriter, r *http.Request) (FuncBufUniv, *bytes.Buffer) {

	b := new(bytes.Buffer)

	fLog1 := func(a ...interface{}) {

		if len(a) > 0 {
			switch t := a[0].(type) {

			case string:
				if len(a) == 1 {
					Pf(b, r, t)
				} else {
					Pf(b, r, t, a[1:]...)
				}

			case interface{}:
				if err, ok := t.(error); ok {
					if err != nil {
						Pf(b, r, "Err %v", err)
					}
				} else {
					log.Printf("first argument must be string or error occ1; is %T\n", t)
				}

			default:
				log.Printf("first argument must be string or error occ2; is %T\n", t)
			}

		}

	}
	return fLog1, b

}
