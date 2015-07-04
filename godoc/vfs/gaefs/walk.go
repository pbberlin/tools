package gaefs

import (
	"errors"
	"fmt"
	"os"
	pth "path"
	"path/filepath"

	"github.com/pbberlin/tools/logif"
)

type WalkFunc func(path string, info os.FileInfo, err error) error

var SkipDir = errors.New("skip this directory") // walk func signalling: omit this dir

// walk recursively descends path, calling walkFn.
func (fs *AeFileSys) walk(path string, info os.FileInfo, walkFn WalkFunc) error {

	err := walkFn(path, info, nil)
	if err != nil {
		if info.IsDir() && err == SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	names, err := fs.Readdirnames(path)
	logif.Pf("%v => %+v", path, names)
	if err != nil {
		return walkFn(path, info, err)
	}

	for _, name := range names {
		filename := pth.Join(path, name)
		fileInfo, err := fs.Lstat(filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != SkipDir {
				return err
			}
		} else {
			err = fs.walk(filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

// Walk walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn. The files are walked in lexical
// order, which makes the output deterministic but means that for very
// large directories Walk can be inefficient.
// Walk does not follow symbolic links.
func (fs *AeFileSys) Walk(root string, walkFn WalkFunc) error {
	info, err := fs.Lstat(root)
	if err != nil {
		// logif.Pf("walk start error %v", err)
		return walkFn(root, nil, err)
	}
	return fs.walk(root, info, walkFn)
}

func EXAMPLES() {}

// example walk func
func exWalkFunc(path string, f os.FileInfo, err error) error {
	tp := "file"
	if f.IsDir() {
		tp = "dir "
	}
	fmt.Printf("Visited: %s %s \n", tp, path)
	return nil
}

// example walk
func exWalk() {
	root := "/"
	err := filepath.Walk(root, exWalkFunc)
	fmt.Printf("filepath.Walk() returned %v\n", err)
}
