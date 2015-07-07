package gaefs

import (
	"os"
	"time"

	"github.com/pbberlin/tools/logif"
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

func (fs *AeFileSys) Create(name string) (AeFile, error) {

	name = cleanseLeadingSlash(name)
	f := AeFile{}
	f.BName = pth.Base(name)

	err := fs.SaveFile(&f, name)
	if err != nil {
		return f, err
	}
	return f, err
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

func (fs AeFileSys) String() string {
	return "gaefs"
}

func (fs AeFileSys) Name() string {
	return fs.String()
}

// conflicts with Open() interface of VFS
func (fs *AeFileSys) Open(name string) (AeFile, error) {
	return AeFile{}, nil
}

func (fs *AeFileSys) OpenFile(name string, flag int, perm os.FileMode) (AeFile, error) {
	return fs.GetFile(name)
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

func (fs *AeFileSys) WriteFile(name string, data []byte, perm os.FileMode) error {

	name = cleanseLeadingSlash(name)

	f, err := fs.Create(name)
	if err != nil {
		return err
	}
	f.Data = data
	err = fs.SaveFile(&f, name)
	return err
}
