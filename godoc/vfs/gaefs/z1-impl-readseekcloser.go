package gaefs

import (
	"io"

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
