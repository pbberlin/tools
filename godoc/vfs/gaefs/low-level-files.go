package gaefs

import (
	"fmt"
	pth "path"
	"sort"
	"strings"
	"time"

	"github.com/pbberlin/tools/logif"

	"appengine/datastore"
)

// similar to ReadDir but returning only files
func (fs *AeFileSys) GetFiles(path string) ([]AeFile, error) {

	path = cleanseLeadingSlash(path)

	var files []AeFile
	var filesDirect []AeFile

	dir, err := fs.GetDirByPath(path)
	if err == datastore.ErrNoSuchEntity {
		return files, err
	} else if err != nil {
		logif.E(err)
		return files, err
	}

	q := datastore.NewQuery(tfil).Ancestor(dir.Key)
	keys, err := q.GetAll(fs.Ctx(), &files)
	if err != nil {
		fs.Ctx().Errorf("Error fetching files children of %v => %v", dir.Key, err)
		return files, err
	}

	for i, v := range files {
		pK := keys[i].Parent()
		if pK != nil && !pK.Equal(dir.Key) {
			// logif.Pf("%15v =>    skp %-17v", "", v.Dir+v.BName)
			continue
		}
		// logif.Pf("%15v => %-24v", "", v.Dir+v.BName)
		filesDirect = append(filesDirect, v)
		// logif.Pf("%s", v.Data)
	}

	sort.Sort(AeFileByName(filesDirect))

	return filesDirect, err
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

	if !strings.HasPrefix(path, fs.RootDir()) {
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
