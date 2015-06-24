package pblog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// Since our helper funcs codeline and is misleading
// we create a logger with time only.
// Source File and code is added for caller(1) by the helpers.
var lg1 *log.Logger

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("#")

	lg1 = log.New(os.Stdout, "#", 0)
}

// LogE mostly saves the if err != nil
//
func LogE(e error, msg ...string) bool {
	if e != nil {
		inner(e, msg...)
		return true
	}
	return false
}

// Fatal exits.
// It does not panic, to save us from the goroutine dumps
// Identical to LogE, except for last lines.
// We cannot do nesting, 'cause runtime.Caller(1) would get shifted.
func Fatal(e error, msg ...string) {
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
	lg1.Printf(s)
}
