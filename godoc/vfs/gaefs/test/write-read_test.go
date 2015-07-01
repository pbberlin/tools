package test

import (
	"fmt"
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

	f := gaefs.File{}
	f.Name = "test.txt"
	f.Content = []byte("\tsome text content\n")
	err = fs.SaveFile(&f, "/xx")
	if err != nil {
		log.Fatal(err)
	}

	rdr, err := fs.Open("xx/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	// defer rdr.Close()
	io.Copy(os.Stdout, rdr)

	bts, err := gaefs.ReadFile(fs, "xx/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	// defer rdr.Close()
	fmt.Printf("2nd: %v", string(bts))

}
