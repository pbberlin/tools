package filesys

import (
	"fmt"
	"time"

	"github.com/pbberlin/tools/logif"

	ds "appengine/datastore"
)

func (fs FileSys) saveFsoByPath(path string, isDir bool) FSysObj {

	fo := FSysObj{}
	fo.IsDir = isDir
	fo.Dir = parent(path)
	fo.Name = fileN(path)
	fo.Mod = time.Now()
	fo.fs = &fs

	perfKey := ds.NewKey(fs.c, t, path, 0, fs.RootDir.Key)
	fo.Key = perfKey
	// fo.SKey = spf("%v", perfKey) // not effKey.Encode()

	effKey, err := ds.Put(fs.c, perfKey, &fo)

	if err != nil {
		fs.c.Errorf("%v  [%s]", err, fo.Key)
	} else {
		// fs.c.Infof("fso ds-saved - key %v [%v]", effKey, effKey.Encode())
	}

	if !perfKey.Equal(effKey) {
		fs.c.Errorf("keys unequal %v - %v", perfKey, effKey)
	}

	// fo.Key = effKey
	// fo.SKey = spf("%v", effKey) // not effKey.Encode()

	fo.MemCacheSet()

	return fo
}

func (fs FileSys) getFsoByPath(path string) (FSysObj, error) {

	fo := new(FSysObj)
	perfKey := ds.NewKey(fs.c, t, path, 0, fs.RootDir.Key)
	err := ds.Get(fs.c, perfKey, fo)
	logif.E(err)
	if fo == nil {
		err = fmt.Errorf("No fso for path %v", path)
		logif.E(err)
	}

	return *fo, err
}
