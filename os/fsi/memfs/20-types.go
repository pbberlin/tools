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

// MountName is an option func, adding a specific mount name to the filesystem
func MountName(mnt string) func(fsi.FileSystem) {
	return func(fs fsi.FileSystem) {
		fst := fs.(*memMapFs)
		fst.mount = mnt
	}
}

// New creates a in-memory filesystem.
// Notice that variadic options are submitted as functions,
// as is explained and justified here:
// http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func New(options ...func(fsi.FileSystem)) *memMapFs {
	m := &memMapFs{
		fos:   map[string]fsi.File{}, // secure init
		mount: "mnt0",
	}
	for _, option := range options {
		option(m)
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
	memDir         map[string]fsi.File
	memDirFetchPos int       // read position for f.Readdir
	fs             *memMapFs // reference to fs
}

// Implements os.FileInfo
type InMemoryFileInfo struct {
	file *InMemoryFile
}
