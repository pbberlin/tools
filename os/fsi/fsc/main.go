// Package fsc - filesystem common - contains convenience functions,
// common to all fsi implementations;
// since only methods from the fsi interface are used;
// sadly we cannot attach methods to interfaces.
package fsc

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

// walk recursively descends path, calling walkFn.
func walk(fs fsi.FileSystem, path string, info os.FileInfo, walkFn WalkFunc) error {

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
	// logif.Pf("%11v => %+v", path, fis)
	if err != nil {
		return walkFn(path, info, err)
	}

	for _, fi := range fis {
		filename := pth.Join(path, fi.Name())
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
// It requires only the fsi.FileSystem interface, and is therefore implementation indepdenent.
//
// It is similar to filepath.Walk(root string, walkFunc)
// In contrast to filepath.Walk, it does not enumerate files, only dirs.
// File enumeration can be added in the WalkFunc.
//
// Directories are walked in order of Readdirnames()
//
// Errors that arise visiting directories can be filtered by walkFn.
//
// Walk does not follow symbolic links.
func Walk(fs fsi.FileSystem, root string, walkFn WalkFunc) error {
	info, err := fs.Lstat(root)
	if err != nil {
		// logif.Pf("walk start error %v", err)
		return walkFn(root, nil, err)
	}
	return walk(fs, root, info, walkFn)
}
