package filesys

import (
	"os"
	"time"
)

func (fs FileSys) mkdir(path string) {
	fs.recurseMkDir(path)
}

func (fs FileSys) touch(p string) {

}

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
