package aefs_sr

import (
	"os"
	"sort"
	"sync/atomic"

	"appengine/datastore"

	"github.com/pbberlin/tools/logif"
	"github.com/pbberlin/tools/os/fsi"
)
import (
	pth "path"
	"path/filepath"
)

func (fs AeFileSys) String() string { return "aefs_sr" }

func (fs AeFileSys) Name() string { return fs.mount }

//---------------------------------------

// Create opens for read-write.
// Open opens for readonly access.
func (fs *AeFileSys) Create(name string) (fsi.File, error) {

	name = cleanseLeadingSlash(name)
	f := AeFile{}
	f.BName = pth.Base(name)

	err := fs.SaveFile(&f, name)
	if err != nil {
		return nil, err
	}

	// return &f, nil
	ff := fsi.File(&f)
	return ff, nil

}

// No distinction between Stat (links are followed)
// and LStat (links go unresolved)
// We don't support links yet, anyway
func (fs *AeFileSys) Lstat(path string) (os.FileInfo, error) {
	return fs.Stat(path)
}

// Strangely, neither MkdirAll nor Mkdir seem to have
// any concept of current working directory.
// They seem to operate relative to root.
func (fs *AeFileSys) Mkdir(name string, perm os.FileMode) error {
	_, err := fs.SaveDirByPath(name)
	return err
}

func (fs *AeFileSys) MkdirAll(path string, perm os.FileMode) error {
	_, err := fs.SaveDirByPath(path)
	return err
}

// Open opens for readonly access.
// Create opens for read-write.

// We could make provisions to ensure exclusive access;

// complies  with   os.Open()
// conflicts with  vfs.Open() signature
// conflicts with file.Open() interface of Afero
func (fs *AeFileSys) Open(name string) (fsi.File, error) {

	name = cleanseLeadingSlash(name)
	f, err := fs.GetFile(name)
	if err != nil {
		return nil, err
	}

	atomic.StoreInt64(&f.at, 0) // why is this not nested into f.Lock()-f.Unlock()?

	if f.closed == false { // already open
		// return ErrFileInUse // instead of waiting for lock?
	}

	f.Lock()
	f.closed = false
	f.Unlock()

	// return &f, nil
	ff := fsi.File(&f)
	return ff, nil
}

func (fs *AeFileSys) OpenFile(name string, flag int, perm os.FileMode) (fsi.File, error) {
	return fs.Open(name)
}

// ReadDir satisfies the vfs interface
// and ioutil.ReadDir.
func (fs *AeFileSys) ReadDir(path string) ([]os.FileInfo, error) {

	path = cleanseLeadingSlash(path)

	var dirs []AeDir
	var fis []os.FileInfo

	dir, err := fs.GetDirByPath(path)
	if path != dir.Dir+dir.BName {
		// panic(spf("path %v must equal dir and base %v %v ", path, dir.Dir, dir.BName))
	}
	logif.Pf("%15v => %24v", path, "")

	if err == datastore.ErrNoSuchEntity {
		return fis, err
	} else if err != nil {
		logif.E(err)
		return fis, err
	}

	q := datastore.NewQuery(tdir).Ancestor(dir.Key)
	keys, err := q.GetAll(fs.Ctx(), &dirs)
	_ = keys
	if err != nil {
		fs.Ctx().Errorf("Error fetching dir children of %v => %v", dir.Key, err)
		return fis, err
	}

	for i, v := range dirs {
		pK := keys[i].Parent()
		if pK != nil && !pK.Equal(dir.Key) {
			logif.Pf("%15v =>    skp %-17v", "", v.Dir+v.BName)
			continue
		}
		logif.Pf("%15v => %-24v", "", v.Dir+v.BName)
		fi := os.FileInfo(v)
		fis = append(fis, fi)
	}

	sort.Sort(FileInfoByName(fis))

	files, err := fs.GetFiles(path)
	if err != nil {
		return nil, err
	}
	for _, v := range files {
		fis = append(fis, os.FileInfo(v))
	}

	return fis, err

}

func (fs *AeFileSys) Remove(name string) error {
	panic(spf("Remove not (yet) implemented for %v", fs))
	return nil
}

func (fs *AeFileSys) RemoveAll(path string) error {

	paths := []string{}
	walkRemove := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			paths = append(paths, path)
		}
		// logif.Pf("Visited: %s %s \n", tp, path)
		return nil
	}

	err := filepath.Walk(path, walkRemove)

	logif.Pf("filepath.Walk() returned %v\n", err)

	for i := 0; i < len(paths); i++ {
		// todo: remove files
		// bottom-up remove dirs
	}

	return nil
}

func (fs *AeFileSys) Rename(oldname, newname string) error {
	panic(spf("Rename not (yet) implemented for %v", fs))
	return nil
}

func (fs *AeFileSys) Stat(path string) (os.FileInfo, error) {
	f, err := fs.GetFile(path)
	if err != nil {
		dir, err := fs.GetDirByPath(path)
		if err != nil {
			return nil, err
		}
		return os.FileInfo(dir), nil
	} else {
		return os.FileInfo(f), nil
	}
}

func (fs *AeFileSys) ReadFile(path string) ([]byte, error) {

	file, err := fs.GetFile(path)
	if err != nil {
		return []byte{}, err
	}
	return file.Data, err
}

func (fs *AeFileSys) WriteFile(name string, data []byte, perm os.FileMode) error {

	name = cleanseLeadingSlash(name)

	f, err := fs.Create(name)
	_ = f
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	ff := f.(*AeFile)

	err = fs.SaveFile(ff, name)
	if err != nil {
		return err
	}

	return err
}
