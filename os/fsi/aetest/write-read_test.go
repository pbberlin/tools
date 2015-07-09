// Package aetest runs test, which require an appengine context;
// and which must be run by goapp test
package aetest

import (
	"io/ioutil"
	"os"
	"testing"

	"appengine/aetest"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/aefs"
	"github.com/pbberlin/tools/os/fsi/aefs_sr"
)

func TestWriteRead(t *testing.T) {

	c, err := aetest.NewContext(nil)
	if err != nil {
		c.Criticalf("%v\n", err)
		t.Fatalf("%v\n", err)
	}
	defer c.Close()

	fs1 := aefs.NewAeFs("rootX", aefs.AeContext(c))
	fs1i := fsi.FileSystem(fs1)

	fs2 := aefs_sr.NewAeFs("rootY", aefs_sr.AeContext(c))
	fs2i := fsi.FileSystem(fs2)

	fss := []fsi.FileSystem{fs1i, fs2i}

	for _, fs := range fss {
		err = fs.Mkdir("/xx", os.ModePerm)
		if err != nil {
			c.Criticalf("%v\n", err)
			t.Fatalf("%v\n", err)
		}

		f, err := fs.Create("xx/test.txt")
		if err != nil {
			c.Criticalf("%v\n", err)
			t.Fatalf("%v\n", err)
		}

		bts0 := []byte("some text content")
		f.WriteString(string(bts0))
		f.Close()

		err = f.Sync()
		if err != nil {
			c.Criticalf("%v\n", err)
			t.Fatalf("%v\n", err)
		}

		f2, err := fs.Open("xx/test.txt")
		if err != nil {
			c.Criticalf("%v\n", err)
			t.Fatalf("%v\n", err)
		}
		defer f2.Close()
		bts1, err := ioutil.ReadAll(f2)
		if err != nil {
			c.Criticalf("%v\n", err)
			t.Fatalf("%v\n", err)
		}
		c.Infof("1st: %v\n", string(bts1))

		bts2, err := fs.ReadFile("xx/test.txt")
		if err != nil {
			c.Criticalf("%v\n", err)
			t.Fatalf("%v\n", err)
		}

		c.Infof("2nd: %v\n", string(bts2))

		if !(string(bts0) == string(bts1) && string(bts1) == string(bts2)) {
			t.Fatalf("there are differences \nwnt%s \ngt1%s \ngt2%s", bts0, bts1, bts2)
		}

	}

}
