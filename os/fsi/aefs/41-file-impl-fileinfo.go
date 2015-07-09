package aefs

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

func (d AeDir) Name() string {
	return d.BName
}
func (f AeFile) Name() string {
	return f.BName
}

func (d AeDir) Size() int64 {
	return int64(len(d.BName))
}
func (f AeFile) Size() int64 {
	return int64(len(f.Data))
}

// no rights implemented
func (d AeDir) Mode() os.FileMode {
	return os.ModePerm
}
func (f AeFile) Mode() os.FileMode {
	return f.MMode
}

func (d AeDir) ModTime() time.Time {
	return d.MModTime
}
func (f AeFile) ModTime() time.Time {
	return f.MModTime
}

func (d AeDir) IsDir() bool {
	return true
}
func (f AeFile) IsDir() bool {
	return false
}
func (d AeDir) Sys() interface{} {
	return nil
}
func (f AeFile) Sys() interface{} {
	return nil
}
