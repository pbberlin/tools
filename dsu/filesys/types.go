package filesys

import (
	"net/http"
	"strings"
	"time"

	"appengine"

	ds "appengine/datastore"
)

type LowLevelArchitecture interface {
	dirByPath(string) (Directory, error)
	saveDirByPath(string) Directory
}

// Filesystem
type FileSys struct {
	w http.ResponseWriter `datastore:"-"`
	r *http.Request       `datastore:"-"`
	c appengine.Context   `datastore:"-"`

	RootDir Directory
	mount   string // name of mount point, for remount
	LowLevelArchitecture
}

type Directory struct {
	Fs    *FileSys // Reference to root
	Dir   string
	Name  string
	IsDir bool
	Mod   time.Time

	Key  *ds.Key
	SKey string // from *ds.Key.Encode()
}

type File struct {
	Fs    *FileSys // Reference to root
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
	fs.mount = mount

	// fs.RootDir = fs.newFsoByParentKey(root, nil, true)
	// fs.LowLevelArchitecture = nested.Arch
	// fs.LowLevelArchitecture = rooted.Arch
	return fs
}

func (fs *FileSys) Ctx() appengine.Context {
	return fs.c
}
