package gaefs

import (
	"bytes"
	"io"
)

func (fs FileSys) Mkdir(path string) {
	panic("not implemented")
}

func (fs FileSys) Touch(p string) {
	panic("not implemented")
}

// TODO: !!!
// must return ReadSeekCloser
// and much more
func (fs FileSys) Open(path string) (io.Reader, error) {

	var b []byte

	file, err := fs.GetFile(path)
	if err != nil {
		return bytes.NewReader(b), err
	}

	b = file.Content
	return bytes.NewReader(b), nil
}
