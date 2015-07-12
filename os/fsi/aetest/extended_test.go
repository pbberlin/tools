// +build extended
// go test -tags=extended

package aetest

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/aefs"
	"github.com/pbberlin/tools/os/fsi/memfs"
	"github.com/pbberlin/tools/os/fsi/osfs"
)

func TestWalk(t *testing.T) {

	bb := new(bytes.Buffer)
	_ = bb
	os.Chdir("c:\\temp")
	pwd, _ := os.Getwd()

	//
	var fs fsi.FileSystem
	if false {
		fsc := aefs.NewAeFs(aefs.MountPointNext(), aefs.AeContext(c))
		fs = fsi.FileSystem(fsc)
	} else if false {
		fsc := &osfs.OsFileSys{}
		fs = fsi.FileSystem(fsc)
	} else {
		fsc := &memfs.MemMapFs{}
		fs = fsi.FileSystem(fsc)
	}

	wpf(os.Stdout, "-----created fs %v %v-----\n", aefs.MountPointLast(), pwd)

	bb = aefs.CreateSys(fs)
	// wpf(os.Stdout, bb.String())

	bb = aefs.RetrieveByReadDir(fs)
	// wpf(os.Stdout, bb.String())

	bb = aefs.RetrieveByQuery(fs)
	// wpf(os.Stdout, bb.String())

	bb = aefs.WalkDirs(fs)
	// wpf(os.Stdout, bb.String())

	bb = aefs.RemoveSubtree(fs)

	time.Sleep(25 * time.Millisecond) // time to have index removals finished

	bb = aefs.WalkDirs(fs)
	// wpf(os.Stdout, bb.String())

	bb = aefs.RetrieveByQuery(fs)

	bb = aefs.RetrieveByReadDir(fs)

	bb = aefs.WalkDirs(fs)

}
