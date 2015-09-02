package osutilpb

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//  for tests, this is something like
// C:\Users\pbberlin\AppData\Local\Temp\go-build722556939\github.com\...
func DirOfExecutable() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(dir)
	return dir
}

func GetFilesByExtension(dir, dotExtension string, verbose bool) []string {

	ret := make([]string, 0)

	if dir == "" {
		dir = "."
	}
	dirname := dir + string(filepath.Separator)

	//  a far shorter way seems to be
	//	files, err := filepath.Glob(dirname + "test*.html")

	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("scanning dir %v\n", dirname)
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == dotExtension {
				if verbose {
					fmt.Printf("  found file %v \n", file.Name())
				}
				ret = append(ret, file.Name())
			}
		}
	}

	return ret
}

// PathDirReverse is the opposite of path.Dir(filepath)
// It returns the *first* dir of filepath; not the last
//
// Relative paths are converted to root
// dir1/dir2  => /dir1, /dir2
func PathDirReverse(filepath string) (dir, remainder string, dirs []string) {

	filepath = strings.Replace(filepath, "\\", "/", -1)

	filepath = path.Join(filepath, "")

	if filepath == "/" || filepath == "" || filepath == "." {
		return "/", "", []string{}
	}

	filepath = strings.TrimPrefix(filepath, "/")

	dirs = strings.Split(filepath, "/")

	if len(dirs) == 1 {
		return "/" + dirs[0], "", []string{}
	} else {
		return "/" + dirs[0], "/" + strings.Join(dirs[1:], "/"), dirs[1:]
	}

	return

}
