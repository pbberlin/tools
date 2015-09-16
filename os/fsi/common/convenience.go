package common

import (
	"os"

	"github.com/pbberlin/tools/os/fsi"
)

func ReadFile(fs fsi.FileSystem, fn string) ([]byte, error) {

	b, err := fs.ReadFile(fn)
	if err != nil {
		return b, err
	}
	return b, nil

}
func WriteFile(fs fsi.FileSystem, fn string, b []byte) error {

	dir, _ := fs.SplitX(fn)

	err := fs.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	err = fs.WriteFile(fn, b, 0)
	if err != nil {
		return err
	}

	return nil
}
