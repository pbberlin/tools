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
	msg := ""
	_ = bb
	os.Chdir("c:\\temp")
	pwd, _ := os.Getwd()
	_ = pwd

	//
	var fs fsi.FileSystem
	if true {
		fsc := aefs.New(aefs.MountPointNext(), aefs.AeContext(c))
		fs = fsi.FileSystem(fsc)
	} else if false {
		fsc := osfs.New()
		fs = fsi.FileSystem(fsc)
	} else {
		fsc := memfs.New()
		fs = fsi.FileSystem(fsc)
	}

	wpf(os.Stdout, "-----created fs %v %v-----\n", fs.Name(), fs.String())

	bb, msg = aefs.CreateSys(fs)
	if msg != "" {
		wpf(os.Stdout, msg+"\n")
		wpf(os.Stdout, bb.String())
	}

	bb, msg = aefs.RetrieveByReadDir(fs)
	if msg != "" {
		wpf(os.Stdout, msg+"\n")
		wpf(os.Stdout, bb.String())
	}

	bb, msg = aefs.RetrieveByQuery(fs)
	if msg != "" {
		wpf(os.Stdout, msg+"\n")
		wpf(os.Stdout, bb.String())
	}

	bb, msg = aefs.WalkDirs(fs)
	if msg != "" {
		wpf(os.Stdout, msg+"\n")
		wpf(os.Stdout, bb.String())
	}

	bb, msg = aefs.RemoveSubtree(fs)
	if msg != "" {
		wpf(os.Stdout, msg+"\n")
		wpf(os.Stdout, bb.String())
	}

	// After removal, give time,
	// to remove directories from index too.
	// Alternatively, the walkFunc should not return
	// err == datastore.ErrNoSuchEntity
	time.Sleep(5 * time.Millisecond)

	bb, msg = aefs.WalkDirs(fs)
	if msg != "" {
		wpf(os.Stdout, msg+"\n")
		wpf(os.Stdout, bb.String())
	}

	bb, msg = aefs.RetrieveByQuery(fs)
	if msg != "" {
		wpf(os.Stdout, msg+"\n")
		wpf(os.Stdout, bb.String())
	}

	bb, msg = aefs.RetrieveByReadDir(fs)
	if msg != "" {
		wpf(os.Stdout, msg+"\n")
		wpf(os.Stdout, bb.String())
	}

	bb, msg = aefs.WalkDirs(fs)
	if msg != "" {
		wpf(os.Stdout, msg+"\n")
		wpf(os.Stdout, bb.String())
	}

}
