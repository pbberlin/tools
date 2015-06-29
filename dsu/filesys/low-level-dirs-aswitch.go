package filesys

var NestedOrRooted = true

func (fs FileSys) GetDirByPath(path string) (Directory, error) {
	if NestedOrRooted {
		return fs.nestedGetDirByPath(path)
	} else {
		return fs.rootedGetDirByPath(path)
	}
}

func (fs FileSys) SaveDirByPath(path string) (Directory, error) {
	if NestedOrRooted {
		return fs.nestedSaveDirByPath(path)
	} else {
		return fs.rootedSaveDirByPath(path)
	}
}
