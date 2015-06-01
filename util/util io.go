package util

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

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

func WriteBytesToFilename(filename string, ptrB *bytes.Buffer) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error opening file %v %v\n", filename, err)
	}
	defer f.Close()

	nBytes, err := f.Write(ptrB.Bytes())

	if err != nil {
		fmt.Printf("error writing bytes to %v: %v\n", filename, err)
	}
	fmt.Printf("wrote %d bytes to %v \n", nBytes, filename)
}