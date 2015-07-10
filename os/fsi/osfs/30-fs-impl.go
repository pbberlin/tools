package osfs

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/pbberlin/tools/os/fsi"
)

func (OsFileSys) Name() string { return "osfs" }

func (fs OsFileSys) String() string {
	hn, err := os.Hostname()
	if err != nil {
		return err.Error()
	}
	return hn
}

//---------------------------------------

func (OsFileSys) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (OsFileSys) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (OsFileSys) Create(name string) (fsi.File, error) {
	return os.Create(name)
}

func (OsFileSys) Lstat(path string) (os.FileInfo, error) {
	return os.Lstat(path)
}

func (OsFileSys) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (OsFileSys) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (OsFileSys) Open(name string) (fsi.File, error) {
	return os.Open(name)
}

func (OsFileSys) OpenFile(name string, flag int, perm os.FileMode) (fsi.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (OsFileSys) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (OsFileSys) Remove(name string) error {
	return os.Remove(name)
}

func (OsFileSys) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (OsFileSys) Rename(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

func (OsFileSys) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OsFileSys) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (OsFileSys) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}
