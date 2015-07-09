// Package fsi - filesystem interface - contains the minimal
// requirements for exchangeable filesystems.
package fsi

import (
	"os"
	"time"
)

// Interface FileSystem is inspired by os.File + io.ioutil,
// informed by godoc.vfs and package afero.
type FileSystem interface {
	Name() string   // the type
	String() string // a mountpoint or a drive letter

	Chmod(name string, mode os.FileMode) error
	Chtimes(name string, atime time.Time, mtime time.Time) error

	Create(name string) (File, error) // read write
	Lstat(path string) (os.FileInfo, error)
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Open(name string) (File, error) // read only
	OpenFile(name string, flag int, perm os.FileMode) (File, error)

	// Notice the distinct methods on File interface:
	//          Readdir(count int) ([]os.FileInfo, error)
	//          Readdirnames(n int) ([]string, error)
	ReadDir(dirname string) ([]os.FileInfo, error) // from io.ioutil

	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldname, newname string) error
	Stat(path string) (os.FileInfo, error)

	// Two convenience methods taken from io.ioutil, that we want to rely on
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error

	// Inspired by filepath.Walk.
	// Also implemented generically, purely on fsi.FileSystem;
	//   in package fsc. However, only functions are possible,
	//   since golang does not support methods on interfaces.
	// Walk(root string, walkFn WalkFunc) error

}

type WalkFunc func(path string, info os.FileInfo, err error) error
