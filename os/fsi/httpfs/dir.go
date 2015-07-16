package httpfs

import (
	"errors"
	"net/http"
	"path"
	filepath "path" //"path/filepath"
	"strings"
)

type httpDir struct {
	basePath string
	fs       HttpFs
}

func (d httpDir) Open(name string) (http.File, error) {
	if strings.Contains(name, "\x00") {
		return nil, errors.New("http: invalid character in file path")
	}
	dir := string(d.basePath)
	if dir == "" {
		dir = "."
	}

	// jpath := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	jpath := filepath.Join(dir, path.Clean("/"+name))
	f, err := d.fs.Open(jpath)
	if err != nil {
		// log.Printf("    httpdir open %-22v Err %v", jpath, err)
		return nil, err
	}
	// log.Printf("    httpdir open %-22v Success", jpath)
	return f, nil
}
