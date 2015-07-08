package gaefs

import (
	"errors"
	"fmt"
	"os"
)

// the usual short notations for fmt.Printf and fmt.Sprintf
var pf func(format string, a ...interface{}) (int, error) = fmt.Printf
var pfRestore func(format string, a ...interface{}) (int, error) = fmt.Printf
var spf func(format string, a ...interface{}) string = fmt.Sprintf

const (
	tdir    = "fsd"      // datastory entity type for filesystem directory
	tdirsep = tdir + "," // nested datastore keys each have this prefix
	tfil    = "fsf"      // datastory entity type for filesystem file
	sep     = "/"        // no, package path does not provide it; yes, we do need it.
)

var (
	ErrFileClosed = errors.New("File is closed")
	ErrFileInUse  = errors.New("File already in use")
	ErrOutOfRange = errors.New("Out of range")
	ErrTooLarge   = errors.New("Too large")

	// can't those be replaced by the original?
	ErrFileNotFound      = os.ErrNotExist
	ErrFileExists        = os.ErrExist
	ErrDestinationExists = os.ErrExist
)

func init() {

	// forcing our implementations
	// to comply with our interfaces

	f := AeFile{}
	ifa := FileI(&f)
	_ = ifa

	ifi := os.FileInfo(&f)
	_ = ifi

	fs := AeFileSys{}
	ifs := FileSystem(&fs)
	_ = ifs

}
