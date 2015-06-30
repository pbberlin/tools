package gaefs

import (
	"os"
	"time"

	"golang.org/x/tools/godoc/vfs"
)

type Opener interface {
	Open(name string) (vfs.ReadSeekCloser, error)
}

// from golang.org/x/tools/godoc/vfs
type FileSystem interface {
	Opener
	Lstat(path string) (os.FileInfo, error)
	Stat(path string) (os.FileInfo, error)
	ReadDir(path string) ([]os.FileInfo, error)
	String() string
}

// golang.org/pkg/os
type FileInfo interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() os.FileMode  // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)
}

// The package types <Directory> and <File> should
// implement some of these methods.
// I was hoping to plug my package into
//    pth.Walk(File,WalkFunc)
// But it seems impossible :(
