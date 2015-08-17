package memfs

import "github.com/pbberlin/tools/os/fsi/common"

// name is the *external* path or filename.
func (fs *memMapFs) SplitX(name string) (dir, bname string) {

	return common.UnixPather(name, fs.RootDir())

}
