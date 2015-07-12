// +build extended
// go test -tags=extended

package aetest

import (
	"bytes"
	"os"
	"testing"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/aefs"
)

func TestWalk(t *testing.T) {

	bb := new(bytes.Buffer)

	fsConcrete := aefs.NewAeFs(aefs.MountPointNext(), aefs.AeContext(c))
	fs := fsi.FileSystem(fsConcrete)
	wpf(bb, "created fs %v\n\n", aefs.MountPointLast())

	// fsConcrete := &osfs.OsFileSys{}
	// fs := fsi.FileSystem(fsConcrete)

	os.Chdir("c:\\temp")
	pwd, _ := os.Getwd()
	wpf(os.Stdout, "-----created fs %v %v-----\n", "mntX", pwd)

	bb = aefs.CreateSys(fs)
	wpf(os.Stdout, bb.String())

	// bb = aefs.RetrieveByReadDir(fs)
	// wpf(os.Stdout, bb.String())

	// bb = aefs.RetrieveByQuery(fs)
	// wpf(os.Stdout, bb.String())

	// bb = aefs.WalkDirs(fs)
	// wpf(os.Stdout, bb.String())

	// bb = aefs.RemoveSubtree(fs)
	// wpf(os.Stdout, bb.String())

}
