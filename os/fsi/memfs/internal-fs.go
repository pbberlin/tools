package memfs

import (
	"log"
	"path"
	// "path/filepath"

	"github.com/pbberlin/tools/os/fsi"
)

// try to find parent in main map-of-files
func (m *memMapFs) findParent(name string) (fsi.File, error) {

	if name == m.RootName() || name == m.RootDir() {
		// dont find root for root
		return nil, nil
	}

	parent, _ := m.SplitX(name)
	if len(parent) > 0 {
		pfile, err := m.Open(parent)
		if err != nil {
			return nil, err
		}
		return pfile, nil
	}
	return nil, nil
}

// adding f to it's parent's directory
func (m *memMapFs) registerWithParent(name string) string {

	if name == "" {
		return ""
	}

	// first try
	pDir, err := m.findParent(name)
	if err != nil {
		// first try -dont care
	}
	if pDir == nil {
		newPar, _ := m.SplitX(name)
		if len(newPar) > 0 {
			// log.Printf("  create parent %-32q for %v\n", newPar, name)
			err := m.MkdirAll(newPar, 0777)
			if err != nil && err != fsi.ErrFileExists {
				log.Printf("Mkdir for %v failed %v", newPar, err)
			}
		}
	}

	// trying again after creation
	pDir, err = m.findParent(name)
	if err != nil {
		log.Printf("  parent for %v should now be there \n", name)
		return ""
	}
	if pDir != nil {
		pDirC := pDir.(*InMemoryFile)

		f, _ := m.Open(name)
		// ff, _ := f.(*InMemoryFile)
		pDirC.memDir[name] = f

		// log.Printf("    fo %-32q got added to %32q\n", name, pDirC.name)

		return pDirC.name

	}

	return ""

}

// recursing upwards
func (m *memMapFs) registerDirs(name string) {

	parent, bname := m.SplitX(name)
	name = path.Join(parent, bname)

	// register upwardly
	cntr := 0
	for {

		cntr++
		if cntr > 10 {
			break //
		}

		par := m.registerWithParent(name)
		// log.Printf("register loop %v %-22q %-22q \n", cntr, name, par)
		if par == "" || par == m.RootDir() || par == m.RootName() {
			break
		}
		name = par

	}
}

// remove f from it's parent's directory
func (m *memMapFs) unRegisterWithParent(name string) error {

	parent, bname := m.SplitX(name)
	name = path.Join(parent, bname)

	pDir, err := m.findParent(name)
	if err != nil {
		return err
	}

	pDirC := pDir.(*InMemoryFile)
	delete(pDirC.memDir, name)
	// log.Printf("unregistered %-22q in %q \n", name, pDirC.name)

	return nil
}
