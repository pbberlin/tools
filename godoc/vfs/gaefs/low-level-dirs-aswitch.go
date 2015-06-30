package gaefs

import "strings"

func (fs FileSys) GetDirByPath(path string) (Directory, error) {
	path = cleanseLeadingSlash(path)

	if fs.rooted {
		return fs.rootedGetDirByPath(path)
	}
	return fs.nestedGetDirByPath(path)
}

func (fs FileSys) SaveDirByPath(path string) (Directory, error) {
	path = cleanseLeadingSlash(path)

	if fs.rooted {
		return fs.rootedSaveDirByPath(path)
	}
	return fs.nestedSaveDirByPath(path)
}

func cleanseLeadingSlash(p string) string {
	for {
		if strings.HasPrefix(p, sep) {
			p = p[1:]
		} else {
			break
		}
	}
	return p
}
