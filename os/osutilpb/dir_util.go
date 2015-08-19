package osutilpb

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

//  for tests, this is something like
// C:\Users\pbberlin\AppData\Local\Temp\go-build722556939\github.com\...
func DirOfExecutable() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
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
