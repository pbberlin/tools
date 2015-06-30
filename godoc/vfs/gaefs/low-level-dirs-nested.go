package gaefs

import (
	pth "path"
	"strings"

	"github.com/pbberlin/tools/logif"

	"appengine/datastore"
)

func (fs FileSys) nestedGetDirByPath(path string) (Directory, error) {

	// prepare
	path = pth.Join(path)
	var err error
	childDir := fs.RootDir

	// moving top down
	for {
		dirs := strings.Split(path, sep)
		logif.Pf("so far 2 %v - %v - --%v--", path, dirs, childDir.Key)
		childDir, err = fs.getDirUnderParent(childDir.Key, dirs[0])
		if err != nil {
			return childDir, err
		}

		dirs = dirs[1:]
		path = pth.Join(dirs...)

		if len(dirs) < 1 {
			break
		}
	}
	return childDir, nil
}

func (fs FileSys) nestedSaveDirByPath(path string) (Directory, error) {

	// prepare
	path = pth.Join(path) // cleanse
	var err error
	childDir := fs.RootDir
	childDirPrev := fs.RootDir

	// moving top down
	for {
		dirs := strings.Split(path, sep)

		childDirPrev = childDir
		childDir, err = fs.getDirUnderParent(childDir.Key, dirs[0])

		if err == datastore.ErrNoSuchEntity {
			childDir, err = fs.saveDirUnderParent(dirs[0], childDirPrev.Key)
		} else if err != nil {
			return childDir, err
		}

		dirs = dirs[1:]
		path = pth.Join(dirs...)
		if len(dirs) < 1 {
			break
		}

	}
	return childDir, nil

}
