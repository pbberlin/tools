// +build walk
// go test -tags=walk

package aetest

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/pbberlin/tools/os/fsi/aefs"
	"github.com/pbberlin/tools/os/fsi/fsc"
	"github.com/pbberlin/tools/util"
)

func TestWalk(t *testing.T) {

	// first the per-node func:
	exWalkFunc := func(path string, f os.FileInfo, err error) error {

		if err != nil {
			fmt.Printf("Visiting path %s => error: %s \n", path, err)
			return err
		}

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
	_ = exWalkFunc

	rt := <-util.Counter
	rts := fmt.Sprintf("mnt%02v", rt)
	bb := aefs.CreateSys(c, rts)
	wpf(os.Stdout, bb.String())

	bb = aefs.RetrieveDirs(c, rts)
	wpf(os.Stdout, bb.String())

	// err = fsc.Walk(fsi, "temp", exWalkFunc)
	// fmt.Printf("Walk() returned: %v\n", err)
}
