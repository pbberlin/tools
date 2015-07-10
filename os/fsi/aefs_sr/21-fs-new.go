package aefs_sr

import (
	"strings"

	"github.com/pbberlin/tools/os/fsi"

	"appengine"
)

// AeContext is an option func, adding ae context to the filesystem
func AeContext(c appengine.Context) func(fsi.FileSystem) {
	return func(fs fsi.FileSystem) {
		fst := fs.(*AeFileSys)
		fst.c = c
	}
}

// Rooted is an option func, switching from nested to rooted storage architecture
func Rooted(isRooted bool) func(fsi.FileSystem) {
	return func(fs fsi.FileSystem) {
		fst := fs.(*AeFileSys)
		fst.rooted = isRooted
	}
}

// NewAeFs creates a new appengine datastore filesystem.
// Notice that variadic options are submitted as functions,
// as is explained and justified here:
// http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func NewAeFs(mount string, options ...func(fsi.FileSystem)) *AeFileSys {

	fs := AeFileSys{}

	if strings.Contains(mount, "/") {
		panic("mount can't have slash in it")
	}
	fs.mount = mount

	for _, option := range options {
		option(&fs)
	}

	if fs.c == nil {
		panic("this type of filesystem needs appengine context")
	}

	var err error
	fs.rootDir, err = fs.saveDirUnderParent(fs.mount, nil)
	if err != nil {
		panic(spf("%v", err))
	}

	return &fs
}

func (fs *AeFileSys) Ctx() appengine.Context {
	return fs.c
}

func (fs *AeFileSys) RootDir() string {
	return fs.mount + sep
}
