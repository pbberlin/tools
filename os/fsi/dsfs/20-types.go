package dsfs

import (
	"os"
	"sync"
	"time"

	"appengine"

	ds "appengine/datastore"
)

// Filesystem
type dsFileSys struct {
	// w http.ResponseWriter `datastore:"-" json:"-"`
	// r *http.Request       `datastore:"-" json:"-"`
	c appengine.Context `datastore:"-" json:"-"`

	mount string // name of mount point, for remount

	dirsorter  func([]os.FileInfo)
	filesorter func([]DsFile)
}

// The distinction between AeDir and AeFile
// brings clarity into the low-level implementation.
// And into the google datastore overlay architecture.
//
// However, AeDir needs many redundant methods,
// since it must be convertable into fsi.File in Open()

// Upper case field names sadly
// inevitable, for ae datastore :(
type DsDir struct {
	fSys *dsFileSys `datastore:"-" json:"-"` // Reference to root
	Key  *ds.Key    `datastore:"-" json:"-"` // throw out? Can be constructed from Dir+BName

	Dir      string
	BName    string      // BaseName - distinct from os.FileInfo method Name()
	isDir    bool        // distinct from os.FileInfo method IsDir()
	MModTime time.Time   `datastore:"ModTime" json:"ModTime"`
	MMode    os.FileMode `datastore:"-" json:"-"` // SaveProperty must be implemented

	memDirFetchPos int // read position for f.Readdir
}

// Upper case field names sadly
// inevitable, for ae datastore :(
type DsFile struct {
	fSys *dsFileSys `datastore:"-" json:"-"` // Reference to root
	Key  *ds.Key    `datastore:"-" json:"-"` // throw out? Can be constructed from Dir+BName.

	Dir      string
	BName    string      // BaseName - distinct from os.FileInfo method Name()
	isDir    bool        // distinct from os.FileInfo method IsDir()
	MModTime time.Time   `datastore:"ModTime" json:"ModTime"`
	MMode    os.FileMode `datastore:"-" json:"-"` // SaveProperty must be implemented

	Data []byte `datastore:"Data" json:"Data"`
	sync.Mutex
	at     int64
	closed bool // default open

	memDirFetchPos int // read position for f.Readdir

}
