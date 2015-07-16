package aefs

import (
	"strings"

	"github.com/pbberlin/tools/os/fsi"

	"appengine"
	"appengine/datastore"
)

// AeContext is an option func, adding ae context to the filesystem
func AeContext(c appengine.Context) func(fsi.FileSystem) {
	return func(fs fsi.FileSystem) {
		fst := fs.(*aeFileSys)
		fst.c = c
	}
}

// old
func NewAeFs(mount string, options ...func(fsi.FileSystem)) *aeFileSys {
	return New(mount, options...)
}

// New creates a new appengine datastore filesystem.
// Notice that variadic options are submitted as functions,
// as is explained and justified here:
// http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func New(mount string, options ...func(fsi.FileSystem)) *aeFileSys {

	fs := aeFileSys{}

	if strings.Contains(mount, "/") {
		panic("mount can't have slash in it")
	}
	fs.mount = mount

	for _, option := range options {
		option(&fs)
	}

	if fs.c == nil {
		panic("this type of filesystem needs appengine context, submitted as option")
	}

	rt, err := fs.dirByPath(mount)
	_ = rt
	if err == datastore.ErrNoSuchEntity {
		// log.Printf("need to creat root %v", mount)
		_, err := fs.saveDirByPath(mount) // fine
		if err != nil {
			fs.c.Errorf("could not create mount %v => %v", mount, err)
		}
	} else if err != nil {
		fs.c.Errorf("could read mount dir %v => %v", mount, err)
	}

	return &fs
}

func (fs *aeFileSys) Ctx() appengine.Context {
	return fs.c
}

func (fs *aeFileSys) RootDir() string {
	return fs.mount + sep
}

func (fs *aeFileSys) RootName() string {
	return fs.mount
}

func Unwrap(fs fsi.FileSystem) (*aeFileSys, bool) {
	fsc, ok := fs.(*aeFileSys)
	return fsc, ok
}
