package aefs

import (
	"bytes"
	"os"
	"strings"

	"appengine/datastore"

	"github.com/pbberlin/tools/os/fsi"
	"github.com/pbberlin/tools/os/fsi/fsc"

	pth "path"
)

// osfs1
// const rel = "."
// const relPsep = rel + sep
// const relOpt = ""

// osfs1
// const rel = "."
// const relPsep = ""
// const relOpt = ""

// memfs1
// const rel = ""
// const relPsep = rel + sep
// const relOpt = rel

// memfs2
// const rel = ""
// const relPsep = rel
// const relOpt = rel

// aefs1
// const rel = "/"
// const relPsep = "/"
// const relOpt = "/"

// aefs2
// const rel = ""
// const relPsep = ""
// const relOpt = ""

const relPsep = ""
const relOpt = ""

var rel = "."

func CreateSys(fs fsi.FileSystem) (*bytes.Buffer, string) {

	bb := new(bytes.Buffer)
	wpf(bb, "--------create-dirs---------\n")

	fc1 := func(p []string) {
		path := pth.Join(p...)
		err := fs.MkdirAll(relOpt+path, os.ModePerm)
		if err != nil {
			wpf(bb, "MkdirAll failed %v\n", err)
		}
	}

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
		wpf(bb, "searching... %q\n", path)
		f, err := fs.Lstat(relOpt + path)
		if err != nil {
			wpf(bb, "   nothing retrieved - err %v\n", err)
		} else {
			wpf(bb, "   fnd %v \n", pth.Join(path, f.Name()))
			gotByPath++
		}
	}
	fc2([]string{"ch1"})
	fc2([]string{"ch1", "ch2"})
	fc2([]string{"ch1", "non-exist-dir"})
	fc2([]string{"ch1", "ch2", "ch3"})
	fc2([]string{"ch1A"})
	fc2([]string{rel})

	wpf(bb, "\nfnd %v of %v dirs \n", gotByPath, wntByPath)

	wpf(bb, "\n-------create and save some files----\n")

	fc4a := func(name, content string) {
		err := fs.WriteFile(relOpt+name, []byte(content), os.ModePerm)
		if err != nil {
			wpf(bb, "WriteFile %v failed %v\n", name, err)
		}
	}
	fc4b := func(name, content string) {
		f, err := fs.Create(relOpt + name)
		if err != nil {
			wpf(bb, "Create %v failed %v\n", name, err)
		}
		if err != nil {
			return
		}
		_, err = f.WriteString(content)
		if err != nil {
			wpf(bb, "WriteString %v failed %v\n", name, err)
		}
		err = f.Close()
		if err != nil {
			wpf(bb, "Sync %v failed %v\n", name, err)
		}
	}

	fc4a("ch1/ch2/file_1", "content 1")
	fc4b("ch1/ch2/file_2", "content 2")
	fc4a("ch1/ch2/ch3/file3", "another content")
	fc4b(relPsep+"file4", "chq content 2")

	// fsc, ok := memfs.Unwrap(fs)
	// if ok {
	// 	fsc.Dump()
	// }
	// return bb, ""

	wpf(bb, "\n-------retrieve files again----\n\n")

	gotNumFiles := 0
	wntNumFiles := 4
	gotSizeFiles := 0
	wntSizeFiles := 9 + 9 + 15 + 13

	fc5 := func(path string) {
		wpf(bb, " srch %q  \n", relOpt+path)
		files, err := fs.ReadDir(relOpt + path)
		if err != nil {
			wpf(bb, "filesByPath %v failed %v\n", path, err)
		}

		for k, v := range files {
			if v.IsDir() {
				continue
			}
			data, err := fs.ReadFile(pth.Join(path, v.Name()))
			if err != nil {
				wpf(bb, "could not get content of %v =>  %v\n", pth.Join(path, v.Name()), err)
			}
			wpf(bb, "     %v  -  %v %s\n", k, pth.Join(path, v.Name()), data)
			gotNumFiles++
			gotSizeFiles += len(data)
		}
	}

	fc5("ch1/ch2")
	fc5("ch1/ch2/ch3")
	fc5(rel)

	wpf(bb, "\n")

	wpf(bb, "fnd %2v of %2v fils \n", gotNumFiles, wntNumFiles)
	wpf(bb, "fnd %2v of %2v fsize \n", gotSizeFiles, wntSizeFiles)
	wpf(bb, "\n")

	testRes := ""
	if gotNumFiles != wntNumFiles {
		testRes += spf("Create:   wnt %2v - got %v\n", wntNumFiles, gotNumFiles)
	}
	if gotSizeFiles != wntSizeFiles {
		testRes += spf("Create:   wnt %2v - got %v\n", wntSizeFiles, gotSizeFiles)
	}
	return bb, testRes
}

