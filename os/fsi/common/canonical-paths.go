package common

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
func cleanseLeadingAndDoublySlashes(p string) string {

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
func UnixPather(name, rootDir string) (dir, bname string) {

	rootName := rootDir
	if len(rootDir) > 1 {
		rootName = rootDir[:len(rootDir)-1]
	}

	name = strings.Replace(name, "\\", "/", -1) // windows to unix; drive letters become directories

	isDirSuffix := ""
	if len(name) > 1 && strings.HasSuffix(name, "/") {
		isDirSuffix = "/"
	}

	name = cleanseLeadingAndDoublySlashes(name)

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

	if len(bname) > 1 {
		bname += isDirSuffix
	}

	return

}

// Remove trailing slash.
// Except from Root.
// We need this when we want Dir-Objects to have a basename
// without trailing slash.
// We also need it, to modify search path towards no-trailing slash.
// In all other cases, SplitX yields exactly the naming we want.
func Filify(bname string) string {
	if len(bname) > 1 {
		return strings.TrimSuffix(bname, "/")
	}
	if bname == "/" {
		return ""
	}
	if bname == "." {
		return ""
	}
	if bname == "" {
		return ""
	}
	return bname
}

// Basically, add trailing slash
func Directorify(name string) string {

	if name == "/" || name == "" || name == "." {
		return name
	}

	if strings.HasSuffix(name, "/") {
		return name
	}

	return name + "/"
}
