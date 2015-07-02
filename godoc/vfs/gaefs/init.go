package gaefs

import "fmt"

var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var pfRestore func(format string, a ...interface{}) (int, error) = fmt.Printf
var spf func(format string, a ...interface{}) string = fmt.Sprintf

var (
	tdir       string = "fsd" // FileSys filesystem directory
	tdirsep    string = tdir + ","
	tfil       string = "fsf" // FileSys filesystem file
	dummdidumm        = ""
)

const sep = "/" // no, package path does not provide it; yes, we do need it.
