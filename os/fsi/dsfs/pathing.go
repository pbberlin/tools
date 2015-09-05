package dsfs

import (
	"strings"

	"github.com/pbberlin/tools/os/fsi/common"
)

// name is the *external* path or filename.
func (fs *dsFileSys) SplitX(name string) (dir, bname string) {
	return common.UnixPather(name, fs.RootDir())
}

// Remove trailing slash.
// Except from Root.
func filyfyBName(bname string) string {
	if len(bname) > 1 {
		return strings.TrimSuffix(bname, "/")
	}
	if bname == "/" {
		return ""
	}
	if bname == "." {
		return ""
	}
	if bname == "" {
		return ""
	}
	return bname
}
