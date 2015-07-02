package gaefs

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

func (d Directory) Name() string {
	return d.BName
}
func (f File) Name() string {
	return f.BName
}

func (d Directory) Size() int64 {
	return int64(len(d.BName))
}
func (f File) Size() int64 {
	return int64(len(f.Content))
}

// no rights implemented
func (d Directory) Mode() os.FileMode {
	return os.ModePerm
}
func (f File) Mode() os.FileMode {
	return f.mode
}

func (d Directory) ModTime() time.Time {
	return d.Mod
}
func (f File) ModTime() time.Time {
	return f.Mod
}

func (d Directory) IsDir() bool {
	return true
}
func (f File) IsDir() bool {
	return false
}
func (d Directory) Sys() interface{} {
	return nil
}
func (f File) Sys() interface{} {
	return nil
}
