package dsfs

import (
	"os"
	"time"
)

// golang.org/pkg/os
type FileInfo interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() os.FileMode  // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)
}

// All of them: NO pointer receiver

func (d DsDir) Name() string {
	return d.BName
}
func (f DsFile) Name() string {
	return f.BName
}

func (d DsDir) Size() int64 {
	return int64(len(d.BName))
}
func (f DsFile) Size() int64 {
	return int64(len(f.Data))
}

// no rights implemented
func (d DsDir) Mode() os.FileMode {
	return os.ModePerm
}
func (f DsFile) Mode() os.FileMode {
	return f.MMode
}

func (d DsDir) ModTime() time.Time {
	return d.MModTime
}
func (f DsFile) ModTime() time.Time {
	return f.MModTime
}

func (d DsDir) IsDir() bool {
	return true
}
func (f DsFile) IsDir() bool {
	return false
}
func (d DsDir) Sys() interface{} {
	return nil
}
func (f DsFile) Sys() interface{} {
	return nil
}
