package aefs

import "github.com/pbberlin/tools/os/fsi/fsc"

// name is the *external* path or filename.
func (fs *aeFileSys) SplitX(name string) (dir, bname string) {

	return fsc.UnixPather(name, fs.RootDir())

}
