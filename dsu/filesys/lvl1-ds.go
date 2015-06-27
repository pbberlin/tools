package filesys

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"appengine"
	ds "appengine/datastore"
)

var t string

func init() {
	fo := FSysObj{}
	t = fmt.Sprintf("%T", fo) // "kind"
	t = "fso"
}

func NewFileSys(w http.ResponseWriter, r *http.Request, root string) FileSys {
	fs := FileSys{}
	fs.w = w
	fs.r = r
	fs.c = appengine.NewContext(r)

	fs.RootDir = fs.newFso(root, nil, true)

	return fs
}

func (fs FileSys) getFso(path string) (FSysObj, error) {

	base := filepath.Base(path)
	par := filepath.Dir(path)

	PKey := ds.NewKey(fs.c, t, fs.RootDir.SKey+","+par, 0, nil)
	suggKey := ds.NewKey(fs.c, t, base, 0, PKey)
	fo := FSysObj{}
	err := ds.Get(fs.c, suggKey, &fo)

	if err != nil {
		return fo, nil
	} else {
		return fo, err
	}
}

func (fs FileSys) newFso(name string, parent *ds.Key, isDir bool) FSysObj {

	fo := FSysObj{}
	fo.IsDir = isDir
	fo.Name = name
	fo.Mod = time.Now()
	fo.fs = &fs

	suggKey := ds.NewKey(fs.c, t, name, 0, parent)
	fo.Key = suggKey
	fo.SKey = spf("%v", suggKey) // not effKey.Encode()

	effKey, err := ds.Put(fs.c, suggKey, &fo)

	if err != nil {
		fs.c.Errorf("%v  [%s]", err, fo.Key)
	} else {
		fs.c.Infof("fso ds-saved - key %v [%v]", effKey, effKey.Encode())
	}

	fo.Key = effKey
	fo.SKey = spf("%v", effKey) // not effKey.Encode()

	fo.MemCacheSet()

	return fo
}

func (fs FileSys) GetFso(path string) (FSysObj, error) {

	fo := FSysObj{}

	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	if path == "" {
		return fs.RootDir, nil
	}

	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-2]
	}

	dir, file := filepath.Split(path)

	if dir == "" && file == "" {
		return fs.RootDir, nil
	}

	fs.c.Infof("%-44v --------- %-28v -- %v", path, dir, file)

	par, err := fs.GetFso(dir) // recurse
	if err != nil {
		fs.c.Errorf("err getting fso: %v", err)
		return fo, err
	}

	err = ds.Get(fs.c, par.Key, &fo)
	if err != nil {
		fs.c.Errorf("err fetching parent: %v", err)
		return fo, err
	}

	return fo, nil
}
