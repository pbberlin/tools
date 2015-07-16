// +build extended
// go test -tags=extended

package fstest

import (
	"bytes"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestWalk(t *testing.T) {

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

	Fss, c := initFileSystems()
	defer c.Close()

	for _, fs := range Fss {

		wpf(os.Stdout, "-----created fs %v %v-----\n", fs.Name(), fs.String())

		bb, msg = CreateSys(fs)
		if msg != "" {
			wpf(os.Stdout, msg+"\n")
			wpf(os.Stdout, bb.String())
			t.Errorf("%v", msg)
		}

		bb, msg = RetrieveByReadDir(fs)
		if msg != "" {
			wpf(os.Stdout, msg+"\n")
			wpf(os.Stdout, bb.String())
			t.Errorf("%v", msg)
		}

		bb, msg = RetrieveByQuery(fs)
		if msg != "" {
			wpf(os.Stdout, msg+"\n")
			wpf(os.Stdout, bb.String())
			t.Errorf("%v", msg)
		}

		bb, msg = WalkDirs(fs)
		if msg != "" {
			wpf(os.Stdout, msg+"\n")
			wpf(os.Stdout, bb.String())
			t.Errorf("%v", msg)
		}

		bb, msg = RemoveSubtree(fs)
		if msg != "" {
			wpf(os.Stdout, msg+"\n")
			wpf(os.Stdout, bb.String())
			t.Errorf("%v", msg)
		}

		// After removal, give time,
		// to remove directories from index too.
		// Alternatively, the walkFunc should not return
		// err == datastore.ErrNoSuchEntity
		time.Sleep(5 * time.Millisecond)

		bb, msg = WalkDirs(fs)
		if msg != "" {
			wpf(os.Stdout, msg+"\n")
			wpf(os.Stdout, bb.String())
			t.Errorf("%v", msg)
		}

		bb, msg = RetrieveByQuery(fs)
		if msg != "" {
			wpf(os.Stdout, msg+"\n")
			wpf(os.Stdout, bb.String())
			t.Errorf("%v", msg)
		}

		bb, msg = RetrieveByReadDir(fs)
		if msg != "" {
			wpf(os.Stdout, msg+"\n")
			wpf(os.Stdout, bb.String())
			t.Errorf("%v", msg)
		}

		bb, msg = WalkDirs(fs)
		if msg != "" {
			wpf(os.Stdout, msg+"\n")
			wpf(os.Stdout, bb.String())
			t.Errorf("%v", msg)
		}
	}
}
