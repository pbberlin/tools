// +build pathing
// go test -tags=pathing

package test

import (
	"fmt"
	"log"
	"path"
	"testing"

	"github.com/pbberlin/tools/os/fsi/common"
	"github.com/pbberlin/tools/os/fsi/memfs"
)

func splitIsWhatWeWant() {

	s := "rt/"

	dir1, f1 := path.Split(s)

	dir2 := path.Dir(s)
	f2 := path.Base(s)

	fmt.Printf("%q %q \n", dir1, f1) //  "rt/"   ""   - good
	fmt.Printf("%q %q \n", dir2, f2) //  "rt"    "rt" - bad

}

func TestPathCleanage(t *testing.T) {

	_, c := initFileSystems()
	defer c.Close()

	cases := [][]string{
		[]string{"", "mntX/", ""},
		[]string{"/", "mntX/", ""},
		[]string{".", "mntX/", ""},
		[]string{"mntX", "mntX/", ""},
		[]string{"mntX/", "mntX/", ""},
		[]string{"mntX/dir1", "mntX/", "dir1"},
		[]string{"mntX/dir1/file2", "mntX/dir1/", "file2"},
		[]string{"///mntX/dir1/dir1///file3/", "mntX/dir1/dir2", "file3"},
		[]string{"mntX/dir1/dir1///file3/", "mntX/dir1/dir2", "file3"},
		[]string{"/dir1/dir1///file3/", "mntX/dir1/dir2", "file3"},
		[]string{"dir1/dir1///file3/", "mntX/dir1/dir2", "file3"},
		[]string{"./dir1/", "mntX/", "dir1"},
		[]string{"c:\\dir1\\dir2", "mntX/", "mntX/c:/dir1/dir2"},
	}

	fs := memfs.New(memfs.MountName("mntX"))
	for _, v := range cases {
		inpt := v[0]
		_ = inpt
		wnt1 := v[1]
		wnt2 := v[2]
		dir, bname := common.UnixPather(v[0], fs.RootDir())
		fullpath := dir + bname

		log.Printf("%-28v %-24v => %-16q %-12q ", inpt, dir, bname, fullpath)

		if wnt1 != dir {
			t.Logf("got %-13v - wnt %v\n", dir, wnt1)
		}
		if wnt2 != bname {
			t.Logf("got %-13v - wnt %v\n", bname, wnt2)
		}
		if wnt1+wnt2 != fullpath {
			t.Logf("got %-13v - wnt %v\n", fullpath, wnt1+wnt2)
		}
	}

}
