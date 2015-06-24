package pblog

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// Since our helper funcs codeline and is misleading
// we create a logger with time only.
// Source File and code is added for caller(1) by the helpers.
var lgDeeper *log.Logger

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("#")

	lgDeeper = log.New(os.Stdout, "#", 0)
}

// LogE mostly saves the if err != nil
//
func LogE(e error, msg ...string) bool {
	if e != nil {
		var s string
		_, file, line, _ := runtime.Caller(1)
		if len(msg) > 0 {
			s = fmt.Sprintf("ERR: %v - %v  \n\tSRC: %s:%d \n", msg[0], e, file, line)
		} else {
			s = fmt.Sprintf("ERR: %v  \n\tSRC: %s:%d \n", e, file, line)
		}
		lgDeeper.Printf(s)
		return true
	}
	return false
}

// Fatal exits.
// It does not panic, to save us from the goroutine dumps
func Fatal(e error, msg ...string) {
	ret := LogE(e, msg...)
	if ret == true {
		os.Exit(1)
	}
}
