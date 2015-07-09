package aefs_sr

import (
	"os"

	"golang.org/x/tools/godoc/vfs"
)

//
// from golang.org/x/tools/godoc/vfs
type T_OS func(string) AeFileSys
type T_ReadFile func(Opener, string) ([]byte, error)

// I wanted my package types <Directory> and <File>
// pluggable into
//    pth.Walk(File,WalkFunc)
// But it seems impossible :(

type FileSystemVFS interface {
	Opener
	Lstat(path string) (os.FileInfo, error)
	Stat(path string) (os.FileInfo, error)
	ReadDir(path string) ([]os.FileInfo, error) // uppercase dir :( - different from file interface
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

// method signiture conflicted with Afero Open()
type Opener interface {
	OpenVFS(name string) (vfs.ReadSeekCloser, error)
}

func OS(mount string) AeFileSys {
	panic(`
		Sadly, google app engine file system requires a
	 	http.Request based context object.
	 	Use NewFs(string, AeContext(c)) instead of OS.
	`)
}
