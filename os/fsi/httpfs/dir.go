package httpfs

import (
	"errors"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

type httpDir struct {
	basePath string
	fs       HttpFs
}

func (d httpDir) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.Contains(name, "\x00") {
		return nil, errors.New("http: invalid character in file path")
	}
	dir := string(d.basePath)
	if dir == "" {
		dir = "."
	}

	f, err := d.fs.Open(filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name))))
	if err != nil {
		return nil, err
	}
	return f, nil
}
