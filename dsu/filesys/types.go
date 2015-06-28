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
	w http.ResponseWriter `datastore:"-"`
	r *http.Request       `datastore:"-"`
	c appengine.Context   `datastore:"-"`

	RootDir Directory `datastore:"-"`
	Mount   string    // name of mount point, for remount
	LowLevelArchitecture
}

type Directory struct {
	Fs    *FileSys `datastore:"-"` // Reference to root
	Dir   string
	Name  string
	IsDir bool
	Mod   time.Time

	Key  *ds.Key
	SKey string // from *ds.Key.Encode()
}

type File struct {
	Fs    *FileSys `datastore:"-"` // Reference to root
	Dir   string
	Name  string
	IsDir bool
	Mod   time.Time

	Key  *ds.Key
	SKey string // from *ds.Key.Encode()
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
