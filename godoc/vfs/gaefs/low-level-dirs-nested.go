package gaefs

import (
	"path/filepath"
	"strings"

	"appengine/datastore"
)

func (fs FileSys) nestedGetDirByPath(path string) (Directory, error) {

	// prepare
	path = filepath.Join(path, "") // cleanse
	var err error
	childDir := fs.RootDir

	// moving top down
	for {
		dirs := strings.Split(path, string(filepath.Separator))
		childDir, err = fs.getDirUnderParent(childDir.Key, dirs[0])
		if err != nil {
			return childDir, err
		}

		dirs = dirs[1:]
		path = filepath.Join(dirs...)
		if len(dirs) < 1 {
			break
		}
	}
	return childDir, nil
}

func (fs FileSys) nestedSaveDirByPath(path string) (Directory, error) {

	// prepare
	path = filepath.Join(path, "") // cleanse
	var err error
	childDir := fs.RootDir
	childDirPrev := fs.RootDir

	// moving top down
	for {
		dirs := strings.Split(path, string(filepath.Separator))

		childDirPrev = childDir
		childDir, err = fs.getDirUnderParent(childDir.Key, dirs[0])

		if err == datastore.ErrNoSuchEntity {
			childDir, err = fs.saveDirUnderParent(dirs[0], childDirPrev.Key)
		} else if err != nil {
			return childDir, err
		}

		dirs = dirs[1:]
		path = filepath.Join(dirs...)
		if len(dirs) < 1 {
			break
		}

	}
	return childDir, nil

}
