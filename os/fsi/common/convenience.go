package common

import (
	"os"

	"github.com/pbberlin/tools/os/fsi"
)

func WriteFile(fs fsi.FileSystem, fn string, b []byte) error {

	dir, _ := fs.SplitX(fn)

	err := fs.MkdirAll(dir, os.ModePerm)
	if err != nil && err != fsi.ErrFileExists {
		return err
	}

	err = fs.WriteFile(fn, b, 0)
	if err != nil {
		return err
	}

	return nil
}
