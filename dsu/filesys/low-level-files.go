package filesys

import (
	"github.com/pbberlin/tools/logif"

	ds "appengine/datastore"
)

func (fs *FileSys) GetFiles(path string) ([]File, error) {

	var files []File

	dir, err := fs.GetDirByPath(path)
	if err == ds.ErrNoSuchEntity {
		return files, err
	} else if err != nil {
		logif.E(err)
		return files, err
	}

	q := ds.NewQuery(tfil).Ancestor(dir.Key)
	keys, err := q.GetAll(fs.Ctx(), &files)
	_ = keys
	if err != nil {
		fs.Ctx().Errorf("Error getching files children of %v => %v", dir.Key, err)
		return files, err
	}

	return files, err
}
