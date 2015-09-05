package dsfs

import (
	"os"
	"time"

	"github.com/pbberlin/tools/os/fsi/common"

	ds "appengine/datastore"
)

// Retrieves a directory in one go.
// Also used to check existence; returning ds.ErrNoSuchEntity
func (fs *dsFileSys) dirByPath(name string) (DsDir, error) {

	dir, bname := fs.SplitX(name)

	fo := DsDir{}
	fo.fSys = fs

	preciseK := ds.NewKey(fs.c, tdir, dir+common.Filify(bname), 0, nil)
	fo.Key = preciseK

	err := fo.MemCacheGet(dir + common.Filify(bname))
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
		fs.Ctx().Errorf("Error getting dir %v => %v", dir+common.Filify(bname), err)
	}

	fo.MemCacheSet()
	return fo, err
}

//
// dirsByPath might not find recently added directories.
// Upon finding nothing, it therefore returns the
// "warning" fsi.EmptyQueryResult
//
// It is currently used by ReadDir
func (fs *dsFileSys) dirsByPath(name string) ([]os.FileInfo, error) {

	dir, bname := fs.SplitX(name)

	var fis []os.FileInfo

	dirs, err := fs.SubtreeByPath(dir+bname, true)
	for _, v := range dirs {
		// log.Printf("%15v => %-24v", "", v.Dir+v.BName)
		fi := os.FileInfo(v)
		fis = append(fis, fi)
	}

	fs.dirsorter(fis)

	return fis, err

}

func (fs *dsFileSys) saveDirByPath(name string) (DsDir, error) {
	dir := DsDir{}
	return fs.saveDirByPathExt(dir, name)
}

func (fs *dsFileSys) saveDirByPathExt(dirObj DsDir, name string) (DsDir, error) {

	fo := DsDir{}
	fo.isDir = true
	fo.MModTime = dirObj.MModTime

	if dirObj.MMode == 0 {
		fo.MMode = 0755
	} else {
		fo.MMode = dirObj.MMode
	}

	if dirObj.MModTime.IsZero() {
		fo.MModTime = time.Now()
	} else {
		fo.MModTime = dirObj.MModTime
	}

	fo.fSys = fs

	dir, bname := fs.SplitX(name)
	fo.Dir = dir
	fo.BName = common.Filify(bname)

	preciseK := ds.NewKey(fs.c, tdir, dir+common.Filify(bname), 0, nil)
	fo.Key = preciseK

	// fs.Ctx().Infof("Saving dir %-14q  %q  %v ", fo.Dir, fo.BName, fo.Key)

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
