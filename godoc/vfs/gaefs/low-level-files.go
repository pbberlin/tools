package gaefs

import (
	"fmt"

	"github.com/pbberlin/tools/logif"

	"appengine/datastore"
)

func (fs *FileSys) GetFiles(path string) ([]File, error) {

	var files []File

	dir, err := fs.GetDirByPath(path)
	if err == datastore.ErrNoSuchEntity {
		return files, err
	} else if err != nil {
		logif.E(err)
		return files, err
	}

	q := datastore.NewQuery(tfil).Ancestor(dir.Key)
	keys, err := q.GetAll(fs.Ctx(), &files)
	_ = keys
	if err != nil {
		fs.Ctx().Errorf("Error getching files children of %v => %v", dir.Key, err)
		return files, err
	}

	return files, err
}

// The nested approach requires recursing directories.
// Retrieval is then possible via recurring dirByPathRecursive()
// or via GetDirByPathQuery()
func (fs *FileSys) SaveFile(f *File, path string) error {

	if f.Name == "" {
		return fmt.Errorf("file needs name")
	}

	f.Fs = fs
	f.Dir = path

	dir, err := fs.GetDirByPath(path)
	if err == datastore.ErrNoSuchEntity {
		return err
	} else if err != nil {
		logif.E(err)
		return err
	}

	suggKey := datastore.NewKey(fs.Ctx(), tfil, f.Name, 0, dir.Key)
	f.Key = suggKey
	f.SKey = spf("%v", suggKey)

	effKey, err := datastore.Put(fs.Ctx(), suggKey, f)
	if err != nil {
		logif.E(err)
		return err
	}
	if !suggKey.Equal(effKey) {
		fs.Ctx().Errorf("keys unequal %v - %v", suggKey, effKey)
	}

	f.MemCacheSet()

	return nil
}
