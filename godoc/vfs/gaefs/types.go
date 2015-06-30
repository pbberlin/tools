package gaefs

import (
	"strings"
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
}

type Directory struct {
	Fs    *FileSys `datastore:"-" json:"-"` // Reference to root
	Dir   string
	Name  string
	IsDir bool
	Mod   time.Time

	Key  *ds.Key `datastore:"-" json:"-"` // throw out? available anyway.
	SKey string  // readable form; not from *ds.Key.Encode()
}

type File struct {
	Fs    *FileSys `datastore:"-" json:"-"` // Reference to root
	Dir   string
	Name  string
	IsDir bool
	Mod   time.Time

	Key     *ds.Key `datastore:"-" json:"-"` // throw out? available anyway.
	SKey    string  // readable form; not from *ds.Key.Encode()
	Content []byte
}

func NewFs(mount string, c appengine.Context, rooted bool) FileSys {
	fs := FileSys{}
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
