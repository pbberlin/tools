package gaefs

import (
	"os"
	"time"
)
import pth "path"

func (fs *AeFileSys) Create(name string) (AeFile, error) {

	name = cleanseLeadingSlash(name)

	f := AeFile{}
	dir, base := pth.Split(name)
	f.name = base
	err := fs.SaveFile(&f, dir)
	if err != nil {
		return f, err
	}
	return f, err
}

func (fs *AeFileSys) Lstat(path string) (os.FileInfo, error) {
	panic(spf("Links not implemented for %v", fs))
	var fi os.FileInfo
	return fi, nil
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
	panic(spf("RemoveAll not (yet) implemented for %v", fs))
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

func (fs *AeFileSys) Chmod(name string, mode os.FileMode) error {
	panic(spf("Chmod not (yet) implemented for %v", fs))
	return nil
}

func (fs *AeFileSys) Chtimes(name string, atime time.Time, mtime time.Time) error {
	panic(spf("Chtimes not (yet) implemented for %v", fs))
	return nil
}
