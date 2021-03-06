package fsi

import (
	"io"
	"os"
)

// Interface File is inspired by os.File,
// and informed by godoc.vfs and package afero
type File interface {
	io.Closer
	// Notice, we dont have an Opener on file.
	// Opener is attached to filesystem one level higher.
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	io.WriterAt

	//Fd() uintptr
	Stat() (os.FileInfo, error)

	// Close writes to disk; Sync is not neccessary
	// Sync() error

	// Readdir and Readdirnames come from os.File.
	// Notice the distinctive signature and the remarks on FileSys.ReadDir(...)
	//
	// Notice the contract https://golang.org/src/os/doc.go
	// For n > 0, returning io.EOF is important,
	// otherwise http.FileServer will repeat queries forever.
	//
	// Example at https://golang.org/src/os/file_windows.go
	Readdir(n int) ([]os.FileInfo, error)
	Readdirnames(n int) ([]string, error)

	WriteString(s string) (ret int, err error)
	Truncate(size int64) error

	Name() string

	// Notice the indirect need to implement os.FileInfo
	// because it is returned by Stat()
	//	  Size() int64 {
	//	  Mode() os.FileMode {
	//	  ModTime() time.Time {
	//	  IsDir() bool {
	//	  Sys() interface{} {

}
