package filesys

import "path/filepath"

func parent(p string) string {
	return (filepath.Dir(p))
}
