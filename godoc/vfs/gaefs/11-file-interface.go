package gaefs

import (
	"io"
	"os"
)

// Interface File is inspired by os.File, informed by godoc.vfs and package afero
type FileI interface {
	io.Closer
	// Notice, we dont have an Opener on file. Opener is attached to filesystem one level higher.
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer

	io.WriterAt
	//Fd() uintptr
	Stat() (os.FileInfo, error)

	Readdir(count int) ([]os.FileInfo, error) // notice distinctive signature FileSys.ReadDir(...)
	Readdirnames(n int) ([]string, error)

	WriteString(s string) (ret int, err error)
	Truncate(size int64) error

	Name() string

	// Notice indirect need to support os.FileInfo by the means of Stat()
	//	Size() int64 {
	//	Mode() os.FileMode {
	//	ModTime() time.Time {
	//	IsDir() bool {
	//	Sys() interface{} {

}
