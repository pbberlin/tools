package dsfs

import (
	"github.com/pbberlin/tools/os/fsi"

	"appengine"
	"appengine/datastore"
)

// AeContext is an option func, adding ae context to the filesystem
func AeContext(c appengine.Context) func(fsi.FileSystem) {
	return func(fs fsi.FileSystem) {
		fst := fs.(*dsFileSys)
		fst.c = c
	}
}

// MountName is an option func, adding a specific mount name to the filesystem
func MountName(mnt string) func(fsi.FileSystem) {
	return func(fs fsi.FileSystem) {
		fst := fs.(*dsFileSys)
		fst.mount = mnt
	}
}

// New creates a new appengine datastore filesystem.
// Notice that variadic options are submitted as functions,
// as is explained and justified here:
// http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func New(options ...func(fsi.FileSystem)) *dsFileSys {

	fs := dsFileSys{}

	for _, option := range options {
		option(&fs)
	}

	if fs.mount == "" {
		fs.mount = MountPointLast()
	}

	if fs.c == nil {
		panic("this type of filesystem needs appengine context, submitted as option")
	}

	rt, err := fs.dirByPath(fs.mount)
	_ = rt
	if err == datastore.ErrNoSuchEntity {
		// log.Printf("need to creat root %v", fs.mount)
		_, err := fs.saveDirByPath(fs.mount) // fine
		if err != nil {
			fs.c.Errorf("could not create mount %v => %v", fs.mount, err)
		}
	} else if err != nil {
		fs.c.Errorf("could read mount dir %v => %v", fs.mount, err)
	}

	return &fs
}

func (fs *dsFileSys) Ctx() appengine.Context {
	return fs.c
}

func (fs *dsFileSys) RootDir() string {
	return fs.mount + sep
}

func (fs *dsFileSys) RootName() string {
	return fs.mount
}

func Unwrap(fs fsi.FileSystem) (*dsFileSys, bool) {
	fsc, ok := fs.(*dsFileSys)
	return fsc, ok
}
