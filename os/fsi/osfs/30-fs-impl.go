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

func (fs *osFileSys) WinGoofify(name string) string {
	if fs.replacePath {
		return strings.Replace(name, "/", "\\", -1)
	}
	return name
}

//---------------------------------------

func (fs *osFileSys) Chmod(name string, mode os.FileMode) error {
	name = fs.WinGoofify(name)
	return os.Chmod(name, mode)
}

func (fs *osFileSys) Chtimes(name string, atime time.Time, mtime time.Time) error {
	name = fs.WinGoofify(name)
	return os.Chtimes(name, atime, mtime)
}

func (fs *osFileSys) Create(name string) (fsi.File, error) {
	name = fs.WinGoofify(name)
	return os.Create(name)
}

func (fs *osFileSys) Lstat(path string) (os.FileInfo, error) {
	path = fs.WinGoofify(path)
	return os.Lstat(path)
}

func (fs *osFileSys) Mkdir(name string, perm os.FileMode) error {
	name = fs.WinGoofify(name)
	return os.Mkdir(name, perm)
}

func (fs *osFileSys) MkdirAll(path string, perm os.FileMode) error {
	path = fs.WinGoofify(path)
	return os.MkdirAll(path, perm)
}

// Type wraps os.File to overwrite os.File.Readdir()
type fileWithCustomReaddir struct {
	//                     fileWithCustomReaddir needs to fullfill fsi.File interface:
	*os.File            // first  - wrap anonymously     => type implements all methods of os.File
	f        *os.File   // second - wrap as named member => to reach back from the overwritten method
	parFS    *osFileSys // to get the current readdirsorter
}

// Overwritten method
func (f fileWithCustomReaddir) Readdir(n int) ([]os.FileInfo, error) {
	fis, err := f.f.Readdir(n) // reaching back to os.File.Readdir()
	f.parFS.readdirsorter(fis) // additional logic - purpose of the entire exercise
	return fis, err
}

func (fs *osFileSys) Open(name string) (fsi.File, error) {
	name = fs.WinGoofify(name)
	osfile, err := os.Open(name)
	wrappedOsFile := fileWithCustomReaddir{osfile, osfile, fs} // wrap it into
	return wrappedOsFile, err
}

func (fs *osFileSys) OpenFile(name string, flag int, perm os.FileMode) (fsi.File, error) {
	name = fs.WinGoofify(name)
	return os.OpenFile(name, flag, perm)
}

func (fs *osFileSys) ReadDir(dirname string) ([]os.FileInfo, error) {
	fis, err := ioutil.ReadDir(dirname)
	fs.readdirsorter(fis)
	return fis, err
}

func (fs *osFileSys) Remove(name string) error {
	name = fs.WinGoofify(name)
	return os.Remove(name)
}

func (fs *osFileSys) RemoveAll(path string) error {
	path = fs.WinGoofify(path)
	return os.RemoveAll(path)
}

func (fs *osFileSys) Rename(oldname, newname string) error {
	oldname = fs.WinGoofify(oldname)
	newname = fs.WinGoofify(newname)
	return os.Rename(oldname, newname)
}

func (fs *osFileSys) Stat(name string) (os.FileInfo, error) {
	name = fs.WinGoofify(name)
	return os.Stat(name)
}

func (fs *osFileSys) ReadFile(filename string) ([]byte, error) {
	filename = fs.WinGoofify(filename)
	return ioutil.ReadFile(filename)
}

func (fs *osFileSys) WriteFile(filename string, data []byte, perm os.FileMode) error {
	filename = fs.WinGoofify(filename)
	return ioutil.WriteFile(filename, data, perm)
}
