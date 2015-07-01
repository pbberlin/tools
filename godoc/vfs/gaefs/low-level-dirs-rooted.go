package gaefs

import (
	"time"

	"github.com/pbberlin/tools/logif"

	pth "path"

	ds "appengine/datastore"
)

// saveDirByPathUnderRoot saves *all* directories under root.
// Retrieval is then possible directly via rootedGetDirByPath()
// or via GetDirByPathQuery().
// The main disadvantage is that "ancestor group updates"
// in google datastore are restricted to frequency "~ five per second".
// If you anticipate fewer directory changes, consider it.
func (fs *FileSys) rootedSaveDirByPath(path string) (Directory, error) {

	if path == "" {
		return fs.RootDir, nil
	}

	fo := Directory{}
	fo.IsDir = true
	dir, base := pth.Split(path)
	fo.Dir = dir
	fo.Name = base
	fo.Mod = time.Now()
	fo.Fs = fs

	perfKey := ds.NewKey(fs.c, tdir, path, 0, fs.RootDir.Key)

	fo.Key = perfKey
	fo.SKey = spf("%v", perfKey) // not effKey.Encode()

	effKey, err := ds.Put(fs.c, perfKey, &fo)

	if err != nil {
		logif.E(err)
		return fo, err
	}
	if !perfKey.Equal(effKey) {
		fs.Ctx().Errorf("keys unequal %v - %v", perfKey, effKey)
	}

	fo.MemCacheSet()

	return fo, nil
}

// Retrieves a directory in one go.
// But only if it was saved with rootedSaveDirBypth.
func (fs *FileSys) rootedGetDirByPath(path string) (Directory, error) {

	if path == "" {
		return fs.RootDir, nil
	}

	fo := Directory{}
	fo.Fs = fs
	perfKey := ds.NewKey(fs.c, tdir, path, 0, fs.RootDir.Key)
	fo.Key = perfKey
	err := ds.Get(fs.c, perfKey, &fo)
	if err == ds.ErrNoSuchEntity {
		return fo, err
	} else if err != nil {
		logif.E(err)
	}
	return fo, err
}
