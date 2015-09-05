package osfs

import (
	"path"
	"path/filepath"

	"github.com/pbberlin/tools/os/fsi/common"
)

// name is the *external* path or filename.
func (fs *osFileSys) SplitX(name string) (dir, bname string) {

	// return common.UnixPather(name, fs.RootDir())

	dir, bname = common.UnixPather(name, "")

	name = fs.WinGoofify(path.Join(dir, bname))
	dir, bname = filepath.Split(name)

	return
}
