package gaefs

func (fs FileSys) GetDirByPath(path string) (Directory, error) {

	if fs.rooted {
		return fs.rootedGetDirByPath(path)
	}
	return fs.nestedGetDirByPath(path)
}

func (fs FileSys) SaveDirByPath(path string) (Directory, error) {
	if fs.rooted {
		return fs.rootedSaveDirByPath(path)
	}
	return fs.nestedSaveDirByPath(path)
}
