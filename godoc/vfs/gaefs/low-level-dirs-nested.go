package gaefs

import (
	pth "path"
	"strings"

	"appengine/datastore"
)

func (fs AeFileSys) nestedGetDirByPath(path string) (AeDir, error) {

	if path == "" {
		return fs.RootDir, nil
	}

	// prepare
	var err error
	childDir := fs.RootDir

	// moving top down
	for {
		dirs := strings.Split(path, sep)
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

func (fs AeFileSys) nestedSaveDirByPath(path string) (AeDir, error) {

	if path == "" {
		return fs.RootDir, nil
	}

	// prepare
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
