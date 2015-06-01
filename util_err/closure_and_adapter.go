package util_err

import (
	"net/http"

	"strings"

	"os"
	"path/filepath"
)

var AppDir string

/* This is a closure,
  an anonymous func
  with surrounding variables

  myIntSeq01(), myIntSeq02()
	 yield independent values for i

They also demonstrate "static instance memory",
as the "global" variables are kept as long as the app lives

*/
func intSeq() func() int {
	i := 0
	return func() int {
		i += 1
		return i
	}
}

var MyIntSeq01 func() int = intSeq()
var MyIntSeq02 func() int = intSeq()

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

// works like an interface - functions just have to fit in the signature
type ExtendedHandler func(http.ResponseWriter, *http.Request, map[string]interface{})

func setAppDir() {

	sEnviron := os.Environ()
	for _, v := range sEnviron {
		if strings.HasPrefix(v, "PWD=") {
			vS := strings.Split(v, "=")
			AppDir = vS[1]
		}
	}
	AppDir = strings.Replace(AppDir, `\`, `/`, -1)

	for i := 0; i < 2; i++ {
		if len(AppDir) > 1 {
			AppDir, _ = filepath.Split(AppDir[:len(AppDir)-1])
			// log.Println("    app dir is " + AppDir)
		}
	}

}

func init() {
	setAppDir()
}
