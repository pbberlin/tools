package aefs

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"appengine/datastore"

	"github.com/pbberlin/tools/logif"
	"github.com/pbberlin/tools/os/fsi"
)
import (
	pth "path"
	"path/filepath"
)

func (fs *AeFileSys) Chmod(name string, mode os.FileMode) error {
	panic(spf("Chmod not (yet) implemented for %v", fs))
	return nil
}

func (fs *AeFileSys) Chtimes(name string, atime time.Time, mtime time.Time) error {
	panic(spf("Chtimes not (yet) implemented for %v", fs))
	return nil
}

// Create opens for read-write.
// Open opens for readonly access.
func (fs *AeFileSys) Create(name string) (fsi.File, error) {

	if name == "" {
		err := fmt.Errorf("name cant be empty string")
		logif.E(err)
		return nil, err
	}

	name = cleanseLeadingSlash(name)
	f := AeFile{}
	f.BName = pth.Base(name)

	err := fs.saveFileByPath(&f, name)
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
	_, err := fs.saveDirByPath(name)
	return err
}

func (fs *AeFileSys) MkdirAll(path string, perm os.FileMode) error {
	_, err := fs.saveDirByPath(path)
	return err
}

func (fs AeFileSys) String() string {
	return "gaefs"
}

func (fs AeFileSys) Name() string {
	return fs.String()
}

// Open opens for readonly access.
// Create opens for read-write.

// We could make provisions to ensure exclusive access;

// complies  with   os.Open()
// conflicts with  vfs.Open() signature
// conflicts with file.Open() interface of Afero
func (fs *AeFileSys) Open(name string) (fsi.File, error) {

	name = cleanseLeadingSlash(name)
	f, err := fs.fileByPath(name)
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
// It is similar to filesByPaths, but returning only dirs
func (fs *AeFileSys) ReadDir(path string) ([]os.FileInfo, error) {
	return fs.dirsByPath(path)
}

func (fs *AeFileSys) Readdirnames(path string) (names []string, err error) {
	fis, err := fs.ReadDir(path)
	names = make([]string, 0, len(fis))
	for _, lp := range fis {
		names = append(names, lp.Name())
	}
	return names, err
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
	f, err := fs.fileByPath(path)
	if err != nil {
		dir, err := fs.dirByPath(path)
		if err != nil {
			return nil, err
		}
		return os.FileInfo(dir), nil
	} else {
		return os.FileInfo(f), nil
	}
}

func (fs *AeFileSys) ReadFile(path string) ([]byte, error) {

	file, err := fs.fileByPath(path)
	if err != nil {
		return []byte{}, err
	}
	return file.Data, err
}

func (fs *AeFileSys) WriteFile(name string, data []byte, perm os.FileMode) error {

	name = cleanseLeadingSlash(name)

	f, err := fs.Create(name)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	// ff := f.(*AeFile)

	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}

// DeleteAll deletes across all roots
// DeleteAll deletes by kind alone.
func (fs *AeFileSys) DeleteAll() (string, error) {

	msg := ""
	{
		q := datastore.NewQuery(tfil).KeysOnly()
		var files []AeFile
		keys, err := q.GetAll(fs.Ctx(), &files)
		if err != nil {
			msg += "could not get file keys\n"
			return msg, err
		}

		err = datastore.DeleteMulti(fs.Ctx(), keys)
		if err != nil {
			msg += "error deleting files\n"
			return msg, err
		}

		msg += spf("%v files deleted\n", len(keys))

	}

	{
		q := datastore.NewQuery(tdir).KeysOnly()
		var dirs []AeDir
		keys, err := q.GetAll(fs.Ctx(), &dirs)
		if err != nil {
			msg += "could not get dir keys\n"
			return msg, err
		}

		err = datastore.DeleteMulti(fs.Ctx(), keys)
		if err != nil {
			msg += "error deleting directories\n"
			return msg, err
		}

		msg += spf("%v directories deleted\n", len(keys))
	}

	return msg, nil
}
