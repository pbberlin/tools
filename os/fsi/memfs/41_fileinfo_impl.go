package memfs

import (
	"os"
	"time"
)

// Implements os.FileInfo
func (s *InMemoryFileInfo) Name() string       { return s.file.Name() }
func (s *InMemoryFileInfo) Mode() os.FileMode  { return s.file.mode }
func (s *InMemoryFileInfo) ModTime() time.Time { return s.file.modtime }
func (s *InMemoryFileInfo) IsDir() bool        { return s.file.dir }
func (s *InMemoryFileInfo) Sys() interface{}   { return nil }
func (s *InMemoryFileInfo) Size() int64 {
	if s.IsDir() {
		return int64(42)
	}
	return int64(len(s.file.data))
}

// byName implements sort.Interface.
type byName []os.FileInfo

func (f byName) Len() int           { return len(f) }
func (f byName) Less(i, j int) bool { return f[i].Name() < f[j].Name() }
func (f byName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

type byDateAsc []os.FileInfo

func (f byDateAsc) Len() int           { return len(f) }
func (f byDateAsc) Less(i, j int) bool { return f[i].ModTime().Before(f[j].ModTime()) }
func (f byDateAsc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

type byDateDesc []os.FileInfo

func (f byDateDesc) Len() int           { return len(f) }
func (f byDateDesc) Less(i, j int) bool { return f[i].ModTime().After(f[j].ModTime()) }
func (f byDateDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
