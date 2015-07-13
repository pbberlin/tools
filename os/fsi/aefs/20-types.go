package aefs

import (
	"os"
	"sync"
	"time"

	"appengine"

	ds "appengine/datastore"
)

// Filesystem
type aeFileSys struct {
	// w http.ResponseWriter `datastore:"-" json:"-"`
	// r *http.Request       `datastore:"-" json:"-"`
	c appengine.Context `datastore:"-" json:"-"`

	mount string // name of mount point, for remount
}

// Upper case field names sadly
// inevitable, for ae datastore :(
type AeDir struct {
	fSys *aeFileSys `datastore:"-" json:"-"` // Reference to root
	Key  *ds.Key    `datastore:"-" json:"-"` // throw out? Can be constructed from Dir+BName

	Dir      string
	BName    string      // BaseName - distinct from os.FileInfo method Name()
	isDir    bool        // distinct from os.FileInfo method IsDir()
	MModTime time.Time   `datastore:"ModTime" json:"ModTime"`
	MMode    os.FileMode `datastore:"-" json:"-"` // SaveProperty must be implemented
}

// Upper case field names sadly
// inevitable, for ae datastore :(
type AeFile struct {
	fSys *aeFileSys `datastore:"-" json:"-"` // Reference to root
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

}
