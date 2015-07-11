package aefs

import (
	"path"
	"strings"

	"appengine/datastore"
)

func OBSOLETE_dirFromKey(key *datastore.Key) string {

	dir := key.String()
	dir = cleanseLeadingSlash(dir)

	dir = strings.Replace(dir, tdirsep, "", -1) // removing all fsd,
	dirs := strings.Split(dir, sep)

	if len(dirs) > 1 {
		// dirs = dirs[1:] // WITH Root
		dir = path.Join(dirs...)
		dir = dir + sep
		return dir
	} else {
		return ""
	}

}

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
func (fs *AeFileSys) pathInternalize(name string) (dir, bname string) {

	name = cleanseLeadingSlash(name)

	// prepend rootdir, if neccessary
	if !strings.HasPrefix(name, fs.RootName()) {
		name = fs.RootDir() + name
	}

	// append slash to rootdir, if alone
	if name == fs.RootName() {
		name += sep
	}

	dir, bname = path.Split(name)

	if dir != fs.RootDir() && bname == "" {
		panic("file needs name")
	}

	if !strings.HasSuffix(dir, sep) {
		dir += sep
	}

	return

	// remove filename
	if strings.HasSuffix(name, bname) {
		name = name[:len(name)-len(bname)]
	}
	return

}
