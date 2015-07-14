// +build suite0
// go test -tags=suite0

package aetest

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestWriteRead(t *testing.T) {

	for _, fs := range Fss {

		c.Infof(" ")
		c.Infof("%v %v\n", fs.Name(), fs.String())
		c.Infof("================================")

		err := fs.Mkdir("/temp/testdir", os.ModePerm)
		if err != nil {
			if !os.IsExist(err) {
				c.Criticalf("%v\n", err)
				t.Fatalf("%v\n", err)
			}
		}

		f, err := fs.Create("/temp/testdir/test.txt")
		if err != nil {
			c.Criticalf("create: %v\n", err)
			t.Fatalf("create: %v\n", err)
		}

		bts0 := []byte("some text content")
		_, err = f.WriteString(string(bts0))
		if err != nil {
			c.Criticalf("writestr: %v\n", err)
			t.Fatalf("writestr: %v\n", err)
		}


		err = f.Close()
		if err != nil {
			c.Criticalf("close: %v\n", err)
			t.Fatalf("close: %v\n", err)
		}

		f2, err := fs.Open("/temp/testdir/test.txt")
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
		// c.Infof("1st: %v\n", string(bts1))

		bts2, err := fs.ReadFile("/temp/testdir/test.txt")
		if err != nil {
			c.Criticalf("%v\n", err)
			t.Fatalf("%v\n", err)
		}

		// c.Infof("2nd: %v\n", string(bts2))

		bts4 := []byte("other stuff")
		err = fs.WriteFile("/temp/testdir/test1.txt", bts4, os.ModePerm)
		if err != nil {
			c.Criticalf("%v\n", err)
			t.Fatalf("%v\n", err)
		}

		bts5, err := fs.ReadFile("/temp/testdir/test1.txt")
		if err != nil {
			c.Criticalf("%v\n", err)
			t.Fatalf("%v\n", err)
		}

		if !(string(bts0) == string(bts1) && string(bts1) == string(bts2)) {
			t.Fatalf("there are differences \nwnt %s \ngt1 %s \ngt2 %s", bts0, bts1, bts2)
		}
		if !(string(bts4) == string(bts5)) {
			t.Fatalf("there are differences \nwnt %s \ngt1 %s", bts4, bts5)
		}

		// fsc, ok := fs.(*memfs.MemMapFs)
		// if ok {
		// 	fsc.Dump()
		// }

	}

}
