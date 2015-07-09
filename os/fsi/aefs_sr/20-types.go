package aefs_sr

import (
	"os"
	"sync"
	"time"

	"appengine"

	ds "appengine/datastore"
)

// Filesystem
type AeFileSys struct {
	// w http.ResponseWriter `datastore:"-" json:"-"`
	// r *http.Request       `datastore:"-" json:"-"`
	c appengine.Context `datastore:"-" json:"-"`

	rooted bool // default would be nested; <nested, rooted>

	mount   string // name of mount point, for remount
	rootDir AeDir  // a directory with the same name
}

// Upper case field names sadly
// inevitable, for ae datastore :(
type AeDir struct {
	Fs   *AeFileSys `datastore:"-" json:"-"` // Reference to root
	SKey string     // readable form; not *ds.Key.Encode()
	Key  *ds.Key    `datastore:"-" json:"-"` // throw out? Can be constructed from SKey

	Dir      string
	BName    string      // BaseName - distinct from os.FileInfo method Name()
	isDir    bool        // distinct from os.FileInfo method IsDir()
	MModTime time.Time   `datastore:"ModTime" json:"ModTime"`
	MMode    os.FileMode `datastore:"-" json:"-"` // SaveProperty must be implemented
}

// Upper case field names sadly
// inevitable, for ae datastore :(
type AeFile struct {
	Fs   *AeFileSys `datastore:"-" json:"-"` // Reference to root
	SKey string     // readable form; not *ds.Key.Encode()
	Key  *ds.Key    `datastore:"-" json:"-"` // throw out? Can be constructed from SKey.

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
