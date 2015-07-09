package aefs

import (
	"os"
	"sort"
	"strings"
	"time"

	"github.com/pbberlin/tools/logif"

	pth "path"

	ds "appengine/datastore"
)

// Retrieves a directory in one go.
func (fs *AeFileSys) dirByPath(path string) (AeDir, error) {

	path = cleanseLeadingSlash(path)

	fo := AeDir{}
	fo.Fs = fs
	preciseK := ds.NewKey(fs.c, tdir, fs.RootDir()+path, 0, nil)
	fo.Key = preciseK
	err := ds.Get(fs.c, preciseK, &fo)
	if err == ds.ErrNoSuchEntity {
		return fo, err
	} else if err != nil {
		logif.E(err)
	}
	return fo, err
}

func (fs *AeFileSys) dirsByPath(path string) ([]os.FileInfo, error) {

	path = cleanseLeadingSlash(path)

	var fis []os.FileInfo

	dirs, err := fs.subdirsByPath(path, true)

	for _, v := range dirs {
		// logif.Pf("%15v => %-24v", "", v.Dir+v.BName)
		fi := os.FileInfo(v)
		fis = append(fis, fi)
	}

	sort.Sort(FileInfoByName(fis))

	return fis, err

}

func (fs *AeFileSys) saveDirByPath(path string) (AeDir, error) {

	path = cleanseLeadingSlash(path)

	fo := AeDir{}
	fo.isDir = true
	dir, base := pth.Split(path)
	fo.Dir = dir
	if !strings.HasPrefix(fo.Dir, fs.RootDir()) {
		fo.Dir = fs.RootDir() + fo.Dir
	}

	fo.BName = base
	fo.MModTime = time.Now()
	fo.Fs = fs

	preciseK := ds.NewKey(fs.c, tdir, path, 0, nil)

	fo.Key = preciseK
	fo.SKey = spf("%v", preciseK) // not effKey.Encode()

	effKey, err := ds.Put(fs.c, preciseK, &fo)

	if err != nil {
		logif.E(err)
		return fo, err
	}

	if !preciseK.Equal(effKey) {
		fs.Ctx().Errorf("keys unequal %v - %v", preciseK, effKey)
	}

	fo.MemCacheSet()

	return fo, nil
}
