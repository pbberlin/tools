package gaefs

import (
	"bytes"

	"golang.org/x/tools/godoc/vfs"
)

func (fs FileSys) Mkdir(path string) {
	panic("not implemented")
}

func (fs FileSys) Touch(p string) {
	panic("not implemented")
}

func OS(mount string) FileSys {
	panic(`
		Sadly, google app engine file system requires a
	 	http.Request based context object.
	 	Use NewFs(string, appengine.Context) instead of OS.
	`)
}

func (fs FileSys) Open(path string) (vfs.ReadSeekCloser, error) {

	var b []byte

	file, err := fs.GetFile(path)
	if err != nil {
		return NopCloser(bytes.NewReader(b)), err
	}

	b = file.Content
	return NopCloser(bytes.NewReader(b)), nil
}
