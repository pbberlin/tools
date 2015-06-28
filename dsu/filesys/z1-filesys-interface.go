package filesys

import (
	"os"
	"time"
)

// The package types <Directory> and <File> should
// implement some of these methods.
// I was hoping to plug my package into
//    filepath.Walk(File,WalkFunc)
// But it seems impossible :(

type FileInfo interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() os.FileMode  // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)

	Parent() string // Added
	Path() string   // Added
}
