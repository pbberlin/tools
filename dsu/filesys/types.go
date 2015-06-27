package filesys

import (
	"net/http"
	"os"
	"time"

	"appengine"

	ds "appengine/datastore"
)

type FileInfo interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() os.FileMode  // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)

	Parent() string // Added
}

// Filesystem
type FileSys struct {
	w http.ResponseWriter `datastore:"-"`
	r *http.Request       `datastore:"-"`
	c appengine.Context   `datastore:"-"`

	RootDir FSysObj
}

// Filesystem Object - a directory or a file
type FSysObj struct {
	fs    *FileSys
	Name  string
	IsDir bool
	Mod   time.Time

	Key  *ds.Key
	SKey string // from *ds.Key.Encode()
}
