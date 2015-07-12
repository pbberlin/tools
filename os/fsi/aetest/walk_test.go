// +build walk
// go test -tags=walk

package aetest

import (
	"bytes"
	"os"
	"testing"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/aefs"
)

func TestWalk(t *testing.T) {

	fsConcrete := aefs.NewAeFs(aefs.MountPointNext(), aefs.AeContext(c))
	fs := fsi.FileSystem(fsConcrete)

	bb := new(bytes.Buffer)
	wpf(bb, "created fs %v\n\n", aefs.MountPointLast())

	bb = aefs.CreateSys(fs)
	wpf(os.Stdout, bb.String())

	bb = aefs.RetrieveByReadDir(fs)
	wpf(os.Stdout, bb.String())

	bb = aefs.RetrieveByQuery(fs)
	wpf(os.Stdout, bb.String())

	bb = aefs.WalkDirs(fs)
	wpf(os.Stdout, bb.String())

	bb = aefs.RemoveSubtree(fs)
	wpf(os.Stdout, bb.String())

}
