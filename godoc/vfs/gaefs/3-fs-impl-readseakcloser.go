package gaefs

import (
	"bytes"
	"io"
	"io/ioutil"

	"golang.org/x/tools/godoc/vfs"
)

// This is inspired by ioutil.NopCloser.
type nopCloser struct {
	io.ReadSeeker
}

func (nopCloser) Close() error { return nil }

// NopCloser wraps a ReadCloser
// adding a no-op Close method.
//
// The fascinating thing is,
// that vfs.ReadSeekCloser demands an io.Closer
// but is content with a provided gaefs.Close() method.
// It does not matter in which package an interface
// is defined - and where it's implementations come from.
func NopCloser(r io.ReadSeeker) vfs.ReadSeekCloser {
	return nopCloser{r}
}

// satisfy vfs.Opener
// Conflicts with Afero Open method
func (fs AeFileSys) OpenVFS(path string) (vfs.ReadSeekCloser, error) {

	var b []byte

	file, err := fs.GetFile(path)
	if err != nil {
		return NopCloser(bytes.NewReader(b)), err
	}

	b = file.data
	return NopCloser(bytes.NewReader(b)), nil
}

func ReadFileVFS(fs AeFileSys, path string) ([]byte, error) {

	// via read-seak-closer:
	rsc, err := fs.OpenVFS(path)
	if err != nil {
		return []byte{}, err
	}

	defer rsc.Close()
	b, err := ioutil.ReadAll(rsc)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}
