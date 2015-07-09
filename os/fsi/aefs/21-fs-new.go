package aefs

import (
	"strings"

	"github.com/pbberlin/tools/logif"

	"appengine"
	"appengine/datastore"
)

// AeContext is an option func, adding ae context to the filesystem
func AeContext(c appengine.Context) func(*AeFileSys) {
	return func(fs *AeFileSys) {
		fs.c = c
	}
}

// NewAeFs creates a new appengine datastore filesystem.
// Notice that variadic options are submitted as functions,
// as is explained and justified here:
// http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func NewAeFs(mount string, options ...func(*AeFileSys)) *AeFileSys {

	fs := AeFileSys{}

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

	_, err := fs.dirByPath(mount)
	if err == datastore.ErrNoSuchEntity {
		_, err := fs.saveDirByPath(mount)
		logif.F(err)
	}

	return &fs
}

func (fs *AeFileSys) Ctx() appengine.Context {
	return fs.c
}

func (fs *AeFileSys) RootDir() string {
	return fs.mount + sep
}

func (fs *AeFileSys) RootName() string {
	return fs.mount
}
