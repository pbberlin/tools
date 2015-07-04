package gaefs

import (
	"fmt"
	pth "path"
	"sort"
	"time"

	"github.com/pbberlin/tools/logif"

	"appengine/datastore"
)

// similar to ReadDir but returning only files
// Todo: Sort files by name
func (fs *AeFileSys) GetFiles(path string) ([]AeFile, error) {

	path = cleanseLeadingSlash(path)

	var files []AeFile

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

	sort.Sort(AeFileByName(files))

	return files, err
}

func (fs *AeFileSys) GetFile(path string) (AeFile, error) {

	path = cleanseLeadingSlash(path)

	fo := AeFile{}
	fo.Fs = fs

	sdir, base := pth.Split(path)

	var dir AeDir
	var err error
	if sdir == "" {
		dir = fs.rootDir
	} else {
		dir, err = fs.GetDirByPath(sdir)
		if err == datastore.ErrNoSuchEntity {
			return fo, err
		} else if err != nil {
			logif.E(err)
			return fo, err
		}
	}

	fileKey := datastore.NewKey(fs.Ctx(), tfil, base, 0, dir.Key)
	fo.Key = fileKey
	err = datastore.Get(fs.c, fileKey, &fo)
	if err == datastore.ErrNoSuchEntity {
		return fo, err
	} else if err != nil {
		s := fmt.Sprintf("%v", fileKey)
		logif.E(err, s)
	}

	return fo, err

}

//
// The nested approach requires recursing directories.
// Retrieval is then possible via recurring dirByPathRecursive()
// or via GetDirByPathQuery()
func (fs *AeFileSys) SaveFile(f *AeFile, path string) error {

	path = cleanseLeadingSlash(path)

	if f.BName == "" {
		return fmt.Errorf("file needs name")
	}

	f.Fs = fs
	f.Dir = path
	f.MModTime = time.Now()

	dir, err := fs.GetDirByPath(path)
	if err == datastore.ErrNoSuchEntity {
		return err
	} else if err != nil {
		logif.E(err)
		return err
	}

	suggKey := datastore.NewKey(fs.Ctx(), tfil, f.BName, 0, dir.Key)
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

	// f.MemCacheSet()

	return nil
}
