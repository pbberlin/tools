package dsfs

import (
	"os"
	"sort"
	"time"

	ds "appengine/datastore"
)

// Retrieves a directory in one go.
// Also used to check existence; returning ds.ErrNoSuchEntity
func (fs *dsFileSys) dirByPath(name string) (AeDir, error) {

	dir, bname := fs.SplitX(name)

	fo := AeDir{}
	fo.fSys = fs

	preciseK := ds.NewKey(fs.c, tdir, dir+bname, 0, nil)
	fo.Key = preciseK

	err := fo.MemCacheGet(dir + bname)
	if err == nil {
		fo.fSys = fs
		fo.Key = preciseK
		// log.Printf("  mcg %-16v %v", dir+bname, fo.Key)
		return fo, nil
	}

	err = ds.Get(fs.c, preciseK, &fo)
	if err == ds.ErrNoSuchEntity {
		// uncomment to see where directories do not exist:
		// log.Printf("no directory: %-20v ", path)
		// runtimepb.StackTrace(4)
		return fo, err
	} else if err != nil {
		fs.Ctx().Errorf("Error getting dir %v => %v", dir+bname, err)
	}

	fo.MemCacheSet()
	return fo, err
}

// dirsByPath might not find recently added directories.
// Upon finding nothing, it therefore returns the
// "warning" fsi.EmptyQueryResult
func (fs *dsFileSys) dirsByPath(name string) ([]os.FileInfo, error) {

	dir, bname := fs.SplitX(name)

	var fis []os.FileInfo

	dirs, err := fs.SubdirsByPath(dir+bname, true)
	for _, v := range dirs {
		// log.Printf("%15v => %-24v", "", v.Dir+v.BName)
		fi := os.FileInfo(v)
		fis = append(fis, fi)
	}

	sort.Sort(FileInfoByName(fis))

	return fis, err

}

func (fs *dsFileSys) saveDirByPath(name string) (AeDir, error) {

	fo := AeDir{}
	fo.isDir = true
	fo.MModTime = time.Now()
	fo.fSys = fs

	dir, bname := fs.SplitX(name)
	fo.Dir = dir
	fo.BName = bname

	preciseK := ds.NewKey(fs.c, tdir, dir+bname, 0, nil)
	fo.Key = preciseK

	effKey, err := ds.Put(fs.c, preciseK, &fo)
	if err != nil {
		fs.Ctx().Errorf("Error saving dir %v => %v", dir+bname, err)
		return fo, err
	}

	if !preciseK.Equal(effKey) {
		fs.Ctx().Errorf("dir keys unequal %v - %v", preciseK, effKey)
	}

	fo.MemCacheSet()

	// recurse upwards
	_, err = fs.dirByPath(fo.Dir)
	if err == ds.ErrNoSuchEntity {
		_, err = fs.saveDirByPath(fo.Dir)
		if err != nil {
			return fo, err
		}
	}

	return fo, nil
}
