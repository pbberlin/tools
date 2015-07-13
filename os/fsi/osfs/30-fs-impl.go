package osfs

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/pbberlin/tools/os/fsi"
)

func (osFileSys) Name() string { return "osfs" }

func (fs osFileSys) String() string {
	hn, err := os.Hostname()
	if err != nil {
		return err.Error()
	}
	return hn
}

//---------------------------------------

func (osFileSys) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (osFileSys) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (osFileSys) Create(name string) (fsi.File, error) {
	return os.Create(name)
}

func (osFileSys) Lstat(path string) (os.FileInfo, error) {
	return os.Lstat(path)
}

func (osFileSys) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (osFileSys) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (osFileSys) Open(name string) (fsi.File, error) {
	return os.Open(name)
}

func (osFileSys) OpenFile(name string, flag int, perm os.FileMode) (fsi.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (osFileSys) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (osFileSys) Remove(name string) error {
	return os.Remove(name)
}

func (osFileSys) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (osFileSys) Rename(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

func (osFileSys) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (osFileSys) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (osFileSys) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}
