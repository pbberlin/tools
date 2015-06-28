package filesys

import (
	"time"

	"github.com/pbberlin/tools/logif"

	ds "appengine/datastore"
)

// saveDirByPathUnderRoot saves *all* directories under root.
// Retrieval is then possible directly via rootedGetDirByPath()
// or via GetDirByPathQuery().
// The main disadvantage is that "ancestor group updates"
// in google datastore are restricted to frequency "~ five per second".
// If you anticipate fewer directory changes, consider it.
func (fs *FileSys) rootedSaveDirByPath(path string) (Directory, error) {

	fo := Directory{}
	fo.IsDir = true
	fo.Dir = parent(path)
	fo.Name = fileN(path)
	fo.Mod = time.Now()
	fo.Fs = fs

	perfKey := ds.NewKey(fs.c, t, path, 0, fs.RootDir.Key)

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
// But only if it was saved with rootedSaveDirByPath.
func (fs FileSys) rootedGetDirByPath(path string) (Directory, error) {
	fo := Directory{}
	perfKey := ds.NewKey(fs.c, t, path, 0, fs.RootDir.Key)
	err := ds.Get(fs.c, perfKey, fo)
	if err == ds.ErrNoSuchEntity {
		return fo, err
	} else if err != nil {
		logif.E(err)
	}
	return fo, err
}
