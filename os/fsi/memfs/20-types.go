package memfs

import (
	"os"
	"sync"
	"time"

	"github.com/pbberlin/tools/os/fsi"
)

type MemMapFs struct {
	data  map[string]fsi.File
	mutex *sync.RWMutex
}

type InMemoryFile struct {
	sync.Mutex
	at      int64
	name    string
	data    []byte
	memDir  MemDir
	dir     bool
	closed  bool
	mode    os.FileMode
	modtime time.Time
}

type InMemoryFileInfo struct {
	file *InMemoryFile
}

type MemDirMap map[string]fsi.File

type MemDir interface {
	Len() int
	Names() []string
	Files() []fsi.File
	Add(fsi.File)
	Remove(fsi.File)
}
