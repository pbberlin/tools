package aefs

import (
	"fmt"
	"io"
	"os"

	"github.com/pbberlin/tools/os/fsi"
)

// the usual short notations for fmt.Printf and fmt.Sprintf
var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var pfRestore func(format string, a ...interface{}) (int, error) = fmt.Printf
var spf func(format string, a ...interface{}) string = fmt.Sprintf
var wpf func(w io.Writer, format string, a ...interface{}) (int, error) = fmt.Fprintf

const (
	tdir      = "fsd"      // datastory entity type for filesystem directory
	tdirsep   = tdir + "," // nested datastore keys each have this prefix
	tfil      = "fsf"      // datastory entity type for filesystem file
	sep       = "/"        // no, package path does not provide it; yes, we do need it.
	doublesep = "//"
)

func init() {

	// forcing our implementations
	// to comply with our interfaces

	f := AeFile{}
	ifa := fsi.File(&f)
	_ = ifa

	ifi := os.FileInfo(&f)
	_ = ifi

	fs := aeFileSys{}
	ifs := fsi.FileSystem(&fs)
	_ = ifs

}
