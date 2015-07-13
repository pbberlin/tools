// Package osfs provides access to the file system
// of the OS; merely wrapping os and ioutil
package osfs

import (
	"os"

	"github.com/pbberlin/tools/os/fsi"
)

func init() {

	// forcing our implementations
	// to comply with our interfaces

	f := os.File{}
	ifa := fsi.File(&f)
	_ = ifa

	var fi os.FileInfo
	ifi := os.FileInfo(fi) // of course idiotic, but we keep the pattern
	_ = ifi

	fs := osFileSys{}
	ifs := fsi.FileSystem(&fs)
	_ = ifs

}
