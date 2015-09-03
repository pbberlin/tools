package memfs

import (
	"os"
	"sort"
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
	fos map[string]fsi.File

	mtx *sync.RWMutex // syncing fos

	// Could be expanded into a kind of "currentDir"
	// for compatibility with osfs.
	// Currently only used as instance name
	ident         string
	readdirsorter func([]os.FileInfo)
	shadow        fsi.FileSystem
}

// Ident is an option func, adding a specific identification to the filesystem
func Ident(mnt string) func(fsi.FileSystem) {
	return func(fs fsi.FileSystem) {
		fst := fs.(*memMapFs)
		fst.ident = mnt
	}
}

// Ident is an option func, adding a specific identification to the filesystem
func ShadowFS(fs fsi.FileSystem) func(fsi.FileSystem) {
	return func(fs fsi.FileSystem) {
		fst := fs.(*memMapFs)
		fst.shadow = fs
	}
}

// Default sort for ReadDir... is ByNameAsc
// We may want to change this; for instance sort byDate
func DirSort(srt string) func(fsi.FileSystem) {
	return func(fs fsi.FileSystem) {
		fst := fs.(*memMapFs)

		switch srt {
		case "byDateAsc":
			fst.readdirsorter = func(fis []os.FileInfo) { sort.Sort(byDateAsc(fis)) }

		case "byDateDesc":
			fst.readdirsorter = func(fis []os.FileInfo) { sort.Sort(byDateDesc(fis)) }
		case "byName":
			fst.readdirsorter = func(fis []os.FileInfo) { sort.Sort(byName(fis)) }
		}
	}
}

// New creates a in-memory filesystem.
// Notice that variadic options are submitted as functions,
// as is explained and justified here:
// http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func New(options ...func(fsi.FileSystem)) *memMapFs {
	m := &memMapFs{
		fos:           map[string]fsi.File{}, // secure init
		ident:         "mnt00",
		readdirsorter: func(fis []os.FileInfo) { sort.Sort(byName(fis)) },
	}
	for _, option := range options {
		option(m)
	}
	return m
}

func (m *memMapFs) RootDir() string {
	return m.ident + sep
}

func (m *memMapFs) RootName() string {
	return m.ident
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
