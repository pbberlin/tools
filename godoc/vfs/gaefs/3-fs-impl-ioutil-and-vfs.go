package gaefs

import (
	"os"
	"sort"

	"appengine/datastore"

	"github.com/pbberlin/tools/logif"
)

func OS(mount string) AeFileSys {
	panic(`
		Sadly, google app engine file system requires a
	 	http.Request based context object.
	 	Use NewFs(string, AeContext(c)) instead of OS.
	`)
}

func ReadFile(fs *AeFileSys, path string) ([]byte, error) {

	file, err := fs.GetFile(path)
	if err != nil {
		return []byte{}, err
	}
	return file.Data, err
}

// ReadDir satisfies the vfs interface
// and ioutil.ReadDir.
// It is similar to GetFiles, but returning only dirs
func (fs *AeFileSys) ReadDir(path string) ([]os.FileInfo, error) {

	path = cleanseLeadingSlash(path)

	var dirs []AeDir
	var fis []os.FileInfo

	dir, err := fs.GetDirByPath(path)
	if path != dir.Dir+dir.BName {
		// panic(spf("path %v must equal dir and base %v %v ", path, dir.Dir, dir.BName))
	}
	logif.Pf("%15v => %24v", path, "")

	if err == datastore.ErrNoSuchEntity {
		return fis, err
	} else if err != nil {
		logif.E(err)
		return fis, err
	}

	q := datastore.NewQuery(tdir).Ancestor(dir.Key)
	keys, err := q.GetAll(fs.Ctx(), &dirs)
	_ = keys
	if err != nil {
		fs.Ctx().Errorf("Error fetching dir children of %v => %v", dir.Key, err)
		return fis, err
	}

	for i, v := range dirs {
		pK := keys[i].Parent()
		if pK != nil && !pK.Equal(dir.Key) {
			logif.Pf("%15v =>    skp %-17v", "", v.Dir+v.BName)
			continue
		}
		logif.Pf("%15v => %-24v", "", v.Dir+v.BName)
		fi := os.FileInfo(v)
		fis = append(fis, fi)
	}

	sort.Sort(FileInfoByName(fis))

	return fis, err

}

func (fs *AeFileSys) Readdirnames(path string) (names []string, err error) {
	fis, err := fs.ReadDir(path)
	names = make([]string, 0, len(fis))
	for _, lp := range fis {
		names = append(names, lp.Name())
	}
	return names, err
}
