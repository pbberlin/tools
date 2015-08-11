// Package fsi serves as universal abstraction
// layer for all go programs that require a file system.

package fsi

import (
	"errors"
	"fmt"
	"os"
)

var (
	// EmptyQueryResult is a warning, that implementations of ReadDir may return,
	// if their results are based on weakly consistent indexes.
	// It is defined here, since common.Walk() wants to ignore it.
	EmptyQueryResult = fmt.Errorf("Query found no results based on weakly consistent index.")

	// If an implementation cannot support a method, it should at least return this testable error.
	NotImplemented = fmt.Errorf("Filesystem does not support this method.")

	ErrRootDirNoFile = fmt.Errorf("rootdir; no file")

	ErrFileClosed = errors.New("File is closed")
	ErrFileInUse  = errors.New("File already in use")
	ErrOutOfRange = errors.New("Out of range")
	ErrTooLarge   = errors.New("Too large")

	ErrFileNotFound      = os.ErrNotExist
	ErrFileExists        = os.ErrExist
	ErrDestinationExists = os.ErrExist
)

// Interface FileSystem is inspired by os.File + io.ioutil,
// informed by godoc.vfs and package afero.
type FileSystem interface {
	Name() string   // the type of filesystem, i.e. "osfs"
	String() string // a mountpoint or a drive letter

	// Nobody restricts you from implementing following methods.
	// But they are not mandatory for our interface:
	// Chmod(name string, mode os.FileMode) error
	// Chtimes(name string, atime time.Time, mtime time.Time) error

	Create(name string) (File, error)       // read write
	Lstat(path string) (os.FileInfo, error) // for common.Walk
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Open(name string) (File, error) // read only
	OpenFile(name string, flag int, perm os.FileMode) (File, error)

	// ReadDir is taken from io.ioutil.
	// It should return directories first, then files second.
	// Notice the distinct methods on File interface:
	//          Readdir(count int) ([]os.FileInfo, error)
	//          Readdirnames(n int) ([]string, error)
	// Those coming from os.File.
	// We base all those methods on a single internal implementation.
	// Readdir may return EmptyQueryResult error as a warning.
	//
	// Compare File.Readdir remarks.
	ReadDir(dirname string) ([]os.FileInfo, error)

	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldname, newname string) error
	Stat(path string) (os.FileInfo, error)

	// Two convenience methods taken from io.ioutil.
	// They are mandatory, because you will need them sooner or later anyway.
	// Thus we require them right from the start and with *standard* signatures.
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error

	// Walk() is inspired by filepath.Walk()
	// Walk is implemented generically, purely on fsi.FileSystem,
	// in package common. Implementing it here is discouraged:
	// Walk(root string, walkFn WalkFunc) error

}
