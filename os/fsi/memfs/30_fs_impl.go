package memfs

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/pbberlin/tools/os/fsi"
)

func (memMapFs) Name() string { return "memfs" } // type
// instance
func (m *memMapFs) String() string {
	return m.ident
}

func (m *memMapFs) createHelper(name string) *InMemoryFile {
	return &InMemoryFile{
		name:    name,
		mode:    os.ModeTemporary,
		modtime: time.Now(),
		memDir:  map[string]fsi.File{},
		fs:      m,
	}
}

func (m *memMapFs) Chmod(name string, mode os.FileMode) error {

	dir, bname := m.SplitX(name)
	name = path.Join(dir, bname)

	f, ok := m.fos[name]
	if !ok {
		return &os.PathError{Op: "chmod", Path: name, Err: fsi.ErrFileNotFound}
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

	dir, bname := m.SplitX(name)
	name = path.Join(dir, bname)

	f, ok := m.fos[name]
	if !ok {
		return &os.PathError{Op: "chtimes", Path: name, Err: fsi.ErrFileNotFound}
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

func (m *memMapFs) Create(name string) (fsi.File, error) {

	dir, bname := m.SplitX(name)
	name = path.Join(dir, bname)

	m.lock()
	m.fos[name] = m.createHelper(name)
	m.unlock()
	m.registerDirs(name)
	return m.fos[name], nil
}

func (fs *memMapFs) Lstat(path string) (os.FileInfo, error) {
	return fs.Stat(path)
}

func (m *memMapFs) Mkdir(name string, perm os.FileMode) error {

	dir, bname := m.SplitX(name)
	name = path.Join(dir, bname)

	m.rlock()
	_, ok := m.fos[name]
	m.runlock()
	if ok {
		return fsi.ErrFileExists
	} else {
		fo := m.createHelper(name)
		fo.dir = true
		m.lock()
		m.fos[name] = fo
		m.unlock()
		m.registerDirs(name)
	}
	return nil
}

func (m *memMapFs) MkdirAll(name string, perm os.FileMode) error {
	return m.Mkdir(name, 0777)
}

func (m *memMapFs) Open(name string) (fsi.File, error) {

	origName := name

	dir, bname := m.SplitX(name)
	name = path.Join(dir, bname)

	m.rlock()
	f, ok := m.fos[name]
	if ok {
		ff, okConv := f.(*InMemoryFile)
		if okConv {
			ff.Open()
		} else {
			return nil, fmt.Errorf("could not convert opened file into InMemoryFile 1")
		}
	}
	m.runlock()

	//
	//
	// Fallback to underlying fs
	if !ok {
		var err error
		f, err = m.lookupUnderlyingFS(name, origName)
		if err != nil {
			// log.Printf("underlying says  %q => %v\n", name, err)
			return nil, err
			return nil, fsi.ErrFileNotFound
		}
		// log.Printf("underlying succ %q\n", name)
		m.Dump()
		return f, nil
	}

	//
	// Regular return
	if ok {
		return f, nil
	} else {
		return nil, fsi.ErrFileNotFound
	}
}

//
func (m *memMapFs) lookupUnderlyingFS(
	nameMemFS string,
	origName string, // orig name has no mountname prefix á la mnt02
) (fsi.File, error) {

	if m.shadow == nil { // no underlying filesystem
		return nil, fsi.ErrFileNotFound
	}

	// mshadt := m.shadow.(*)

	fshad, err := m.shadow.Open(nameMemFS)
	if err != nil {
		return nil, fsi.ErrFileNotFound
	}
	defer fshad.Close()

	inf, err := fshad.Stat()
	if err != nil {
		return nil, fmt.Errorf("fileinfo from shadow failed: %v", err)
	}

	//
	// special case
	// resource is a directory
	if inf.IsDir() {
		// return nil, fmt.Errorf("is dir")
		err = m.MkdirAll(origName, 0755)
		if err != nil {
			return nil, err
		}
		m.rlock()
		dir, ok := m.fos[nameMemFS]
		m.runlock()
		if !ok {
			m.Dump()
			return nil, fmt.Errorf("dir created with MkDir, but not in fos map %q %q", nameMemFS, origName)
		}
		return dir, nil
	}

	nameMemFS = origName

	var dst fsi.File

	// regular file
	err = m.MkdirAll(path.Dir(nameMemFS), 0755)
	if err != nil {
		return nil, err
	}
	log.Printf("  from underlying: created front dir  %q \n", path.Dir(nameMemFS))

	dst, err = m.Create(nameMemFS)
	if err != nil {
		return nil, err
	}
	log.Printf("  from underlying: created front file %q \n", nameMemFS)

	n, err := io.Copy(dst, fshad)
	_ = n
	if err != nil {
		return nil, err
	}
	// log.Printf("copied %v for %v\n", n, name)

	err = dst.Close()
	if err != nil {
		return nil, err
	}

	//
	// reopen
	ff, okConv := dst.(*InMemoryFile)
	if okConv {
		ff.Open()
	} else {
		return nil, fmt.Errorf("could not convert opened file into InMemoryFile 2")
	}

	return dst, nil

}

func (m *memMapFs) OpenFile(name string, flag int, perm os.FileMode) (fsi.File, error) {
	return m.Open(name)
}

func (fs *memMapFs) ReadDir(name string) ([]os.FileInfo, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (m *memMapFs) Remove(name string) error {

	dir, bname := m.SplitX(name)
	name = path.Join(dir, bname)

	m.rlock()
	defer m.runlock()

	if _, ok := m.fos["name"]; ok {
		m.lock()
		delete(m.fos, name)
		m.unlock()
		m.unRegisterWithParent(name) // should be inside lock-unlock - but causes deadlock
	}
	return nil
}

func (m *memMapFs) RemoveAll(name string) error {

	dir, bname := m.SplitX(name)
	name = path.Join(dir, bname)

	log.Printf("starting memfs removeall %v", name)

	m.rlock()
	defer m.runlock()
	for p, _ := range m.fos {
		// log.Printf("    removeall checking %v", p)
		if strings.HasPrefix(p, name) {
			log.Printf("    removeall deleting %v", p)
			m.runlock()
			m.lock()
			delete(m.fos, p)
			m.unlock()
			m.rlock()
			m.unRegisterWithParent(name) // now readlocked, therefore ok
		}
	}
	return nil
}

func (m *memMapFs) Rename(name, newname string) error {

	dir, bname := m.SplitX(name)
	name = path.Join(dir, bname)

	{
		dir, bname := m.SplitX(newname)
		newname = path.Join(dir, bname)
	}

	m.rlock()
	defer m.runlock()
	if _, ok := m.fos[name]; ok {
		if _, ok := m.fos[newname]; !ok {
			m.runlock()
			m.lock()
			m.fos[newname] = m.fos[name]
			delete(m.fos, name)
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

func (fs *memMapFs) ReadFile(name string) ([]byte, error) {

	f, err := fs.Open(name)
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

	keys := make([]string, 0, len(m.fos))
	for key, _ := range m.fos {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, k := range keys {
		f := m.fos[k]
		ff, ok := f.(*InMemoryFile)
		y, _ := f.Stat()
		if ok {
			names := ""
			for _, v := range ff.memDir {
				vf, _ := v.(*InMemoryFile)
				names += vf.name + "  "
			}
			log.Printf("%-38q %5v %4v %-20v\n", ff.name, y.IsDir(), y.Size(), names)
		}
	}
}
