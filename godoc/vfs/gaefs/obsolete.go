package gaefs

import "strings"
import pth "path"

// before using constructDirKey
func (fs AeFileSys) nestedGetDirByPath_OOOOOOLD(path string) (AeDir, error) {

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
