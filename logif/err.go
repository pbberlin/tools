// Package logif wraps logging errors, saving some err != nil boilerplate.
package logif

import (
	"fmt"
	"log"
	"os"

	"github.com/pbberlin/tools/net/http/loghttp"
	"github.com/pbberlin/tools/runtimepb"
)

func init() {
	restoreLogFlags()
	log.SetPrefix("")
}

func restoreLogFlags() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.SetFlags(log.Lshortfile)
	log.SetFlags(0)
}

var lnp = log.New(os.Stderr, "", 0) // logger no prefix; os.Stderr shows up in appengine devserver; os.Stdout does not

// E (Err) mostly saves the if err != nil
//
func E(e error, msg ...string) bool {
	if e != nil {
		inner(e, msg...)
		return true
	}
	return false
}

// F (Fatal) exits.
// It does not panic, to save us from the goroutine dumps
// Identical to LogE, except for last lines.
// Behold the nesting, leading to runtime.Caller(2).
func F(e error, msg ...string) {
	if e != nil {
		inner(e, msg...)
		os.Exit(1)
	}
}

func inner(e error, msg ...string) {
	var s string
	// _, file, line, _ := runtime.Caller(2) // TWO steps up
	// dir := filepath.Dir(file)
	// dirLast := filepath.Base(dir)
	// file = filepath.Join(dirLast, filepath.Base(file))

	line, file := runtimepb.LineFileXUp(2) // TWO steps up, inner() and logif.F / logif.E
	if len(msg) > 0 {
		s = fmt.Sprintf("ERR: %v - %v  \n\tSRC: %s:%d ", msg[0], e, file, line)
	} else {
		s = fmt.Sprintf("ERR: %v  \n\tSRC: %s:%d ", e, file, line)
	}

	// Since codeline points to *this* helper-func
	// we would like to logger with time only.
	if loghttp.C == nil {
		lnp.Printf(s)
	} else {
		// This is of course criminal,
		// since loghttp.C is not syncronized.
		// The message might appear under a *wrong* request.
		// But it's the only way to make
		// ordinary log messages available.
		loghttp.C.Infof(fmt.Sprintf("%s - volat req assign", s))
	}

}

func Pf(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	line, file := runtimepb.LineFileXUp(1)
	s = fmt.Sprintf("%v - %v:%v", s, file, line)
	if loghttp.C == nil {
		lnp.Printf(s)
	} else {
		loghttp.C.Infof(fmt.Sprintf("%s - volat req assign", s))
	}
}
