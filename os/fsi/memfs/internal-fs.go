package memfs

import (
	"path"
	"path/filepath"

	"github.com/pbberlin/tools/os/fsi"
)

func (m *MemMapFs) registerDirs(f fsi.File) {
	var x = f.Name()
	for x != "/" {
		f := m.registerWithParent(f)
		if f == nil {
			break
		}
		x = f.Name()
	}
}

func (m *MemMapFs) unRegisterWithParent(f fsi.File) fsi.File {
	parent := m.findParent(f)
	pmem := parent.(*InMemoryFile)
	pmem.memDir.Remove(f)
	return parent
}

func (m *MemMapFs) registerWithParent(f fsi.File) fsi.File {
	if f == nil {
		return nil
	}
	parent := m.findParent(f)
	if parent != nil {
		pmem := parent.(*InMemoryFile)
		pmem.memDir.Add(f)
	} else {
		pdir := filepath.Dir(path.Clean(f.Name()))
		m.Mkdir(pdir, 0777)
	}
	return parent
}

func (m *MemMapFs) findParent(f fsi.File) fsi.File {
	dirs, _ := path.Split(f.Name())
	if len(dirs) > 1 {
		_, parent := path.Split(path.Clean(dirs))
		if len(parent) > 0 {
			pfile, err := m.Open(parent)
			if err != nil {
				return pfile
			}
		}
	}
	return nil
}
