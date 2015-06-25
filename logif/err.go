package logif

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"appengine"
)

func init() {
	setFlags()
	log.SetPrefix("#")
}

func setFlags() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.SetFlags(log.Lshortfile)
}

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
	_, file, line, _ := runtime.Caller(2) // TWO steps up
	dir := filepath.Dir(file)
	dirLast := filepath.Base(dir)
	file = filepath.Join(dirLast, filepath.Base(file))
	if len(msg) > 0 {
		s = fmt.Sprintf("ERR: %v - %v  \n\tSRC: %s:%d \n", msg[0], e, file, line)
	} else {
		s = fmt.Sprintf("ERR: %v  \n\tSRC: %s:%d \n", e, file, line)
	}

	// Since codeline points to *this* helper-func
	// we would like to logger with time only.
	//   lg1 = log.New(os.Stdout, "#", 0)
	// but it would not be  written under appengine, because of os.Stdout

	log.SetFlags(0)
	log.Printf(s)
	setFlags()

}

func SafeGaeCheck(r *http.Request) (appengine.Context, error) {
	c := checkPanicking(r)
	if c != nil {
		return c, nil
	} else {
		return nil, fmt.Errorf("Request is not appengine")
	}
}

func checkPanicking(r *http.Request) appengine.Context {
	defer func() {
		recover()
	}()
	c := appengine.NewContext(r)
	return c
}
