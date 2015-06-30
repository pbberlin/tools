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

	files, err := fs.GetFiles(path)
	if err != nil {
		return bytes.NewReader(b), err
	}
	b = files[0].Content
	return bytes.NewReader(b), nil
}
