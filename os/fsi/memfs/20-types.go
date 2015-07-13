package memfs

import (
	"os"
	"sync"
	"time"

	"github.com/pbberlin/tools/os/fsi"
)

// It's unexported; only New() is possible.
// If access to underlying methods neccessary from outside, use Unwrap()
type memMapFs struct {
	fos   map[string]fsi.File // file objects - a vector of all files and directories
	mutex *sync.RWMutex
}

func New() *memMapFs {
	m := &memMapFs{
		fos: map[string]fsi.File{}, //
	}
	return m
}

func Unwrap(fs fsi.FileSystem) (*memMapFs, bool) {
	fsc, ok := fs.(*memMapFs)
	return fsc, ok
}

type InMemoryFile struct {
	sync.Mutex
	at      int64
	name    string
	data    []byte
	memDir  MemDir // directory contents
	dir     bool
	closed  bool
	mode    os.FileMode
	modtime time.Time
}

type InMemoryFileInfo struct {
	file *InMemoryFile
}

// Implemented by MemDirMap
type MemDir interface {
	Add(fsi.File)
	Len() int
	Files() []fsi.File
	Names() []string
	Remove(fsi.File)
}

type MemDirMap map[string]fsi.File // implements MemDir
