package gaefs

func (fs AeFileSys) GetDirByPath(path string) (AeDir, error) {
	path = cleanseLeadingSlash(path)

	if fs.Rooted {
		return fs.rootedGetDirByPath(path)
	}
	return fs.nestedGetDirByPath(path)
}

func (fs AeFileSys) SaveDirByPath(path string) (AeDir, error) {
	path = cleanseLeadingSlash(path)

	if fs.Rooted {
		return fs.rootedSaveDirByPath(path)
	}
	return fs.nestedSaveDirByPath(path)
}
