package memfs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/pbberlin/tools/os/fsi"
)

func (memMapFs) Name() string { return "MemMapFS" }
func (memMapFs) String() string {
	hn, err := os.Hostname()
	if err != nil {
		return err.Error()
	}
	return hn
}

func MemFileCreate(name string) *InMemoryFile {
	return &InMemoryFile{name: name, mode: os.ModeTemporary, modtime: time.Now()}
}

func (m *memMapFs) Create(name string) (fsi.File, error) {
	m.lock()
	m.fos[name] = MemFileCreate(name)
	m.unlock()
	m.registerDirs(m.fos[name])
	return m.fos[name], nil
}

func (fs *memMapFs) Lstat(path string) (os.FileInfo, error) {
	return fs.Stat(path)
}

func (m *memMapFs) Mkdir(name string, perm os.FileMode) error {
	m.rlock()
	_, ok := m.fos[name]
	m.runlock()
	if ok {
		return fsi.ErrFileExists
	} else {
		m.lock()
		m.fos[name] = &InMemoryFile{name: name, memDir: &MemDirMap{}, dir: true}
		m.unlock()
		m.registerDirs(m.fos[name])
	}
	return nil
}

func (m *memMapFs) MkdirAll(name string, perm os.FileMode) error {

	return m.Mkdir(name, 0777)

	name = strings.TrimSpace(name)
	dirs := strings.Split(path.Clean(name), sep)
	// log.Printf("  MkdirAll %-22v => %v", name, dirs)
	for _, v := range dirs {
		err := m.Mkdir(v, 0777)
		// log.Printf("    MkdirAll %q %v", v, err)
		if err != nil && err != fsi.ErrFileExists {
			return err
		}
	}
	return nil
}

func (m *memMapFs) Open(name string) (fsi.File, error) {
	m.rlock()
	f, ok := m.fos[name]
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

func (m *memMapFs) OpenFile(name string, flag int, perm os.FileMode) (fsi.File, error) {
	return m.Open(name)
}

func (fs *memMapFs) ReadDir(path string) ([]os.FileInfo, error) {

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

func (m *memMapFs) Remove(name string) error {
	m.rlock()
	defer m.runlock()

	if _, ok := m.fos["name"]; ok {
		m.lock()
		delete(m.fos, name)
		m.unlock()
	}
	return nil
}

func (m *memMapFs) RemoveAll(path string) error {
	m.rlock()
	defer m.runlock()
	for p, _ := range m.fos {
		if strings.HasPrefix(p, path) {
			m.runlock()
			m.lock()
			delete(m.fos, p)
			m.unlock()
			m.rlock()
		}
	}
	return nil
}

func (m *memMapFs) Rename(oldname, newname string) error {
	m.rlock()
	defer m.runlock()
	if _, ok := m.fos[oldname]; ok {
		if _, ok := m.fos[newname]; !ok {
			m.runlock()
			m.lock()
			m.fos[newname] = m.fos[oldname]
			delete(m.fos, oldname)
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

func (m *memMapFs) Stat(name string) (os.FileInfo, error) {
	f, err := m.Open(name)
	if err != nil {
		return nil, err
	}
	return &InMemoryFileInfo{file: f.(*InMemoryFile)}, nil
}

func (m *memMapFs) Chmod(name string, mode os.FileMode) error {
	f, ok := m.fos[name]
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

func (m *memMapFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	f, ok := m.fos[name]
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

func (fs *memMapFs) ReadFile(path string) ([]byte, error) {

	f, err := fs.Open(path)
	if err != nil {
		return []byte{}, err
	}
	f1 := f.(*InMemoryFile)
	return f1.data, nil
}

func (fs *memMapFs) WriteFile(name string, data []byte, perm os.FileMode) error {

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

func (m *memMapFs) List() {
	for _, x := range m.fos {
		y, _ := x.Stat()
		fmt.Println(x.Name(), y.Size())
	}
}

func (m *memMapFs) Dump() {
	for _, f := range m.fos {

		ff, ok := f.(*InMemoryFile)
		y, _ := f.Stat()
		if ok && ff.memDir != nil {
			log.Printf("%-36q %4v   %-20v\n", f.Name(), y.Size(), ff.memDir.Names())
		} else {
			log.Printf("%-36q %4v\n", f.Name(), y.Size())
		}
	}
}
