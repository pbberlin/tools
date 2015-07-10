package memfs

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/pbberlin/tools/os/fsi"
)

func (MemMapFs) Name() string { return "MemMapFS" }
func (MemMapFs) String() string {
	hn, err := os.Hostname()
	if err != nil {
		return err.Error()
	}
	return hn
}

func MemFileCreate(name string) *InMemoryFile {
	return &InMemoryFile{name: name, mode: os.ModeTemporary, modtime: time.Now()}
}

func (m *MemMapFs) Create(name string) (fsi.File, error) {
	m.lock()
	m.getData()[name] = MemFileCreate(name)
	m.unlock()
	m.registerDirs(m.getData()[name])
	return m.getData()[name], nil
}

func (fs *MemMapFs) Lstat(path string) (os.FileInfo, error) {
	return fs.Stat(path)
}

func (m *MemMapFs) Mkdir(name string, perm os.FileMode) error {
	m.rlock()
	_, ok := m.getData()[name]
	m.runlock()
	if ok {
		return fsi.ErrFileExists
	} else {
		m.lock()
		m.getData()[name] = &InMemoryFile{name: name, memDir: &MemDirMap{}, dir: true}
		m.unlock()
		m.registerDirs(m.getData()[name])
	}
	return nil
}

func (m *MemMapFs) MkdirAll(path string, perm os.FileMode) error {
	return m.Mkdir(path, 0777)
}

func (m *MemMapFs) Open(name string) (fsi.File, error) {
	m.rlock()
	f, ok := m.getData()[name]
	ff, ok := f.(*InMemoryFile)
	if ok {
		ff.Open()
	}
	m.runlock()

	if ok {
		return f, nil
	} else {
		return nil, fsi.ErrFileNotFound
	}
}

func (m *MemMapFs) OpenFile(name string, flag int, perm os.FileMode) (fsi.File, error) {
	return m.Open(name)
}

func (fs *MemMapFs) ReadDir(path string) ([]os.FileInfo, error) {

	f, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Sort(byName(list))
	return list, nil

	// ===================================
	// var files1 map[string]File
	// files1 = fs.getData()
	// _ = files1

	// // or
	// var m MemDirMap
	// files2 := m.Files()
	// names2 := m.Names()
	// _, _ = files2, names2

	// return []os.FileInfo{}, nil

}

func (m *MemMapFs) Remove(name string) error {
	m.rlock()
	defer m.runlock()

	if _, ok := m.getData()["name"]; ok {
		m.lock()
		delete(m.getData(), name)
		m.unlock()
	}
	return nil
}

func (m *MemMapFs) RemoveAll(path string) error {
	m.rlock()
	defer m.runlock()
	for p, _ := range m.getData() {
		if strings.HasPrefix(p, path) {
			m.runlock()
			m.lock()
			delete(m.getData(), p)
			m.unlock()
			m.rlock()
		}
	}
	return nil
}

func (m *MemMapFs) Rename(oldname, newname string) error {
	m.rlock()
	defer m.runlock()
	if _, ok := m.getData()[oldname]; ok {
		if _, ok := m.getData()[newname]; !ok {
			m.runlock()
			m.lock()
			m.getData()[newname] = m.getData()[oldname]
			delete(m.getData(), oldname)
			m.unlock()
			m.rlock()
		} else {
			return fsi.ErrDestinationExists
		}
	} else {
		return fsi.ErrFileNotFound
	}
	return nil
}

func (m *MemMapFs) Stat(name string) (os.FileInfo, error) {
	f, err := m.Open(name)
	if err != nil {
		return nil, err
	}
	return &InMemoryFileInfo{file: f.(*InMemoryFile)}, nil
}

func (m *MemMapFs) Chmod(name string, mode os.FileMode) error {
	f, ok := m.getData()[name]
	if !ok {
		return &os.PathError{"chmod", name, fsi.ErrFileNotFound}
	}

	ff, ok := f.(*InMemoryFile)
	if ok {
		m.lock()
		ff.mode = mode
		m.unlock()
	} else {
		return errors.New("Unable to Chmod Memory File")
	}
	return nil
}

func (m *MemMapFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	f, ok := m.getData()[name]
	if !ok {
		return &os.PathError{"chtimes", name, fsi.ErrFileNotFound}
	}

	ff, ok := f.(*InMemoryFile)
	if ok {
		m.lock()
		ff.modtime = mtime
		m.unlock()
	} else {
		return errors.New("Unable to Chtime Memory File")
	}
	return nil
}

func (fs *MemMapFs) ReadFile(path string) ([]byte, error) {

	f, err := fs.Open(path)
	if err != nil {
		return []byte{}, err
	}
	f1 := f.(*InMemoryFile)
	return f1.data, nil
}

func (fs *MemMapFs) WriteFile(name string, data []byte, perm os.FileMode) error {

	f, err := fs.Create(name)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// other
// -----------------------------------------

func (m *MemMapFs) List() {
	for _, x := range m.data {
		y, _ := x.Stat()
		fmt.Println(x.Name(), y.Size())
	}
}
