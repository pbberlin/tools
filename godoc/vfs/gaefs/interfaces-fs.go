package gaefs

import (
	"os"

	"golang.org/x/tools/godoc/vfs"
)

// method signiture conflicted with Afero Open()
type Opener interface {
	OpenVFS(name string) (vfs.ReadSeekCloser, error)
}

// I wanted my package types <Directory> and <File>
// pluggable into
//    pth.Walk(File,WalkFunc)
// But it seems impossible :(

//
// from golang.org/x/tools/godoc/vfs
type T_OS func(string) FileSys
type T_Readfile func(Opener, string) ([]byte, error)
type FileSystem interface {
	Opener
	Lstat(path string) (os.FileInfo, error)
	Stat(path string) (os.FileInfo, error)
	ReadDir(path string) ([]os.FileInfo, error)
	String() string
}

// from afero - https://github.com/spf13/afero
// Stat() and Open(...) overlap with vfs interface
type FileSystemI2 interface {
	Create(name string) (File, error)
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Name() string
	Open(name string) (File, error)
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldname, newname string) error
	Stat(name string) (os.FileInfo, error)
}
