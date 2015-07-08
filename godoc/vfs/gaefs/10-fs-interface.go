package gaefs

import "os"

// Interface FileSystem is inspired by os.File + io.ioutil,
// informed by godoc.vfs and package afero.
type FileSystem interface {
	Create(name string) (*AeFile, error) // read write
	Lstat(path string) (os.FileInfo, error)
	Mkdir(name string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Name() string
	Open(name string) (*AeFile, error) // read only
	OpenFile(name string, flag int, perm os.FileMode) (*AeFile, error)

	// Notice the distinct methods on File interface:
	//          Readdir(count int) ([]os.FileInfo, error)
	//          Readdirnames(n int) ([]string, error)
	ReadDir(dirname string) ([]os.FileInfo, error) // from io.ioutil

	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldname, newname string) error
	Stat(path string) (os.FileInfo, error)
	String() string

	// Two convenience methods taken from io.ioutil, that we want to rely on
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error

	// Inspired by filepath.Walk.
	// Could be implemented generically on interface FileSystem;
	//   with interface methods only. But golang does not support methods on interfaces.
	Walk(root string, walkFn WalkFunc) error
}
