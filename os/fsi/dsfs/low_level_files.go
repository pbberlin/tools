package dsfs

import (
	"fmt"
	"strings"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/common"
	"github.com/pbberlin/tools/runtimepb"

	"appengine/datastore"
)

func (fs *dsFileSys) fileByPath(name string) (DsFile, error) {

	dir, bname := fs.SplitX(name)

	fo := DsFile{}
	fo.fSys = fs

	//
	if dir == fs.RootDir() && bname == "" {
		return fo, fsi.ErrRootDirNoFile
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
func (fs *dsFileSys) filesByPath(name string) ([]DsFile, error) {

	dir, bname := fs.SplitX(name)

	var files []DsFile

	foDir, err := fs.dirByPath(dir + common.Filify(bname))
	if err == datastore.ErrNoSuchEntity {
		return files, err
	} else if err != nil {
		fs.Ctx().Errorf("Error reading dir for files %v => %v", dir+bname, err)
		return files, err
	}

	// fs.Ctx().Infof("  Files by Path %-20v - got dir %-10v - %v", dir+bname, foDir.Name(), foDir.Key)

	q := datastore.NewQuery(tfil).Ancestor(foDir.Key)
	keys, err := q.GetAll(fs.Ctx(), &files)
	if err != nil {
		fs.Ctx().Errorf("Error fetching files children of %v => %v", foDir.Key, err)
		return files, err
	}

	// fs.Ctx().Infof("  Files by Path %-20v - got files %v", dir+bname, len(files))

	for i := 0; i < len(files); i++ {
		files[i].Key = keys[i]
		files[i].fSys = fs
	}

	fs.filesorter(files)

	return files, err
}

//
//
// Path is the directory, BName contains the base name.
func (fs *dsFileSys) saveFileByPath(f *DsFile, name string) error {

	dir, bname := fs.SplitX(name)
	f.Dir = dir
	// bname was only submitted in the fileobject
	// correct previous
	if f.BName != "" && f.BName != bname {
		dir = dir + bname // re-append
		if !strings.HasSuffix(dir, sep) {
			dir += sep
		}
		bname = f.BName
		f.Dir = dir
	}
	f.BName = common.Filify(bname)

	// f.MModTime = time.Now()
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
		fs.Ctx().Errorf("file keys unequal %v - %v; %v %s", suggKey, effKey, f.Dir+f.BName, f.Data)
		runtimepb.StackTrace(6)
	}

	// f.MemCacheSet()

	return nil
}