func RetrieveByReadDir(fs fsi.FileSystem) (*bytes.Buffer, string) {

	bb := new(bytes.Buffer)
	wpf(bb, "--------retrieve by readDir---------\n\n")

	wnt1 := []int{2, 3, 2, 5}
	wnt2 := []int{2, 2, 5}
	got := []int{}

	fc3 := func(path string) {
		wpf(bb, "searching %q\n", path)
		children, err := fs.ReadDir(path)
		if err != nil {
			wpf(bb, "   nothing retrieved - err %v\n", err)
		} else {
			for k, v := range children {
				wpf(bb, "  child #%-2v        %-24v\n", k, pth.Join(path, v.Name()))
			}
			got = append(got, len(children))
		}
		wpf(bb, "\n")
	}

	fc3(`ch1/ch2/ch3`)
	fc3(`ch1/ch2`)
	fc3(`ch1`)
	fc3(rel)

	testRes := ""
	if spf("%+v", wnt1) != spf("%+v", got) &&
		spf("%+v", wnt2) != spf("%+v", got) {
		testRes = spf("ReadDir:  wnt %v or %v - got %v", wnt1, wnt2, got)
	}
	return bb, testRes

}

func RetrieveByQuery(fs fsi.FileSystem) (*bytes.Buffer, string) {

	bb := new(bytes.Buffer)

	wnt1 := []int{1, 1, 2, 1, 4, 2, 4, 13}
	wnt2 := []int{2, 2, 4, 11}
	got := []int{}

	fsConcrete, ok := fs.(*aeFileSys)
	if !ok {
		wpf(bb, "--------retrieve by query UNSUPPORTED---------\n\n")
		return bb, ""
	}

	wpf(bb, "--------retrieve by query---------\n\n")

	fc3 := func(path string, direct bool) {
		mode := "direct"
		if !direct {
			mode = "all"
		}
		wpf(bb, "searching %-6v  %q\n", mode, path)
		children, err := fsConcrete.subdirsByPath(path, direct)
		if err != nil {
			wpf(bb, "   nothing retrieved - err %v\n", err)
		} else {
			for k, v := range children {
				wpf(bb, "  child #%-2v        %-24v\n", k, pth.Join(path, v.Name()))
			}
			got = append(got, len(children))
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

	testRes := ""
	if spf("%+v", wnt1) != spf("%+v", got) &&
		spf("%+v", wnt2) != spf("%+v", got) {
		testRes = spf("IdxQuery: wnt %v or %v - got %v", wnt1, wnt2, got)
	}
	return bb, testRes

}

func WalkDirs(fs fsi.FileSystem) (*bytes.Buffer, string) {

	bb := new(bytes.Buffer)
	wpf(bb, "-------filewalk----\n\n")

	wnt := []int{16, 6, 3}
	wnt2 := []int{13, 3, 0}
	got := []int{}

	cntr := 0
	walkFunc := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			wpf(bb, "error on visiting %q => %v \n", path, err)
			if err == datastore.ErrNoSuchEntity || err == os.ErrNotExist {
				return nil // dont break the walk on this, it's just a stale directory
			}
			return err // this would break the walk on any error; notably dir-index entries, that have been deleted since.
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
		cntr++
		// if cntr > 100 {
		// 	return fmt.Errorf("too many walk recursions%v", cntr)
		// }
		return nil
	}

	var err error

	cntr = 0
	err = fsc.Walk(fs, rel, walkFunc)
	wpf(bb, "fs.Walk() returned %v\n\n", err)
	got = append(got, cntr)

	cntr = 0
	err = fsc.Walk(fs, "ch1/ch2", walkFunc)
	wpf(bb, "fs.Walk() returned %v\n\n", err)
	got = append(got, cntr)

	cntr = 0
	err = fsc.Walk(fs, "ch1/ch2/ch3", walkFunc)
	wpf(bb, "fs.Walk() returned %v\n\n", err)
	got = append(got, cntr)

	testRes := ""
	if spf("%+v", wnt) != spf("%+v", got) &&
		spf("%+v", wnt2) != spf("%+v", got) {
		testRes = spf("WalkDir:  wnt %v or %v - got %v", wnt, wnt2, got)
	}

	return bb, testRes
}

func RemoveSubtree(fs fsi.FileSystem) (*bytes.Buffer, string) {

	bb := new(bytes.Buffer)

	wpf(bb, "-------removedir----\n\n")
	err := fs.RemoveAll("ch1/ch2/ch3")
	wpf(bb, "fs.RemoveAll() returned %v\n\n", err)

	testRes := ""
	if err != nil {
		testRes = spf("RemoveTree: %v", err)
	}
	return bb, testRes
}
