package osfs

import (
	"os"
	"runtime"
	"sort"

	"github.com/pbberlin/tools/os/fsi"
)

type osFileSys struct {
	replacePath   bool
	readdirsorter func([]os.FileInfo)
}

func New(options ...func(fsi.FileSystem)) *osFileSys {

	repl := false
	if runtime.GOOS == "windows" {
		repl = true
	}

	o := &osFileSys{}
	o.replacePath = repl
	o.readdirsorter = func(fis []os.FileInfo) {} // unchanged

	for _, option := range options {
		option(o) // apply options over defaults
	}

	return o
}

// Default sort for ReadDir... is ByNameAsc
// We may want to change this; for instance sort byDate
func DirSort(srt string) func(fsi.FileSystem) {
	return func(fs fsi.FileSystem) {
		fst := fs.(*osFileSys)

		switch srt {
		case "byDateAsc":
			fst.readdirsorter = func(fis []os.FileInfo) {
				sort.Sort(byDateAsc(fis))
			}

		case "byDateDesc":
			fst.readdirsorter = func(fis []os.FileInfo) {
				sort.Sort(byDateDesc(fis))
			}
		case "byName":
			fst.readdirsorter = func(fis []os.FileInfo) {} // unchanged
		}
	}
}
