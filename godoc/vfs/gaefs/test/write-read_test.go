package test

import (
	"io"
	"log"
	"os"
	"testing"

	"appengine/aetest"

	"github.com/pbberlin/tools/godoc/vfs/gaefs"
)

func TestWriteRead(t *testing.T) {

	c, err := aetest.NewContext(nil)
	if err != nil {
		log.Printf("%v\n", err)
		t.Fatal(err)
	}
	defer c.Close()

	fs := gaefs.NewFs("rootX", c, false)

	dir, err := fs.SaveDirByPath("/xx")
	_ = dir
	if err != nil {
		log.Fatal(err)
	}

	f, err := fs.Create("xx/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString("some text content")
	f.Close()
	err = fs.SaveFile(&f, "xx")
	if err != nil {
		log.Fatal(err)
	}

	f2, err := fs.Open("xx/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()
	io.Copy(os.Stdout, &f2)

	bts, err := gaefs.ReadFile(&fs, "xx/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	// defer rdr.Close()
	log.Printf("2nd: %v", string(bts))

}
