// +build walk
// go test -tags=walk

package aetest

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/aefs"
	"github.com/pbberlin/tools/os/fsi/fsc"
)

func TestWalk(t *testing.T) {

	// first the per-node func:
	exWalkFunc := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Visiting path %s => error: %s \n", path, err)
			return err
		}
		if strings.HasSuffix(path, "my secret directory") {
			return fsc.SkipDir // do not delve deeper
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

	bb := aefs.CreateSys(c)
	wpf(os.Stdout, bb.String())

	// bb = aefs.RetrieveByReadDir(c)
	// wpf(os.Stdout, bb.String())

	// bb = aefs.RetrieveByQuery(c)
	// wpf(os.Stdout, bb.String())

	// bb = aefs.WalkDirs(c)
	// wpf(os.Stdout, bb.String())

	// bb = aefs.WalkDirs(c)
	// wpf(os.Stdout, bb.String())

	fs := aefs.NewAeFs(aefs.MountPointLast(), aefs.AeContext(c))
	fsIn := fsi.FileSystem(fs)
	err := fsc.Walk(fsIn, fs.RootDir(), exWalkFunc)
	log.Printf("Walk() returned: %v\n", err)
}
