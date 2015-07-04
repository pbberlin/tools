package gaefs

import (
	pth "path"
	"strings"

	"appengine/datastore"
)

func (fs AeFileSys) nestedGetDirByPath(path string) (AeDir, error) {
	k := fs.constructDirKey(path)
	d, err := fs.getDirByExactKey(k)
	return d, err
}

// We can derive a directory key
// without loading any of the parental directory objects.
// But
func (fs AeFileSys) constructDirKey(path string) (k *datastore.Key) {

	// always starting with root
	k = datastore.NewKey(fs.Ctx(), tdir, fs.RootDir.BName, 0, nil)

	if path == "" {
		return
	}

	// moving top down
	for {
		dirs := strings.Split(path, sep)
		k = datastore.NewKey(fs.Ctx(), tdir, dirs[0], 0, k)
		dirs = dirs[1:]
		path = pth.Join(dirs...)
		if len(dirs) < 1 {
			break
		}
	}
	return
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
			childDir, err = fs.saveDirUnderParent(dirs[0], childDirPrev.Key) // create it
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
