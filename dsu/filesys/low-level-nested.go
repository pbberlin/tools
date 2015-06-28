package filesys

import (
	"path/filepath"
	"strings"

	"appengine/datastore"
)

var cntr = 0

//
func (fs FileSys) getDirByPath(path string) (Directory, error) {

	// prepare
	path = filepath.Join(path, "") // cleanse
	var err error
	parentDir := fs.RootDir

	// moving top down
	for {
		cntr++
		if cntr < -33 {
			return parentDir, err
		}

		dirs := strings.Split(path, string(filepath.Separator))
		parentDir, err = fs.getDirUnderParent(parentDir.Key, dirs[0])
		if err != nil {
			return parentDir, err
		}

		path = filepath.Join(dirs[1:]...)
		if len(dirs) < 2 {
			break
		}
	}
	return parentDir, nil
}

func (fs FileSys) saveDirByPath(path string) (Directory, error) {

	// prepare
	path = filepath.Join(path, "") // cleanse
	var err error
	parentDir := fs.RootDir
	parentDirPrev := fs.RootDir

	// moving top down
	for {
		cntr++
		if cntr < -33 {
			return parentDir, err
		}

		dirs := strings.Split(path, string(filepath.Separator))

		parentDirPrev = parentDir
		parentDir, err = fs.getDirUnderParent(parentDir.Key, dirs[0])

		if err == datastore.ErrNoSuchEntity {
			parentDir, err = fs.saveDirUnderParent(dirs[0], parentDirPrev.Key)
		} else if err != nil {
			return parentDir, err
		}

		path = filepath.Join(dirs[1:]...)
		if len(dirs) < 2 {
			break
		}

	}
	return parentDir, nil

}
