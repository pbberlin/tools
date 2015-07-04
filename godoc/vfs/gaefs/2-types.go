package gaefs

import (
	"os"
	"strings"
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

	Rooted bool // default would be nested; <nested, rooted>

	RootDir AeDir
	Mount   string // name of mount point, for remount
	// Opener         // forcing implementation of Open()
}

// Upper case field names sadly
// inevitable, for ae datastore :(
type AeDir struct {
	Fs       *AeFileSys `datastore:"-" json:"-"` // Reference to root
	Dir      string
	BName    string      // BaseName - distinct from os.FileInfo method Name()
	isDir    bool        // distinct from os.FileInfo method IsDir()
	MModTime time.Time   `datastore:"ModTime" json:"ModTime"`
	MMode    os.FileMode `datastore:"-" json:"-"` // SaveProperty must be implemented

	SKey string  // readable form; not from *ds.Key.Encode()
	Key  *ds.Key `datastore:"-" json:"-"` // throw out? Can be constructed from SKey
}

// Upper case field names sadly
// inevitable, for ae datastore :(
type AeFile struct {
	Fs       *AeFileSys `datastore:"-" json:"-"` // Reference to root
	Dir      string
	BName    string      // BaseName - distinct from os.FileInfo method Name()
	isDir    bool        // distinct from os.FileInfo method IsDir()
	MModTime time.Time   `datastore:"ModTime" json:"ModTime"`
	MMode    os.FileMode `datastore:"-" json:"-"` // SaveProperty must be implemented

	Data []byte `datastore:"Data" json:"Data"`
	sync.Mutex
	at     int64
	closed bool // default open

	SKey string  // readable form; not from *ds.Key.Encode()
	Key  *ds.Key `datastore:"-" json:"-"` // throw out? Can be constructed from SKey.
}

func NewFs(mount string, c appengine.Context, rooted bool) AeFileSys {
	fs := AeFileSys{}
	// fs.Opener = fs.Open // implicit
	fs.c = c
	fs.Rooted = rooted
	if strings.Contains(mount, "/") {
		panic("mount can't have slash in it")
	}
	fs.Mount = mount

	var err error
	fs.RootDir, err = fs.saveDirUnderParent(fs.Mount, nil)
	if err != nil {
		panic(spf("%v", err))
	}
	return fs
}

func (fs *AeFileSys) Ctx() appengine.Context {
	return fs.c
}
