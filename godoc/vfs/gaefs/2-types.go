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

	rooted bool // default would be nested; <nested, rooted>

	RootDir AeDir
	Mount   string // name of mount point, for remount
	Opener         // implicit
}

type AeDir struct {
	Fs      *AeFileSys `datastore:"-" json:"-"` // Reference to root
	dir     string
	name    string // BaseeName - distinct from os.FileInfo method Name()
	isDir   bool   // distinct from os.FileInfo method IsDir()
	modTime time.Time
	mode    os.FileMode

	Key  *ds.Key `datastore:"-" json:"-"` // throw out? available anyway.
	SKey string  // readable form; not from *ds.Key.Encode()
}

type AeFile struct {
	Fs      *AeFileSys `datastore:"-" json:"-"` // Reference to root
	dir     string     // memDir  MemDir
	name    string     // BaseeName - distinct from os.FileInfo method Name()
	isDir   bool       // distinct from os.FileInfo method IsDir()
	modTime time.Time
	mode    os.FileMode

	data []byte
	sync.Mutex
	at     int64
	closed bool // default open

	Key  *ds.Key `datastore:"-" json:"-"` // throw out? available anyway.
	SKey string  // readable form; not from *ds.Key.Encode()
}

func NewFs(mount string, c appengine.Context, rooted bool) AeFileSys {
	fs := AeFileSys{}
	// fs.Opener = fs.Open // implicit
	fs.c = c
	fs.rooted = rooted
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
