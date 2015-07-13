package memfs

import (
	"os"
	"sync"
	"time"

	"github.com/pbberlin/tools/os/fsi"
)

// The main type is unexported.
// Use New().
// For conversion from fsi.FileSystem from outside, use Unwrap()
type memMapFs struct {
	// fos - file objects - a map containing all files and directories.
	// It must be keyed by the full path, otherwise uniqueness suffers
	fos   map[string]fsi.File
	mutex *sync.RWMutex

	// mount is a directory prefix, similar to a base directory.
	// Useful to as a kind of current dir; to keep memfs exchangeable with osfs
	mount string
}

func New() *memMapFs {
	m := &memMapFs{
		fos:   map[string]fsi.File{}, // secure init
		mount: "mnt0",
	}
	return m
}

func (m *memMapFs) RootDir() string {
	return m.mount + sep
}

func (m *memMapFs) RootName() string {
	return m.mount
}

func Unwrap(fs fsi.FileSystem) (*memMapFs, bool) {
	fsc, ok := fs.(*memMapFs)
	return fsc, ok
}

// Implements fsi.File
type InMemoryFile struct {
	sync.Mutex
	at      int64
	closed  bool
	data    []byte
	dir     bool
	mode    os.FileMode
	modtime time.Time
	name    string

	// memDir -  in-memory directory
	// For directories it contains the children;
	// For for files:  it contains siblings.
	memDir map[string]fsi.File
	fs     *memMapFs // reference to fs
}

// Implements os.FileInfo
type InMemoryFileInfo struct {
	file *InMemoryFile
}
