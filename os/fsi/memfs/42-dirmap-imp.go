package memfs

import "github.com/pbberlin/tools/os/fsi"

func (m MemDirMap) Add(f fsi.File) { m[f.Name()] = f }
func (m MemDirMap) Len() int       { return len(m) }
func (m MemDirMap) Files() (files []fsi.File) {
	for _, f := range m {
		files = append(files, f)
	}
	return files
}
func (m MemDirMap) Names() (names []string) {
	for x := range m {
		names = append(names, x)
	}
	return names
}
func (m MemDirMap) Remove(f fsi.File) { delete(m, f.Name()) }
