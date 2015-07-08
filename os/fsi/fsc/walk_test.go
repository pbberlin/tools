package fsc

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/gaefs"
)

func TestWalk(t *testing.T) {

	// first the per-node func:
	exWalkFunc := func(path string, f os.FileInfo, err error) error {

		if strings.HasSuffix(path, "my secret directory") {
			return SkipDir
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

	mnt := "mnt00"
	fs := gaefs.NewAeFs(mnt) // add appengine context
	fsiX := fsi.FileSystem(fs)

	err := fsiX.Walk(mnt, exWalkFunc)
	fmt.Printf("Walk() returned %v\n", err)
}
