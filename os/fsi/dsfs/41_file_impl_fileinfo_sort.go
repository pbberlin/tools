package dsfs

import "os"

// byName implements sort.Interface.
type FileInfoByName []os.FileInfo

func (f FileInfoByName) Len() int           { return len(f) }
func (f FileInfoByName) Less(i, j int) bool { return f[i].Name() < f[j].Name() }
func (f FileInfoByName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

type DsFileByName []DsFile

func (f DsFileByName) Len() int           { return len(f) }
func (f DsFileByName) Less(i, j int) bool { return f[i].Name() < f[j].Name() }
func (f DsFileByName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

//
//
type FileInfoByDateAsc []os.FileInfo

func (f FileInfoByDateAsc) Len() int           { return len(f) }
func (f FileInfoByDateAsc) Less(i, j int) bool { return f[i].ModTime().Before(f[j].ModTime()) }
func (f FileInfoByDateAsc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

type DsFileByDateAsc []DsFile

func (f DsFileByDateAsc) Len() int           { return len(f) }
func (f DsFileByDateAsc) Less(i, j int) bool { return f[i].ModTime().Before(f[j].ModTime()) }
func (f DsFileByDateAsc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

//
//
type FileInfoByDateDesc []os.FileInfo

func (f FileInfoByDateDesc) Len() int           { return len(f) }
func (f FileInfoByDateDesc) Less(i, j int) bool { return f[i].ModTime().After(f[j].ModTime()) }
func (f FileInfoByDateDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

type DsFileByDateDesc []DsFile

func (f DsFileByDateDesc) Len() int           { return len(f) }
func (f DsFileByDateDesc) Less(i, j int) bool { return f[i].ModTime().After(f[j].ModTime()) }
func (f DsFileByDateDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
