package gaefs

import (
	"fmt"
	"time"

	"github.com/pbberlin/tools/logif"

	"appengine/datastore"
)

func (fs *FileSys) getDirByExactKey(exactKey *datastore.Key) (Directory, error) {
	fo := Directory{}
	fo.Fs = fs
	fo.Key = exactKey
	err := datastore.Get(fs.c, exactKey, &fo)
	if err == datastore.ErrNoSuchEntity {
		return fo, err
	} else if err != nil {
		s := fmt.Sprintf("%v", exactKey)
		logif.E(err, s)
	}
	return fo, err
}

func (fs *FileSys) getDirUnderParent(parKey *datastore.Key, childName string) (Directory, error) {
	childKey := datastore.NewKey(fs.Ctx(), tdir, childName, 0, parKey)
	return fs.getDirByExactKey(childKey)
}

// The nested approach requires recursing directories.
// Retrieval is then possible via recurring dirByPathRecursive()
// or via GetDirByPathQuery()
func (fs *FileSys) saveDirUnderParent(name string, parent *datastore.Key) (Directory, error) {

	fo := Directory{}
	fo.IsDirectory = true
	fo.BName = name

	fo.Dir = dirFromKey(parent)
	fo.Mod = time.Now()
	fo.Fs = fs

	suggKey := datastore.NewKey(fs.Ctx(), tdir, name, 0, parent)
	fo.Key = suggKey
	fo.SKey = spf("%v", suggKey)

	effKey, err := datastore.Put(fs.Ctx(), suggKey, &fo)

	if err != nil {
		logif.E(err)
		return fo, err
	}
	if !suggKey.Equal(effKey) {
		fs.Ctx().Errorf("keys unequal %v - %v", suggKey, effKey)
	}

	fo.MemCacheSet()

	return fo, nil
}
