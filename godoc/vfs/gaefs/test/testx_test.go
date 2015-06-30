package test

import (
	"io"
	"log"
	"os"
	"testing"

	"appengine/aetest"

	"github.com/pbberlin/tools/godoc/vfs/gaefs"
)

func Test1(t *testing.T) {

	c, err := aetest.NewContext(nil)
	if err != nil {
		log.Printf("%v\n", err)
		t.Fatal(err)
	}
	defer c.Close()

	fs := gaefs.NewFs("rootX", c, false)

	// note: paths must start with a slash!
	f, err := fs.Open("/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	// defer f.Close()
	io.Copy(os.Stdout, f)
}
