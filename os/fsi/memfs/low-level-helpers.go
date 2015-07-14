package memfs

import "github.com/pbberlin/tools/os/fsi/fsc"

// name is the *external* path or filename.
func (fs *memMapFs) pathInternalize(name string) (dir, bname string) {

	return fsc.PathInternalize(name, fs.RootDir(), fs.RootName())

}
