// +build extended
// go test -tags=extended

package aetest

import (
	"bytes"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/pbberlin/tools/os/fsi/aefs"
)

func TestWalk(t *testing.T) {

	defer c.Close()

	testRoot := "c:\\temp"
	if runtime.GOOS != "windows" {
		testRoot = "/tmp"
	}

	bb := new(bytes.Buffer)
	msg := ""
	_ = bb

	os.Chdir(testRoot)
	pwd, _ := os.Getwd()
	if pwd == testRoot {
		os.RemoveAll(pwd)
	}

	for _, fs := range Fss {

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
}
