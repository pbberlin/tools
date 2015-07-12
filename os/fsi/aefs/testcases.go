package aefs

import (
	"bytes"
	"os"
	"strings"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/fsc"

	pth "path"

	"appengine"
)

func CreateSys(c appengine.Context) *bytes.Buffer {

	bb := new(bytes.Buffer)

	fs := NewAeFs(MountPointNext(), AeContext(c))
	wpf(bb, "created fs %v\n\n", fs.RootDir())

	fc1 := func(p []string) {
		path := pth.Join(p...)
		err := fs.MkdirAll(path, os.ModePerm)
		if err != nil {
			wpf(bb, "MkdirAll failed %v\n", err)
		}
	}

	wpf(bb, "--------create-dirs---------\n")

	fc1([]string{"ch1"})
	fc1([]string{"ch1", "ch2"})
	fc1([]string{"ch1", "ch2a"})
	fc1([]string{"ch1", "ch2", "ch3"})
	fc1([]string{"ch1", "ch2", "ch3", "ch4"})
	fc1([]string{"ch1A"})
	fc1([]string{"ch1B"})
	fc1([]string{"d1", "d2", "d3_secretdir", "neverwalked"})
	fc1([]string{"d1", "d2", "d3a", "willwalk"})

	wpf(bb, "\n--------retrieve-dirs---------\n")

	// retrieval
	gotByPath := 0
	wntByPath := 5
	fc2 := func(p []string) {
		path := pth.Join(p...)
		wpf(bb, "searching... %v\n", path)
		f, err := fs.Lstat(path)
		if err != nil {
			wpf(bb, "   nothing retrieved - err %v\n", err)
		} else {
			wpf(bb, "   fnd %v \n", f.Name())
			gotByPath++
		}
	}
	fc2([]string{"ch1"})
	fc2([]string{"ch1", "ch2"})
	fc2([]string{"ch1", "non-exist-dir"})
	fc2([]string{"ch1", "ch2", "ch3"})
	fc2([]string{"ch1A"})
	fc2([]string{fs.RootDir()})

	wpf(bb, "\nfnd %v of %v dirs \n", gotByPath, wntByPath)

	wpf(bb, "\n-------create and save some files----\n")

	fc4a := func(name, content string) {
		err := fs.WriteFile(name, []byte(content), os.ModePerm)
		if err != nil {
			wpf(bb, "WriteFile %v failed %v\n", name, err)
		}
	}
	fc4b := func(name, content string) {
		f, err := fs.Create(name)
		if err != nil {
			wpf(bb, "Create %v failed %v\n", name, err)
		}
		if err != nil {
			return
		}
		_, err = f.WriteString(content)
		if err != nil {
			wpf(bb, "WriteString %v failed %v\n", f.Name(), err)
		}
		err = f.Sync()
		if err != nil {
			wpf(bb, "Sync %v failed %v\n", f.Name(), err)
		}
	}

	fc4a("ch1/ch2/file_1", "content 1")
	fc4b("ch1/ch2/file_2", "content 2")
	fc4a("ch1/ch2/ch3/file3", "another content")
	fc4b(fs.RootDir()+"file4", "chq content 2")

	wpf(bb, "\n-------retrieve files again----\n")

	gotNumFiles := 0
	wntNumFiles := 4
	gotSizeFiles := 0
	wntSizeFiles := 9 + 9 + 15 + 13

	fc5 := func(path string) {
		wpf(bb, " srch %v  \n", path)
		files, err := fs.ReadDir(path)
		if err != nil {
			wpf(bb, "filesByPath %v failed %v\n", path, err)
		}

		for k, v := range files {
			if v.IsDir() {
				continue
			}
			data, err := fs.ReadFile(v.Name())
			if err != nil {
				wpf(bb, "could not get content of %v =>  %v\n", v.Name(), err)
			}
			wpf(bb, "     %v  -  %v %s\n", k, v.Name(), data)
			gotNumFiles++
			gotSizeFiles += len(data)
		}
	}

	fc5("ch1/ch2")
	fc5("ch1/ch2/ch3")
	fc5("")

	wpf(bb, "\n")
	wpf(bb, "fnd %2v of %2v fils \n", gotNumFiles, wntNumFiles)
	wpf(bb, "fnd %2v of %2v fsize \n", gotSizeFiles, wntSizeFiles)
	wpf(bb, "\n")

	return bb
}

