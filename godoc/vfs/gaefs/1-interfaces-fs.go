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
type T_OS func(string) AeFileSys
type T_ReadFile func(Opener, string) ([]byte, error)

type FileSystemVFS interface {
	Opener
	Lstat(path string) (os.FileInfo, error)
	Stat(path string) (os.FileInfo, error)
	ReadDir(path string) ([]os.FileInfo, error)
	String() string
}

// from afero - https://github.com/spf13/afero
// Stat() and Open(...) overlap with vfs interface
type FileSystemI2 interface {
	Create(name string) (AeFile, error)
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Name() string
	Open(name string) (AeFile, error)
	OpenFile(name string, flag int, perm os.FileMode) (AeFile, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldname, newname string) error
	Stat(name string) (os.FileInfo, error)
}

// merged
type FileSystem interface {
	Create(name string) (AeFile, error)
	Lstat(path string) (os.FileInfo, error)
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Name() string
	Open(name string) (AeFile, error)
	OpenFile(name string, flag int, perm os.FileMode) (AeFile, error)
	ReadDir(path string) ([]os.FileInfo, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldname, newname string) error
	Stat(path string) (os.FileInfo, error)
	String() string
}
