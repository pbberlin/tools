package aefs

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"appengine/datastore"
)

func (fs *aeFileSys) fileByPath(name string) (AeFile, error) {

	dir, bname := fs.pathInternalize(name)

	fo := AeFile{}
	fo.fSys = fs

	//
	if dir == fs.RootDir() && bname == "" {
		return fo, fmt.Errorf("rootdir; no file")
	}

	foDir, err := fs.dirByPath(dir)
	if err == datastore.ErrNoSuchEntity {
		return fo, err
	} else if err != nil {
		fs.Ctx().Errorf("Error reading dir for file %v => %v", dir+bname, err)
		return fo, err
	}

	fileKey := datastore.NewKey(fs.Ctx(), tfil, bname, 0, foDir.Key)
	fo.Key = fileKey
	err = datastore.Get(fs.c, fileKey, &fo)
	if err == datastore.ErrNoSuchEntity {
		return fo, err
	} else if err != nil {
		s := fmt.Sprintf("%v", fileKey)
		fs.Ctx().Errorf("Error reading file %v (%v) => %v", dir+bname, s, err)
	}

	return fo, err

}

// similar to ReadDir but returning only files
func (fs *aeFileSys) filesByPath(name string) ([]AeFile, error) {

	dir, bname := fs.pathInternalize(name)

	var files []AeFile

	foDir, err := fs.dirByPath(dir + bname)
	if err == datastore.ErrNoSuchEntity {
		return files, err
	} else if err != nil {
		fs.Ctx().Errorf("Error reading dir for files %v => %v", dir+bname, err)
		return files, err
	}

	q := datastore.NewQuery(tfil).Ancestor(foDir.Key)
	keys, err := q.GetAll(fs.Ctx(), &files)
	if err != nil {
		fs.Ctx().Errorf("Error fetching files children of %v => %v", foDir.Key, err)
		return files, err
	}

	for i := 0; i < len(files); i++ {
		files[i].Key = keys[i]
		files[i].fSys = fs
	}

	sort.Sort(AeFileByName(files))

	return files, err
}

//
//
// Path is the directory, BName contains the base name.
func (fs *aeFileSys) saveFileByPath(f *AeFile, name string) error {

	dir, bname := fs.pathInternalize(name)
	f.Dir = dir
	// bname was only submitted in the fileobject only
	// correct previous
	if f.BName != "" && f.BName != bname {
		dir = dir + bname // re-append
		if !strings.HasSuffix(dir, sep) {
			dir += sep
		}
		bname = f.BName
		f.Dir = dir
	}
	f.BName = bname

	f.MModTime = time.Now()
	f.fSys = fs

	//
	// -------------now the datastore part-------------------------

	foDir, err := fs.dirByPath(dir)
	if err == datastore.ErrNoSuchEntity {
		return err
	} else if err != nil {
		fs.Ctx().Errorf("Error reading dir for file %v => %v", dir+bname, err)
		return err
	}

	suggKey := datastore.NewKey(fs.Ctx(), tfil, f.BName, 0, foDir.Key)
	f.Key = suggKey

	effKey, err := datastore.Put(fs.Ctx(), suggKey, f)
	if err != nil {
		fs.Ctx().Errorf("Error saving file %v => %v", dir+bname, err)
		return err
	}
	if !suggKey.Equal(effKey) {
		fs.Ctx().Errorf("keys unequal %v - %v", suggKey, effKey)
	}

	// f.MemCacheSet()

	return nil
}
