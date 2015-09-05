package dsfs

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"appengine/datastore"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/common"
)

func (fs dsFileSys) Name() string { return "dsfs" }

func (fs dsFileSys) String() string { return fs.mount }

//---------------------------------------

func (fs *dsFileSys) Chmod(name string, mode os.FileMode) error {

	f, err := fs.fileByPath(name)
	if err == nil {
		f.MMode = mode

		err := f.Sync()
		if err != nil {
			return err
		}
		return nil
	} else {
		dir, err := fs.dirByPath(name)
		if err != nil {
			return err
		}
		dir.MMode = mode
		_, err = fs.saveDirByPathExt(dir, dir.Dir+dir.BName)
		if err != nil {
			return err
		}

	}

	return nil
}

func (fs *dsFileSys) Chtimes(name string, atime time.Time, mtime time.Time) error {

	f, err := fs.fileByPath(name)
	if err == nil {
		f.MModTime = atime
		err := f.Sync()
		if err != nil {
			return err
		}
		return nil
	} else {
		dir, err := fs.dirByPath(name)
		if err != nil {
			return err
		}
		dir.MModTime = atime
		_, err = fs.saveDirByPathExt(dir, dir.Dir+dir.BName)
		if err != nil {
			return err
		}
	}

	return nil

}

// Create opens for read-write.
// Open opens for readonly access.
func (fs *dsFileSys) Create(name string) (fsi.File, error) {

	// WriteFile & Create
	dir, bname := fs.SplitX(name)

	f := DsFile{}
	f.fSys = fs
	f.BName = common.Filify(bname)
	f.Dir = dir
	f.MModTime = time.Now()
	f.MMode = 0644

	// let all the properties by set by fs.saveFileByPath
	err := f.Sync()
	if err != nil {
		return nil, err
	}

	// return &f, nil
	ff := fsi.File(&f)
	return ff, err

}

// No distinction between Stat (links are followed)
// and LStat (links go unresolved)
// We don't support links yet, anyway
func (fs *dsFileSys) Lstat(path string) (os.FileInfo, error) {
	fi, err := fs.Stat(path)
	return fi, err
}

// Strangely, neither MkdirAll nor Mkdir seem to have
// any concept of current working directory.
// They seem to operate relative to root.
func (fs *dsFileSys) Mkdir(name string, perm os.FileMode) error {
	_, err := fs.saveDirByPath(name)
	return err
}

func (fs *dsFileSys) MkdirAll(path string, perm os.FileMode) error {
	_, err := fs.saveDirByPath(path)
	return err
}

// Open() open existing file for readonly access.
// Create() should be used   for read-write.

// We could make provisions to ensure exclusive access;

