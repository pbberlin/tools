package rooted

import (
	"time"

	"github.com/pbberlin/tools/logif"

	ds "appengine/datastore"
)

type Rooted FileSys

func (fs Rooted) saveDirByPath(path string) Directory {

	fo := Directory{}
	fo.IsDir = true
	fo.Dir = parent(path)
	fo.Name = fileN(path)
	fo.Mod = time.Now()
	fo.fs = &fs

	perfKey := ds.NewKey(fs.c, t, path, 0, fs.RootDir.Key)
	fo.Key = perfKey
	// fo.SKey = spf("%v", perfKey) // not effKey.Encode()

	effKey, err := ds.Put(fs.c, perfKey, &fo)

	if err != nil {
		Fs.c.Errorf("%v  [%s]", err, fo.Key)
	} else {
		// Fs.c.Infof("fso ds-saved - key %v [%v]", effKey, effKey.Encode())
	}

	if !perfKey.Equal(effKey) {
		Fs.c.Errorf("keys unequal %v - %v", perfKey, effKey)
	}

	// fo.Key = effKey
	// fo.SKey = spf("%v", effKey) // not effKey.Encode()

	fo.MemCacheSet()

	return fo
}

func (fs Rooted) dirByPath(path string) (Directory, error) {

	fo := new(Directory)
	perfKey := ds.NewKey(fs.c, t, path, 0, fs.RootDir.Key)
	err := ds.Get(fs.c, perfKey, fo)
	if err == ds.ErrNoSuchEntity {
		return *fo, err
	} else if err != nil {
		logif.E(err)
	}

	return *fo, err
}
