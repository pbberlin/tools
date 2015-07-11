package aefs

import (
	"bytes"
	"os"

	pth "path"

	"github.com/pbberlin/tools/logif"

	"appengine"
)

func RetrieveDirs(c appengine.Context, rts string) *bytes.Buffer {

	bb := new(bytes.Buffer)

	fs := NewAeFs(rts, AeContext(c))
	wpf(bb, "created fs %v\n", rts)

	wpf(bb, "--------retrieve by query---------\n")

	fc3 := func(path string, direct bool) {
		wpf(bb, "searching ---  %v\n", path)
		children, err := fs.subdirsByPath(path, direct)
		if err != nil {
			wpf(bb, "   nothing retrieved - err %v\n", err)
		} else {
			for k, v := range children {
				wpf(bb, "child #%2v %-24v\n", k, v.Name())
			}
		}
	}

	fc3(spf(`ch1/ch2/ch3`), false)
	fc3(spf(`ch1/ch2/ch3`), true)
	fc3(spf(`ch1`), false)
	fc3(spf(`ch1`), true)
	fc3(spf(``), true)
	fc3(spf(``), false)

	return bb

}

func CreateSys(c appengine.Context, rts string) *bytes.Buffer {

	bb := new(bytes.Buffer)

	fs := NewAeFs(rts, AeContext(c))
	wpf(bb, "created fs %v\n", rts)

	fc1 := func(p []string) {
		path := pth.Join(p...)
		err := fs.MkdirAll(path, os.ModePerm)
		logif.E(err)
	}

	wpf(bb, "--------create-dirs---------\n")

	fc1([]string{"ch1"})
	fc1([]string{"ch1", "ch2"})
	fc1([]string{"ch1", "ch2a"})
	fc1([]string{"ch1", "ch2", "ch3"})
	fc1([]string{"ch1", "ch2", "ch3", "ch4"})
	fc1([]string{"ch1A"})
	fc1([]string{"ch1B"})
	fc1([]string{"ch1", "d2", "d3", "d4"})
	fc1([]string{"d1", "d2", "d3", "d4"})

	wpf(bb, "\n--------retrieve-dirs---------\n")

	// retrieval
	fc2 := func(p []string) {
		path := pth.Join(p...)
		wpf(bb, "searching... %v\n", path)
		f, err := fs.dirByPath(path)
		if err != nil {
			wpf(bb, "   nothing retrieved - err %v\n", err)
		} else {
			wpf(bb, "   fnd %v \n", f.Name())
		}
	}
	fc2([]string{"ch1"})
	fc2([]string{"ch1", "ch2"})
	fc2([]string{"ch1", "cha2"})
	fc2([]string{"ch1", "ch2", "ch3"})
	fc2([]string{"ch1A"})
	fc2([]string{fs.RootDir()})

	wpf(bb, "\n-------create and save some files----\n")

	fc4a := func(name, content string) {
		err := fs.WriteFile(name, []byte(content), os.ModePerm)
		logif.E(err)
	}
	fc4b := func(name, content string) {
		f, err := fs.Create(name)
		logif.E(err)
		if err != nil {
			return
		}
		_, err = f.WriteString(content)
		logif.E(err)
		err = f.Sync()
		logif.E(err)
	}

	fc4a("ch1/ch2/file1", "content 1")
	fc4b("ch1/ch2/file2", "content 2")
	fc4a("ch1/ch2/ch3/file3", "another content")
	fc4b(fs.RootDir()+"file4", "chq content 2")

	wpf(bb, "\n-------retrieve files again----\n")

	fc5 := func(path string) {
		files, err := fs.filesByPath(fs.RootDir() + path)
		logif.E(err)
		wpf(bb, " srch %v  \n", fs.RootDir()+path)
		for k, v := range files {
			wpf(bb, "     %v  -  %v %s\n", k, v.Name(), v.Data)
		}
	}

	fc5("ch1/ch2")
	fc5("ch1/ch2/ch3")
	fc5("")

	return bb
}
