// +build pathing
// go test -tags=pathing

package tests

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

	log.SetFlags(0)

	_, c := initFileSystems()
	defer c.Close()

	cases := [][]string{
		[]string{"", "mntX/", ""},
		[]string{"/", "mntX/", ""},
		[]string{".", "mntX/", ""},
		[]string{"./", "mntX/", ""},
		[]string{"mntX", "mntX/", ""},
		[]string{"mntX/", "mntX/", ""},
		// 5...
		[]string{"file", "mntX/", "file"},
		[]string{"dir/", "mntX/", "dir/"},
		[]string{"./dir1/", "mntX/", "dir1/"},
		[]string{"mntX/dir1", "mntX/", "dir1"},
		[]string{"mntX/dir1/file2", "mntX/dir1/", "file2"},
		// 10
		[]string{"mntY/dir1/file2", "mntX/mntY/dir1/", "file2"},
		[]string{"///mntX/dir1/dir2///dir3/", "mntX/dir1/dir2/", "dir3/"},
		[]string{"mntX/dir1/dir2///file3", "mntX/dir1/dir2/", "file3"},
		[]string{"/dir1/dir2///file3", "mntX/dir1/dir2/", "file3"},
		[]string{"dir1/dir2///dir3/", "mntX/dir1/dir2/", "dir3/"},
		// 15
		[]string{"c:\\dir1\\dir2", "mntX/c:/dir1/", "dir2"},
	}

	fs := memfs.New(
		memfs.Ident("mntX"),
	)
	for i, v := range cases {
		inpt := v[0]
		_ = inpt
		wnt1 := v[1]
		wnt2 := v[2]
		dir, bname := common.UnixPather(v[0], fs.RootDir())
		fullpath := dir + bname

		log.Printf("#%2v %-30q %-24q => %-16q %-12q ", i, inpt, dir, bname, fullpath)

		if wnt1 != dir {
			t.Errorf("dir   #%2v got %-24v - wnt %-16v\n", i, dir, wnt1)
		}
		if wnt2 != bname {
			t.Errorf("bname #%2v got %-24v - wnt %-16v\n", i, bname, wnt2)
		}
		if wnt1+wnt2 != fullpath {
			t.Errorf("fullp #%2v got %-24v - wnt %-16v\n", i, fullpath, wnt1+wnt2)
		}
	}

}
