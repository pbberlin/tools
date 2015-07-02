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
type FileSys struct {
	// w http.ResponseWriter `datastore:"-" json:"-"`
	// r *http.Request       `datastore:"-" json:"-"`
	c appengine.Context `datastore:"-" json:"-"`

	rooted bool // default would be nested; <nested, rooted>

	RootDir Directory
	Mount   string // name of mount point, for remount
	Opener         // implicit
}

type Directory struct {
	Fs          *FileSys `datastore:"-" json:"-"` // Reference to root
	Dir         string
	BName       string // BaseeName - distinct from os.FileInfo method Name()
	IsDirectory bool   // distinct from os.FileInfo method IsDir()
	Mod         time.Time

	Key  *ds.Key `datastore:"-" json:"-"` // throw out? available anyway.
	SKey string  // readable form; not from *ds.Key.Encode()
}

type File struct {
	Fs          *FileSys `datastore:"-" json:"-"` // Reference to root
	Dir         string
	BName       string // BaseeName - distinct from os.FileInfo method Name()
	IsDirectory bool   // distinct from os.FileInfo method IsDir()
	Mod         time.Time

	Key     *ds.Key `datastore:"-" json:"-"` // throw out? available anyway.
	SKey    string  // readable form; not from *ds.Key.Encode()
	Content []byte

	/*
		for
			io.Closer
			io.Reader
			io.ReaderAt
			io.Seeker
			io.Writer
			io.WriterAt
	*/
	sync.Mutex
	at int64
	// name    string
	// data    []byte
	// memDir  MemDir
	// dir     bool
	closed bool
	mode   os.FileMode
	// modtime time.Time

}

func NewFs(mount string, c appengine.Context, rooted bool) FileSys {
	fs := FileSys{}
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

func (fs *FileSys) Ctx() appengine.Context {
	return fs.c
}
