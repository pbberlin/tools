// Package common - contains convenience functions,
// built on top of fsi interfaces;
// thus independent of implementations.
// Sadly we cannot attach methods to interfaces.
// Generic filesystem utils should be implemented here.
// find, compact, export ...
package common

import (
	"errors"
	"os"
	pth "path"

	"github.com/pbberlin/tools/os/fsi"
)

// SkipDir is an "error", which a walk-function can
// return, in order to signal, that walk should not traverse into this dir.
var SkipDir = errors.New("skip this directory")

type WalkFunc func(path string, info os.FileInfo, err error) error

var cntr = 0

// walk recursively descends path, calling walkFn.
func walk(fs fsi.FileSystem, path string, info os.FileInfo, walkFn WalkFunc) error {

	// cntr++
	// if cntr > 20 {
	// 	return fmt.Errorf("too many recursions")
	// }

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

	fis, err := fs.ReadDir(path)
	// fnd := ""
	// for i := 0; i < len(fis); i++ {
	// 	fnd += fis[i].Name() + ", "
	// }
	// log.Printf("readdir of %-26v  => %v, %v", path, len(fis), fnd)
	if err != nil && err != fsi.EmptyQueryResult {
		return walkFn(path, info, err)
	}

	//
	for _, fi := range fis {
		filename := pth.Join(path, pth.Base(fi.Name()))

		fileInfo, err := fs.Lstat(filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != SkipDir {
				return err
			}
		} else {
			err = walk(fs, filename, fileInfo, walkFn)
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
//
// It requires only the fsi.FileSystem interface, and is therefore implementation independent.
//
// It is similar to filepath.Walk(root string, walkFunc)
//
// Directories are crawled in order of fs.ReadDir()
// Walk crawls directories first, files second.
//
// Errors that arise visiting directories can be filtered by walkFn.
//
// Walk does not follow symbolic links.
func Walk(fs fsi.FileSystem, root string, walkFn WalkFunc) error {
	info, err := fs.Lstat(root)
	if err != nil {
		// log.Printf("walk start error %10v %v", root, err)
		return walkFn(root, nil, err)
	}
	// log.Printf("walk start fnd %v", info.Name())
	return walk(fs, root, info, walkFn)
}
