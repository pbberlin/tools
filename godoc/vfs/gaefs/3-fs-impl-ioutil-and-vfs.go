package gaefs

import (
	"os"
	"sort"

	"appengine/datastore"

	"github.com/pbberlin/tools/logif"

	pth "path"
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
// Todo: Sort dirs by name
func (fs *AeFileSys) ReadDir(path string) ([]os.FileInfo, error) {

	path = cleanseLeadingSlash(path)

	var dirs []AeDir
	var fis []os.FileInfo

	dir, err := fs.GetDirByPath(path)
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
		fs.Ctx().Errorf("Error getching dir children of %v => %v", dir.Key, err)
		return fis, err
	}

	for _, v := range dirs {
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

func (fs *AeFileSys) WriteFile(filename string, data []byte, perm os.FileMode) error {
	f := AeFile{}
	f.BName = pth.Base(filename)
	f.Dir = pth.Dir(filename)
	f.Data = data

	err := fs.SaveFile(&f, filename)
	return err
}
