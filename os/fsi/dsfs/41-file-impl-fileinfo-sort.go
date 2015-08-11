package dsfs

import "os"

// byName implements sort.Interface.
type FileInfoByName []os.FileInfo

func (f FileInfoByName) Len() int           { return len(f) }
func (f FileInfoByName) Less(i, j int) bool { return f[i].Name() < f[j].Name() }
func (f FileInfoByName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

type AeFileByName []AeFile

func (f AeFileByName) Len() int           { return len(f) }
func (f AeFileByName) Less(i, j int) bool { return f[i].Name() < f[j].Name() }
func (f AeFileByName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
