package gaefs

import (
	"io"
	"os"
)

// File represents a file in the filesystem.
type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	io.WriterAt
	//Fd() uintptr
	Stat() (os.FileInfo, error)
	// Readdir(count int) ([]os.FileInfo, error)
	// Readdirnames(n int) ([]string, error)
	WriteString(s string) (ret int, err error)
	Truncate(size int64) error
	Name() string
}
