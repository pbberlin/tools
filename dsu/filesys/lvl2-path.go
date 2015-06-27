package filesys

import (
	"path/filepath"
	"strings"

	ds "appengine/datastore"
)

func dir(p string) string {
	return (filepath.Dir(p))
}

func (fs FileSys) mkdir(path string) {

	// base := filepath.Base(path)
	// par := filepath.Dir(path)

	// newDir := fs.newFso(base, par, true)

}

func (f FileSys) touch(p string) {

}

func (f FileSys) CreatePath(path string) {

	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	par := dir(path)
	f.recurseMkDir(par)

	f.touch(path)
}

func (f FileSys) recurseMkDir(path string) {
	if path == "" {
		return
	}
	par, err := f.GetFso(path)
	_ = par
	if err == ds.ErrNoSuchEntity {
		par := dir(path)
		f.recurseMkDir(par)

		f.mkdir(path)
	}
}
