package fsc

import (
	"path"
	"strings"
)

const doublesep = "//"
const sep = "/"

// "" 	=> ""
// "/" 	=> ""
// "/aa/bb/" 	=> "aa/bb"
// "aa/bb/" 	=> "aa/bb"
func cleanseLeadingSlash(p string) string {

	// path.Join does not clean spaces
	p = strings.TrimSpace(p)

	// Make any // single /  -  make leading // single /
	// Remove trailing /
	// i.e.: 		"//a\n a//a//"
	// 				"/a\n a/a"
	p = path.Join(p)

	// now remove the remaining possible leading sep
	if strings.HasPrefix(p, sep) {
		p = p[1:]
	}

	return p
}

// name is the *external* path or filename.
func PathInternalize(name, rootDir, rootName string) (dir, bname string) {

	name = cleanseLeadingSlash(name)

	// exchange current dir "." for root
	if strings.HasPrefix(name, ".") {
		name = rootDir + name[1:]
		name = strings.Replace(name, doublesep, sep, -1)
	}

	// prepend rootdir, if neccessary
	if !strings.HasPrefix(name, rootName) {
		name = rootDir + name
	}

	// append slash to rootdir, if alone
	if name == rootName {
		name += sep
	}

	dir, bname = path.Split(name)

	if !strings.HasSuffix(dir, sep) {
		dir += sep
	}

	return

}
