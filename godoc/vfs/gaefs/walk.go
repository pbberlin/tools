package gaefs

import (
	"errors"
	"fmt"
	"os"
	pth "path"
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
	// logif.Pf("%11v => %+v", path, names)
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
// directory in the tree, including root.
// But it does not enumerate files.
// The directories are walked in lexical order.
//
// It is similar to filepath.Walk(root string, walkFunc)
//
// Errors that arise visiting directories can be filtered by walkFn.
//
// Walk does not follow symbolic links.
func (fs *AeFileSys) Walk(root string, walkFn WalkFunc) error {
	info, err := fs.Lstat(root)
	if err != nil {
		// logif.Pf("walk start error %v", err)
		return walkFn(root, nil, err)
	}
	return fs.walk(root, info, walkFn)
}

// example walk
func ExampleWalk() {

	// first the per-node func:
	exWalkFunc := func(path string, f os.FileInfo, err error) error {

		if err == fmt.Errorf("My special error") {
			err = nil
		} else if err != nil {
			return err // calling off the walk
		}

		tp := "file"
		if f.IsDir() {
			tp = "dir "
		}
		fmt.Printf("Visited: %s %s \n", tp, path)
		return nil
	}

	mnt := "mnt00"
	fs := NewAeFs(mnt) // add appengine context

	err := fs.Walk(mnt, exWalkFunc)
	fmt.Printf("Walk() returned %v\n", err)
}
