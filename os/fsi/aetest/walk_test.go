// +build walk
// go test -tags=walk

package aetest

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/fsc"
	"github.com/pbberlin/tools/os/fsi/memfs"
)

func TestWalk(t *testing.T) {

	// first the per-node func:
	exWalkFunc := func(path string, f os.FileInfo, err error) error {

		if strings.HasSuffix(path, "my secret directory") {
			return fsc.SkipDir
		}

		if err == os.ErrInvalid {
			return err // calling off the walk
		}

		tp := "file"
		if f.IsDir() {
			tp = "dir "
		}
		fmt.Printf("Visited: %s %s \n", tp, path)
		return nil
	}

	fs := &memfs.MemMapFs{}
	fsi := fsi.FileSystem(fs)
	fsi.Mkdir("/", os.ModePerm)
	fsi.Mkdir("/temp", os.ModePerm)

	err := fsc.Walk(fsi, "/", exWalkFunc)
	fmt.Printf("Walk() returned %v\n", err)
}
