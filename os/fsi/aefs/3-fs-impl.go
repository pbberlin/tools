package aefs

import (
	"fmt"
	"os"
	"sync/atomic"

	"appengine/datastore"

	"github.com/pbberlin/tools/logif"
	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/fsc"
)
import pth "path"

func (fs AeFileSys) Name() string { return "aefs" }

func (fs AeFileSys) String() string { return fs.mount }

//---------------------------------------

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

// ReadDir satisfies the godoc.vfs interface
// and ioutil.ReadDir.
func (fs *AeFileSys) ReadDir(path string) ([]os.FileInfo, error) {
	dirs, err := fs.dirsByPath(path)
	if err != nil {
		return nil, err
	}
	files, err := fs.filesByPath(path)
	if err != nil {
		return nil, err
	}
	for _, v := range files {
		dirs = append(dirs, os.FileInfo(v))
	}
	return dirs, nil
}

func (fs *AeFileSys) Remove(name string) error {

	logif.Pf("trying to remove %-20v", name)

	f, err := fs.fileByPath(name)
	if err != nil {
		logif.Pf("   fkey %v", f.Key)
		err = datastore.Delete(fs.Ctx(), f.Key)
		if err != nil {
			return fmt.Errorf("error removing file %v", err)
		}
	} else {
		d, err := fs.dirByPath(name)
		if err != nil {
			logif.Pf("   dkey %v", d.Key)
			err = datastore.Delete(fs.Ctx(), d.Key)
			if err != nil {
				return fmt.Errorf("error removing dir %v", err)
			}
		}
	}
	return nil

}

func (fs *AeFileSys) RemoveAll(path string) error {

	paths := []string{}
	walkRemove := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() || true {
			paths = append(paths, path)
		}
		// logif.Pf("Visited: %s %s \n", tp, path)
		return nil
	}

	err := fsc.Walk(fs, path, walkRemove)
	logif.E(err)

	for i := 0; i < len(paths); i++ {
		iRev := len(paths) - 1 - i
		err := fs.Remove(paths[iRev])
		if err != nil {
			return err
		}
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
