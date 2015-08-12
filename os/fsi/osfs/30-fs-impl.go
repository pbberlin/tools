package osfs

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/pbberlin/tools/os/fsi"
)

func (fs *osFileSys) Name() string { return "osfs" }

func (fs osFileSys) String() string {
	hn, err := os.Hostname()
	if err != nil {
		return err.Error()
	}
	return hn
}

func (fs *osFileSys) Repl(name string) string {
	if fs.replacePath {
		return strings.Replace(name, "/", "\\", -1)
	}
	return name
}

//---------------------------------------

func (fs *osFileSys) Chmod(name string, mode os.FileMode) error {
	name = fs.Repl(name)
	return os.Chmod(name, mode)
}

func (fs *osFileSys) Chtimes(name string, atime time.Time, mtime time.Time) error {
	name = fs.Repl(name)
	return os.Chtimes(name, atime, mtime)
}

func (fs *osFileSys) Create(name string) (fsi.File, error) {
	name = fs.Repl(name)
	return os.Create(name)
}

func (fs *osFileSys) Lstat(path string) (os.FileInfo, error) {
	path = fs.Repl(path)
	return os.Lstat(path)
}

func (fs *osFileSys) Mkdir(name string, perm os.FileMode) error {
	name = fs.Repl(name)
	return os.Mkdir(name, perm)
}

func (fs *osFileSys) MkdirAll(path string, perm os.FileMode) error {
	path = fs.Repl(path)
	return os.MkdirAll(path, perm)
}

func (fs *osFileSys) Open(name string) (fsi.File, error) {
	name = fs.Repl(name)
	return os.Open(name)
}

func (fs *osFileSys) OpenFile(name string, flag int, perm os.FileMode) (fsi.File, error) {
	name = fs.Repl(name)
	return os.OpenFile(name, flag, perm)
}

func (fs *osFileSys) ReadDir(dirname string) ([]os.FileInfo, error) {
	dirname = fs.Repl(dirname)
	return ioutil.ReadDir(dirname)
}

func (fs *osFileSys) Remove(name string) error {
	name = fs.Repl(name)
	return os.Remove(name)
}

func (fs *osFileSys) RemoveAll(path string) error {
	path = fs.Repl(path)
	return os.RemoveAll(path)
}

func (fs *osFileSys) Rename(oldname, newname string) error {
	oldname = fs.Repl(oldname)
	newname = fs.Repl(newname)
	return os.Rename(oldname, newname)
}

func (fs *osFileSys) Stat(name string) (os.FileInfo, error) {
	name = fs.Repl(name)
	return os.Stat(name)
}

func (fs *osFileSys) ReadFile(filename string) ([]byte, error) {
	filename = fs.Repl(filename)
	return ioutil.ReadFile(filename)
}

func (fs *osFileSys) WriteFile(filename string, data []byte, perm os.FileMode) error {
	filename = fs.Repl(filename)
	return ioutil.WriteFile(filename, data, perm)
}
