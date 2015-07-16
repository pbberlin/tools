// package fstest runs tests on all fsi subpackages;
// some of which require an appengine context;
// and which must be run by goapp test.
package fstest

import (
	"fmt"
	"io"
	"log"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/aefs"
	"github.com/pbberlin/tools/os/fsi/memfs"
	"github.com/pbberlin/tools/os/fsi/osfs"

	"appengine/aetest"
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

func initFileSystems() (fss []fsi.FileSystem, c aetest.Context) {

	var err error
	c, err = aetest.NewContext(nil)
	if err != nil {
		log.Fatal(err)
	}

	// defer c.Close()
	// Do not here
	// but instead at the start of the test-funcs

	// We cant make variadic options generic,
	// since they need the concrete filesystem type.
	fs1 := aefs.New(aefs.MountPointNext(), aefs.AeContext(c))

	fs3 := osfs.New()

	fs4 := memfs.New(memfs.MountName("m"))

	fss = []fsi.FileSystem{fs1, fs3, fs4}

	return fss, c
}
