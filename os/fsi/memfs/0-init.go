package memfs

import (
	"os"
	"sync"

	"github.com/pbberlin/tools/os/fsi"
)

var mux = &sync.Mutex{}

func init() {

	// forcing our implementations
	// to comply with our interfaces

	f := InMemoryFile{}
	ifa := fsi.File(&f)
	_ = ifa

	fi := InMemoryFileInfo{}
	ifi := os.FileInfo(&fi)
	_ = ifi

	fs := MemMapFs{}
	ifs := fsi.FileSystem(&fs)
	_ = ifs

}
