package aefs

import (
	"fmt"
	pth "path"
	"sort"
	"strings"
	"time"

	"github.com/pbberlin/tools/logif"

	"appengine/datastore"
)

func (fs *AeFileSys) fileByPath(path string) (AeFile, error) {

	path = cleanseLeadingSlash(path)

	fo := AeFile{}
	fo.Fs = fs

	sdir, base := pth.Split(path)

	dir, err := fs.dirByPath(sdir)
	if err == datastore.ErrNoSuchEntity {
		return fo, err
	} else if err != nil {
		logif.E(err)
		return fo, err
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

// similar to ReadDir but returning only files
func (fs *AeFileSys) filesByPath(path string) ([]AeFile, error) {

	path = cleanseLeadingSlash(path)

	var files []AeFile

	dir, err := fs.dirByPath(path)
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
		fs.Ctx().Errorf("Error fetching files children of %v => %v", dir.Key, err)
		return files, err
	}

	sort.Sort(AeFileByName(files))

	return files, err
}

//
//
// Path is the directory, BName contains the base name.
func (fs *AeFileSys) saveFileByPath(f *AeFile, path string) error {

	path = cleanseLeadingSlash(path)

	if f.BName == "" {
		return fmt.Errorf("file needs name")
	}

	if !strings.HasPrefix(path, fs.RootName()) {
		path = fs.RootDir() + path
	}

	if strings.HasSuffix(path, f.BName) {
		path = path[:len(path)-len(f.BName)]
	}

	f.Dir = path

	if !strings.HasSuffix(f.Dir, sep) {
		f.Dir += sep
	}

	f.MModTime = time.Now()
	f.Fs = fs

	//
	// -------------now the datastore part-------------------------

	logif.Pf("%q %q", f.Dir, f.BName)

	dir, err := fs.dirByPath(f.Dir)
	if err == datastore.ErrNoSuchEntity {
		return err
	} else if err != nil {
		logif.E(err)
		return err
	}

	suggKey := datastore.NewKey(fs.Ctx(), tfil, f.BName, 0, dir.Key)
	f.Key = suggKey

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
