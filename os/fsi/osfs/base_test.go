package osfs

import (
	"os"
	"testing"
)

func TestOsFileSys(t *testing.T) {

	fs := osFileSys{}
	_ = fs
	// Run code and tests requiring the appengine.Context using c.

	err := fs.Mkdir("/temp", os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			t.Fatalf("%v\n", err)
		}
	}

	f, err := fs.Create("/temp/test.txt")
	defer f.Close()
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	_, err = f.WriteString("oh, Ashley, oh.")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = os.Remove("/temp/test.txt")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

}