func RetrieveByQuery(c appengine.Context) *bytes.Buffer {

	bb := new(bytes.Buffer)

	fs := NewAeFs(MountPointLast(), AeContext(c))
	wpf(bb, "created fs %v\n", fs.RootDir())

	wpf(bb, "--------retrieve by query---------\n\n")

	fc3 := func(path string, direct bool) {
		mode := "direct"
		if !direct {
			mode = "all"
		}
		wpf(bb, "searching %-6v  %q\n", mode, path)
		children, err := fs.subdirsByPath(path, direct)
		if err != nil {
			wpf(bb, "   nothing retrieved - err %v\n", err)
		} else {
			for k, v := range children {
				wpf(bb, "  child #%-2v        %-24v\n", k, v.Name())
			}
		}
		wpf(bb, "\n")
	}

	fc3(`ch1/ch2/ch3`, false)
	fc3(`ch1/ch2/ch3`, true)
	fc3(`ch1/ch2`, false)
	fc3(`ch1/ch2`, true)
	fc3(`ch1`, false)
	fc3(`ch1`, true)
	fc3(``, true)
	fc3(``, false)

	return bb

}

func RetrieveByReadDir(c appengine.Context) *bytes.Buffer {

	bb := new(bytes.Buffer)

	fs := NewAeFs(MountPointLast(), AeContext(c))
	wpf(bb, "created fs %v\n", fs.RootDir())

	wpf(bb, "--------retrieve by readDir---------\n\n")

	fc3 := func(path string) {
		wpf(bb, "searching %q\n", path)
		children, err := fs.ReadDir(path)
		if err != nil {
			wpf(bb, "   nothing retrieved - err %v\n", err)
		} else {
			for k, v := range children {
				wpf(bb, "  child #%-2v        %-24v\n", k, v.Name())
			}
		}
		wpf(bb, "\n")
	}

	fc3(`ch1/ch2/ch3`)
	fc3(`ch1/ch2`)
	fc3(`ch1`)
	fc3(``)

	return bb

}

func WalkDirs(c appengine.Context) *bytes.Buffer {

	bb := new(bytes.Buffer)
	fsConcrete := NewAeFs(MountPointLast(), AeContext(c))
	fs := fsi.FileSystem(fsConcrete)

	wpf(bb, "created fs %v\n", MountPointLast())

	wpf(bb, "-------filewalk----\n\n")

	walkFunc := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			wpf(bb, "error on visiting %s => %v \n", path, err)
			return err
		}
		if strings.HasSuffix(path, "_secretdir") {
			return fsc.SkipDir // do not delve deeper
		}
		if err == os.ErrInvalid {
			return err // calling off the walk
		}
		tp := "file"
		if f != nil {
			if f.IsDir() {
				tp = "dir "
			}
		}
		wpf(bb, "Visited: %s %s \n", tp, path)
		return nil
	}

	var err error

	err = fsc.Walk(fs, "/", walkFunc)
	wpf(bb, "fs.Walk() returned %v\n\n", err)

	err = fsc.Walk(fs, "ch1/ch2", walkFunc)
	wpf(bb, "fs.Walk() returned %v\n\n", err)

	err = fsc.Walk(fs, "ch1/ch2/ch3", walkFunc)
	wpf(bb, "fs.Walk() returned %v\n\n", err)

	return bb
}

func RemoveSubtree(c appengine.Context) *bytes.Buffer {
	bb := new(bytes.Buffer)
	fsConcrete := NewAeFs(MountPointLast(), AeContext(c))
	fs := fsi.FileSystem(fsConcrete)

	wpf(bb, "created fs %v\n", MountPointLast())

	wpf(bb, "-------removedir----\n\n")
	err := fs.RemoveAll("ch1/ch2/ch3")
	wpf(bb, "fs.RemoveAll() returned %v\n\n", err)

	return bb
}
