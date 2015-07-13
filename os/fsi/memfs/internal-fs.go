package memfs

import (
	"log"
	"path"
	filepath "path"
	"strings"
	// "path/filepath"

	"github.com/pbberlin/tools/os/fsi"
)

// try to find parent in main map-of-files
func (m *memMapFs) findParent(f fsi.File) fsi.File {

	dir, _ := path.Split(f.Name())
	dir = cleanseLeadingSlash(dir)

	if len(dir) > 1 {
		parent := ""
		if false { // original implementation
			_, parent = path.Split(path.Clean(dir))
		} else {
			parent = path.Clean(dir)
		}
		// log.Printf("%-22v has parent %v ", f.Name(), parent)
		if len(parent) > 0 {
			pfile, err := m.Open(parent)
			if err != nil {
				// log.Printf("parent %-8q not found for %v", parent, f.Name())
			} else {
				// log.Printf("parent %-8q found     for %v", parent, f.Name())
				return pfile
			}
		}
	}
	return nil
}

// adding f to it's parent's directory
func (m *memMapFs) registerWithParent(f fsi.File) fsi.File {
	if f == nil {
		return nil
	}

	parent := m.findParent(f)

	if parent == nil {
		cleanName := cleanseLeadingSlash(f.Name())
		pdir, base := filepath.Split(cleanName)
		pdir = cleanseLeadingSlash(pdir)
		_ = base
		if len(pdir) > 0 {
			log.Printf("  create parent %-24q for %v\n", pdir, cleanName)
			err := m.MkdirAll(pdir, 0777)
			if err != nil && err != fsi.ErrFileExists {
				log.Printf("Mkdir failed %v", err)
			}
		}
	}

	parent = m.findParent(f) // trying again
	if parent != nil {
		pmem := parent.(*InMemoryFile)
		// log.Printf("  added\n")
		pmem.memDir.Add(f)

		log.Printf("    par %-32q got added %32q\n", parent.Name(), f.Name())

	}

	return parent
}

// recursing upwards
func (m *memMapFs) registerDirs(f fsi.File) {
	var x = f.Name()
	x = cleanseLeadingSlash(x)
	cntr := 0
	for {
		log.Printf("register loop %v %v\n", cntr, x)
		cntr++
		if cntr > 4 {
			break
		}
		f = m.registerWithParent(f)
		if f == nil {
			break
		}
		x = f.Name()
	}
}

// remove f from it's parent's directory
func (m *memMapFs) unRegisterWithParent(f fsi.File) fsi.File {
	parent := m.findParent(f)
	pmem := parent.(*InMemoryFile)
	pmem.memDir.Remove(f)
	return parent
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