// complies  with   os.Open()
// conflicts with  vfs.Open() signature
// conflicts with file.Open() interface of Afero
func (fs *dsFileSys) Open(name string) (fsi.File, error) {

	// explicitly requesting directory?
	_, bname := fs.SplitX(name)
	if strings.HasSuffix(bname, "/") {
		dir, err := fs.dirByPath(name)
		if err == nil {
			ff := fsi.File(&dir)
			return ff, nil
		}
	}

	// otherwise: try file, then directory
	f, err := fs.fileByPath(name)

	if err != nil && err != datastore.ErrNoSuchEntity && err != fsi.ErrRootDirNoFile {
		return nil, err
	}
	if err == datastore.ErrNoSuchEntity || err == fsi.ErrRootDirNoFile {
		// http.FileServer requires, that
		// we return a directory here.
		// It's also compliant to os.Open(),
		// where "os.File" means directories too.
		dir, err2 := fs.dirByPath(name)
		if err2 != nil {
			return nil, err
		}
		ff := fsi.File(&dir)
		return ff, nil
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

func (fs *dsFileSys) OpenFile(name string, flag int, perm os.FileMode) (fsi.File, error) {
	return fs.Open(name)
}

// See fsi.FileSystem interface.

//
// ReadDir might not find recently added directories.
func (fs *dsFileSys) ReadDir(name string) ([]os.FileInfo, error) {

	dirs, err := fs.dirsByPath(name)
	// fs.Ctx().Infof("dsfs readdir %-20v dirs %v", name, len(dirs))
	if err != nil && err != fsi.EmptyQueryResult {
		return nil, err
	}
	fs.dirsorter(dirs)

	files, err := fs.filesByPath(name)
	// fs.Ctx().Infof("dsfs readdir %-20v fils %v %v", name, len(files), err)
	if err != nil {
		return nil, err
	}
	fs.filesorter(files)

	for _, v := range files {
		dirs = append(dirs, os.FileInfo(v))
	}
	return dirs, nil
}

func (fs *dsFileSys) Remove(name string) error {

	// log.Printf("trying to remove %-20v", name)
	f, err := fs.fileByPath(name)
	if err == nil {
		// log.Printf("   found file %v", f.Dir+f.BName)
		// log.Printf("   fkey %-26v", f.Key)
		err = datastore.Delete(fs.Ctx(), f.Key)
		if err != nil {
			return fmt.Errorf("error removing file %v", err)
		}
	}

	d, err := fs.dirByPath(name)
	if err == nil {
		// log.Printf("   dkey %v", d.Key)
		err = datastore.Delete(fs.Ctx(), d.Key)
		d.MemCacheDelete()
		if err != nil {
			return fmt.Errorf("error removing dir %v", err)
		}
	}

	return nil

}

func (fs *dsFileSys) RemoveAll(path string) error {

	paths := []string{}
	walkRemove := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			// do nothing; don't break the walk
			fs.Ctx().Errorf("Error walking %v => %v", path, err)
		} else {
			if f != nil { // && f.IsDir() to constrain
				paths = append(paths, path)
			}
		}
		return nil
	}

	err := common.Walk(fs, path, walkRemove)
	if err != nil {
		fs.Ctx().Errorf("Error removing %v => %v", path, err)
	}

	// Walk crawls directories first, files second.
	// Intuitively removal in reverse order should always work. Or does it not?
	for i := 0; i < len(paths); i++ {
		iRev := len(paths) - 1 - i
		err := fs.Remove(paths[iRev])
		if err != nil {
			fs.Ctx().Errorf("Error removing %v => %v", paths[iRev], err)
			return err
		}
		fs.Ctx().Infof("removed path %v", paths[iRev])
	}

	return nil
}

func (fs *dsFileSys) Rename(oldname, newname string) error {
	// we could use a walk similar to remove all
	return fsi.NotImplemented
}

func (fs *dsFileSys) Stat(path string) (os.FileInfo, error) {

	f, err := fs.fileByPath(path)
	if err != nil && err != datastore.ErrNoSuchEntity && err != fsi.ErrRootDirNoFile {
		log.Fatalf("OTHER ERROR %v", err)

		return nil, err
	}
	if err == datastore.ErrNoSuchEntity || err == fsi.ErrRootDirNoFile {
		// log.Printf("isno file err %-24q =>  %v", path, err)
		dir, err := fs.dirByPath(path)
		if err != nil {
			return nil, err
		}
		fiDir := os.FileInfo(dir)
		// log.Printf("Stat for dire %-24q => %-24v, %v", path, fiDir.Name(), err)
		return fiDir, nil
	}

	fiFi := os.FileInfo(f)
	// log.Printf("Stat for file %-24q => %-24v, %v", path, fiFi.Name(), err)
	return fiFi, nil
}

func (fs *dsFileSys) ReadFile(path string) ([]byte, error) {

	file, err := fs.fileByPath(path)
	if err != nil {
		return []byte{}, err
	}
	return file.Data, err
}

// Only one save operation required
func (fs *dsFileSys) WriteFile(name string, data []byte, perm os.FileMode) error {

	// WriteFile & Create
	dir, bname := fs.SplitX(name)
	f := DsFile{}
	f.Dir = dir
	f.BName = common.Filify(bname)
	f.fSys = fs
	f.MModTime = time.Now()

	var err error
	_, err = f.Write(data)
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}
