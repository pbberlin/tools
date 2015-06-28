package rooted

import "path/filepath"

func parent(path string) string {
	return (filepath.Dir(path))
}
func fileN(path string) string {
	return (filepath.Base(path))
}
