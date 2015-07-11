// +build cleaning
// go test -tags=cleaning

package aefs

import (
	"fmt"
	"log"
	"path"
	"testing"
)

func splitIsWhatWeWant() {

	s := "rt/"

	dir1, f1 := path.Split(s)

	dir2 := path.Dir(s)
	f2 := path.Base(s)

	fmt.Printf("%q %q \n", dir1, f1) //  "rt/"   ""
	fmt.Printf("%q %q \n", dir2, f2) //  "rt"    "rt"

}

func TestPathCleanage(t *testing.T) {

	cases := [][]string{
		[]string{"", "mntX/", ""},
		[]string{"/", "mntX/", ""},
		[]string{"mntX", "mntX/", ""},
		[]string{"mntX/", "mntX/", ""},
		[]string{"mntX/dir1", "mntX/", "dir1"},
		[]string{"mntX/dir1/file2", "mntX/dir1/", "file2"},
		[]string{"///mntX/dir1/dir1///file3/", "mntX/dir1/dir2", "file3"},
		[]string{"mntX/dir1/dir1///file3/", "mntX/dir1/dir2", "file3"},
		[]string{"/dir1/dir1///file3/", "mntX/dir1/dir2", "file3"},
		[]string{"dir1/dir1///file3/", "mntX/dir1/dir2", "file3"},
	}

	fs := AeFileSys{mount: "mntX"}
	for _, v := range cases {
		inpt := v[0]
		_ = inpt
		wnt1 := v[1]
		wnt2 := v[2]
		got1, got2, got3 := fs.pathInternalize(v[0])

		log.Printf("%-28v %-24v => %-16q %-12q ", inpt, got1, got2, got3)

		if wnt1 != got1 {
			t.Logf("got %-13v - wnt %v\n", got1, wnt1)
		}
		if wnt2 != got2 {
			t.Logf("got %-13v - wnt %v\n", got2, wnt2)
		}
		if wnt1+wnt2 != got3 {
			t.Logf("got %-13v - wnt %v\n", got3, wnt1+wnt2)
		}
	}
}
