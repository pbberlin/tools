package osutilpb

import "testing"

var inputWant = [][]string{
	[]string{"", "/", ""},
	[]string{"/dir1", "/dir1", ""},
	[]string{"/dir1/", "/dir1", ""},
	[]string{"/dir1/dir2", "/dir1", "/dir2"},
	[]string{"dir1/dir2", "/dir1", "/dir2"},
	[]string{"/dir1/dir2/file.html", "/dir1", "/dir2/file.html"},

	[]string{"/dir2/file.html", "/dir2", "/file.html"},
}

func Test1(t *testing.T) {
	for k, v := range inputWant {
		got1, got2 := PathDirReverse(v[0])

		if got1 != v[1] || got2 != v[2] {
			t.Errorf("Failed test #%2v: inp: %-12v got %q %q - want %q %q\n", k, v[0], got1, got2, v[1], v[2])

		}
	}
}
