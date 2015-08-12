package osfs

import "os"

type byDateAsc []os.FileInfo

func (f byDateAsc) Len() int           { return len(f) }
func (f byDateAsc) Less(i, j int) bool { return f[i].ModTime().Before(f[j].ModTime()) }
func (f byDateAsc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

type byDateDesc []os.FileInfo

func (f byDateDesc) Len() int           { return len(f) }
func (f byDateDesc) Less(i, j int) bool { return f[i].ModTime().After(f[j].ModTime()) }
func (f byDateDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
