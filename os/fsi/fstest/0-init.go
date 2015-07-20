// package fstest runs tests on all fsi subpackages;
// some of which require an appengine context;
// and which must be run by goapp test.
package fstest

import (
	"fmt"
	"io"
)

var spf func(format string, a ...interface{}) string = fmt.Sprintf
var wpf func(w io.Writer, format string, a ...interface{}) (int, error) = fmt.Fprintf

var dot = []string{
	"fs.go",
	"fs_test.go",
	"httpFs.go",
	"memfile.go",
	"memmap.go",
}

var testDir = "/temp/fun"
var testName = "testF.txt"
