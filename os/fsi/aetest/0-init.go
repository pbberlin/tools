// Package aetest runs test, which require an appengine context;
// and which must be run by goapp test
package aetest

import (
	"log"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/aefs"
	"github.com/pbberlin/tools/os/fsi/aefs_sr"
	"github.com/pbberlin/tools/os/fsi/memfs"
	"github.com/pbberlin/tools/os/fsi/osfs"

	"appengine/aetest"
)

var dot = []string{
	"fs.go",
	"fs_test.go",
	"httpFs.go",
	"memfile.go",
	"memmap.go",
}

var testDir = "/temp/fun"
var testName = "testF.txt"

var Fss = []fsi.FileSystem{}
var c aetest.Context

func init() {

	var err error
	c, err = aetest.NewContext(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Do not
	// defer c.Close()

	// We can have variadic option funcs.
	// But seems we can not make those generic,
	// since we need the concrete filesystem type
	// one way or another.

	fs1 := aefs.NewAeFs("rootX", aefs.AeContext(c))
	fs1i := fsi.FileSystem(fs1)
	_ = fs1i

	fs2 := aefs_sr.NewAeFs("rootY", aefs_sr.AeContext(c))
	fs2i := fsi.FileSystem(fs2)
	_ = fs2i

	fs3 := &osfs.OsFileSys{}
	fs3i := fsi.FileSystem(fs3)
	_ = fs3i

	fs4 := &memfs.MemMapFs{}
	fs4i := fsi.FileSystem(fs4)
	_ = fs4i

	Fss = []fsi.FileSystem{fs1i, fs2i, fs3i, fs4i}
	// fss := []fsi.FileSystem{fs4i}

}
