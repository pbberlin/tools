package filesys

import (
	"net/http"
	"strings"
	"time"

	"appengine"

	ds "appengine/datastore"
)

// Make it pluggable between nested and rooted
// Currently rooted is prefixed with rooted...
type LowLevelArchitecture interface {
	getDirByPath(string) (Directory, error)
	saveDirByPath(string) (Directory, error)
}

// Filesystem
type FileSys struct {
	w http.ResponseWriter `datastore:"-" json:"-"`
	r *http.Request       `datastore:"-" json:"-"`
	c appengine.Context   `datastore:"-" json:"-"`

	RootDir Directory
	Mount   string // name of mount point, for remount
	LowLevelArchitecture
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

func NewFileSys(w http.ResponseWriter, r *http.Request, mount string) FileSys {
	fs := FileSys{}
	fs.w = w
	fs.r = r
	fs.c = appengine.NewContext(r)
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
